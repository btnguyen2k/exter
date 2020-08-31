package gvabe

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/prom"

	"main/src/goapi"
	"main/src/gvabe/bo/app"
	"main/src/gvabe/bo/session"
	"main/src/gvabe/bo/user"
	"main/src/henge"
)

func _createSqlConnect(dbtype string) *prom.SqlConnect {
	switch dbtype {
	case "sqlite":
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		return henge.NewSqliteConnection(dir, dbname)
	case "pg", "pgsql", "postgres", "postgresql":
		url := goapi.AppConfig.GetString("gvabe.db.pgsql.url")
		return henge.NewPgsqlConnection(url, goapi.AppConfig.GetString("timezone"))
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func _createAppDao(sqlc *prom.SqlConnect) app.AppDao {
	return app.NewAppDaoSql(sqlc, app.TableApp)
}

func _createUserDao(sqlc *prom.SqlConnect) user.UserDao {
	return user.NewUserDaoSql(sqlc, user.TableUser)
}

func _createSessionDao(sqlc *prom.SqlConnect) session.SessionDao {
	return session.NewSessionDaoSql(sqlc, session.TableSession)
}

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype)
	switch dbtype {
	case "sqlite":
		henge.InitSqliteTable(sqlc, user.TableUser, nil)
		henge.InitSqliteTable(sqlc, app.TableApp, map[string]string{app.ColApp_UserId: "VARCHAR(32)"})
		henge.CreateIndex(sqlc, app.TableApp, false, []string{app.ColApp_UserId})
		henge.InitSqliteTable(sqlc, session.TableSession, map[string]string{
			session.ColSession_IdSource:    "VARCHAR(32)",
			session.ColSession_AppId:       "VARCHAR(32)",
			session.ColSession_UserId:      "VARCHAR(32)",
			session.ColSession_SessionType: "VARCHAR(32)",
			session.ColSession_Expiry:      "DATETIME",
		})
		henge.CreateIndex(sqlc, session.TableSession, false, []string{session.ColSession_IdSource})
		henge.CreateIndex(sqlc, session.TableSession, false, []string{session.ColSession_AppId})
		henge.CreateIndex(sqlc, session.TableSession, false, []string{session.ColSession_Expiry})
	case "pg", "pgsql", "postgres", "postgresql":
		henge.InitPgsqlTable(sqlc, user.TableUser, nil)
		henge.InitPgsqlTable(sqlc, app.TableApp, map[string]string{app.ColApp_UserId: "VARCHAR(32)"})
		henge.CreateIndex(sqlc, app.TableApp, false, []string{app.ColApp_UserId})
		henge.InitPgsqlTable(sqlc, session.TableSession, map[string]string{
			session.ColSession_IdSource:    "VARCHAR(32)",
			session.ColSession_AppId:       "VARCHAR(32)",
			session.ColSession_UserId:      "VARCHAR(32)",
			session.ColSession_SessionType: "VARCHAR(32)",
			session.ColSession_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
		henge.CreateIndex(sqlc, session.TableSession, false, []string{session.ColSession_IdSource})
		henge.CreateIndex(sqlc, session.TableSession, false, []string{session.ColSession_AppId})
		henge.CreateIndex(sqlc, session.TableSession, false, []string{session.ColSession_Expiry})
	}
	userDao = _createUserDao(sqlc)
	appDao = _createAppDao(sqlc)
	sessionDao = _createSessionDao(sqlc)

	_initUsers()
	_initApps()
}

func _initUsers() {
	systemAppOwnerId = goapi.AppConfig.GetString("gvabe.init.system_app_owner_id")
	systemAppOwnerId = strings.ToLower(strings.TrimSpace(systemAppOwnerId))
	if systemAppOwnerId == "" {
		panic("owner id of system app not found at config [gvabe.init.system_app_owner_id]")
	}
	systemAppOwner, err := userDao.Get(systemAppOwnerId)
	if err != nil {
		panic("error while getting user [" + systemAppOwnerId + "]: " + err.Error())
	}
	if systemAppOwner == nil {
		log.Printf("System app owner [%s] not found, creating one...", systemAppOwnerId)
		systemAppOwner = user.NewUser(goapi.AppVersionNumber, systemAppOwnerId)
		result, err := userDao.Create(systemAppOwner)
		if err != nil {
			panic("error while creating user [" + systemAppOwnerId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create user [%s]", systemAppOwnerId)
		}
	}
}

func _initApps() {
	systemApp, err := appDao.Get(systemAppId)
	if err != nil {
		panic("error while getting app [" + systemAppId + "]: " + err.Error())
	}
	if systemApp == nil {
		log.Printf("System app [%s] not found, creating one...", systemAppId)
		systemApp = app.NewApp(goapi.AppVersionNumber, systemAppId, systemAppOwnerId, systemAppDesc)
		attrsPublic := systemApp.GetAttrsPublic()
		attrsPublic.IdentitySources = enabledLoginChannels
		attrsPublic.Tags = []string{systemAppId}
		systemApp.SetAttrsPublic(attrsPublic)
		result, err := appDao.Create(systemApp)
		if err != nil {
			panic("error while creating app [" + systemAppId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create app [%s]", systemAppId)
		}
	}

	// AB#13: sync the public key with exter app record in database
	if systemApp != nil {
		pubBlock := &pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PublicKey(rsaPubKey),
		}
		publicPEM := pem.EncodeToMemory(pubBlock)
		attrsPublic := systemApp.GetAttrsPublic()
		attrsPublic.RsaPublicKey = string(publicPEM)
		systemApp.SetAttrsPublic(attrsPublic)
		result, err := appDao.Update(systemApp)
		if err != nil {
			panic("error while syncing RSA public key for app [" + systemAppId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot sync RSA public key for app [%s]", systemAppId)
		}
	}
}

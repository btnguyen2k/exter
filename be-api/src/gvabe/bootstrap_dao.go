package gvabe

import (
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/prom"

	"main/src/goapi"
	"main/src/gvabe/bo/app"
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

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype)
	switch dbtype {
	case "sqlite":
		henge.InitSqliteTable(sqlc, user.TableUser, nil)
		henge.InitSqliteTable(sqlc, app.TableApp, map[string]string{app.ColApp_UserId: "VARCHAR(64)"})
		henge.CreateIndex(sqlc, app.TableApp, false, []string{app.ColApp_UserId})
	case "pg", "pgsql", "postgres", "postgresql":
		henge.InitPgsqlTable(sqlc, user.TableUser, nil)
		henge.InitPgsqlTable(sqlc, app.TableApp, map[string]string{app.ColApp_UserId: "VARCHAR(64)"})
		henge.CreateIndex(sqlc, app.TableApp, false, []string{app.ColApp_UserId})
	}
	userDao = _createUserDao(sqlc)
	appDao = _createAppDao(sqlc)

	_initUsers()
	_initApps()
}

func _initUsers() {
	systemAdminId = goapi.AppConfig.GetString("gvabe.init.system_admin_id")
	systemAdminId = strings.ToLower(strings.TrimSpace(systemAdminId))
	if systemAdminId == "" {
		panic("system admin user account id is not found at config [gvabe.init.system_admin_id]")
	}
	systemAdminUser, err := userDao.Get(systemAdminId)
	if err != nil {
		panic("error while getting user [" + systemAdminId + "]: " + err.Error())
	}
	if systemAdminUser == nil {
		log.Printf("System admin user [%s] not found, creating one...", systemAdminId)
		systemAdminUser = user.NewUser(goapi.AppVersionNumber, systemAdminId)
		result, err := userDao.Create(systemAdminUser)
		if err != nil {
			panic("error while creating user [" + systemAdminId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create user [%s]", systemAdminId)
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
		systemApp = app.NewApp(goapi.AppVersionNumber, systemAppId, systemAdminId, systemAppDesc)
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
}

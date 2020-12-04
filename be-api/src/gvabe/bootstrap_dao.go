package gvabe

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"

	"main/src/goapi"
	"main/src/gvabe/bo/app"
	"main/src/gvabe/bo/session"
	"main/src/gvabe/bo/user"
)

func _createSqlConnect(dbtype string) *prom.SqlConnect {
	var poolOpts *prom.SqlPoolOptions = nil
	timezone := goapi.AppConfig.GetString("timezone")
	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	var dsn, driver string
	var dbflavor = prom.FlavorDefault
	switch dbtype {
	case "sqlite":
		driver = "sqlite3"
		dbflavor = prom.FlavorSqlite
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		os.MkdirAll(dir, 0711)
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		dsn = dir + "/" + dbname + ".db"
	case "pg", "pgsql", "postgres", "postgresql":
		driver = "pgx"
		dbflavor = prom.FlavorPgSql
		dsn = goapi.AppConfig.GetString("gvabe.db.pgsql.url")
	}
	if driver != "" && dsn != "" {
		dsn = strings.ReplaceAll(dsn, "${loc}", urlTimezone)
		dsn = strings.ReplaceAll(dsn, "${tz}", urlTimezone)
		dsn = strings.ReplaceAll(dsn, "${timezone}", urlTimezone)
		sqlc, err := prom.NewSqlConnectWithFlavor(driver, dsn, 2345, poolOpts, dbflavor)
		if err != nil {
			panic(err)
		}
		return sqlc
	}
	return nil
}

func _createAwsDynamodbConnect(dbtype string) *prom.AwsDynamodbConnect {
	switch dbtype {
	case "dynamodb", "awsdynamodb", "aws_dynamodb", "aws-dynamodb":
		region := goapi.AppConfig.GetString("gvabe.db.dynamodb.region")
		cfg := &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewEnvCredentials(),
		}
		endpoint := goapi.AppConfig.GetString("gvabe.db.dynamodb.endpoint")
		if endpoint != "" {
			cfg.Endpoint = aws.String(endpoint)
			if strings.HasPrefix(endpoint, "http://") {
				cfg.DisableSSL = aws.Bool(true)
			}
		}
		adc, err := prom.NewAwsDynamodbConnect(cfg, nil, nil, 2345)
		if err != nil {
			panic(err)
		}
		return adc
	}
	return nil
}

func _createAppDaoAwsDynamodb(dync *prom.AwsDynamodbConnect) app.AppDao {
	return app.NewAppDaoAwsDynamodb(dync, app.TableApp)
}
func _createAppDaoSql(sqlc *prom.SqlConnect) app.AppDao {
	return app.NewAppDaoSql(sqlc, app.TableApp)
}

func _createUserDaoAwsDynamodb(dync *prom.AwsDynamodbConnect) user.UserDao {
	return user.NewUserDaoAwsDynamodb(dync, user.TableUser)
}
func _createUserDaoSql(sqlc *prom.SqlConnect) user.UserDao {
	return user.NewUserDaoSql(sqlc, user.TableUser)
}

func _createSessionDaoSql(sqlc *prom.SqlConnect) session.SessionDao {
	return session.NewSessionDaoSql(sqlc, session.TableSession)
}

func _createSessionDaoAwsDynamodb(dync *prom.AwsDynamodbConnect) session.SessionDao {
	return session.NewSessionDaoAwsDynamodb(dync, session.TableSession)
}

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype)
	dync := _createAwsDynamodbConnect(dbtype)
	if sqlc == nil && dync == nil {
		panic(fmt.Sprintf("unsupported database type: %s", dbtype))
	}
	switch dbtype {
	case "sqlite":
		henge.InitSqliteTable(sqlc, user.TableUser, nil)
		henge.InitSqliteTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "VARCHAR(32)"})
		henge.InitSqliteTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "VARCHAR(32)",
			session.SqlCol_Session_AppId:       "VARCHAR(32)",
			session.SqlCol_Session_UserId:      "VARCHAR(32)",
			session.SqlCol_Session_SessionType: "VARCHAR(32)",
			session.SqlCol_Session_Expiry:      "TIMESTAMP",
		})
	case "pg", "pgsql", "postgres", "postgresql":
		henge.InitPgsqlTable(sqlc, user.TableUser, nil)
		henge.InitPgsqlTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "VARCHAR(32)"})
		henge.InitPgsqlTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "VARCHAR(32)",
			session.SqlCol_Session_AppId:       "VARCHAR(32)",
			session.SqlCol_Session_UserId:      "VARCHAR(32)",
			session.SqlCol_Session_SessionType: "VARCHAR(32)",
			session.SqlCol_Session_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
	case "mysql":
		henge.InitMysqlTable(sqlc, user.TableUser, nil)
		henge.InitMysqlTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "VARCHAR(32)"})
		henge.InitMysqlTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "VARCHAR(32)",
			session.SqlCol_Session_AppId:       "VARCHAR(32)",
			session.SqlCol_Session_UserId:      "VARCHAR(32)",
			session.SqlCol_Session_SessionType: "VARCHAR(32)",
			session.SqlCol_Session_Expiry:      "DATETIME",
		})
	case "mssql":
		henge.InitMssqlTable(sqlc, user.TableUser, nil)
		henge.InitMssqlTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "NVARCHAR(32)"})
		henge.InitMssqlTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "NVARCHAR(32)",
			session.SqlCol_Session_AppId:       "NVARCHAR(32)",
			session.SqlCol_Session_UserId:      "NVARCHAR(32)",
			session.SqlCol_Session_SessionType: "NVARCHAR(32)",
			session.SqlCol_Session_Expiry:      "DATETIMEOFFSET",
		})
	case "oracle":
		henge.InitOracleTable(sqlc, user.TableUser, nil)
		henge.InitOracleTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "NVARCHAR2(32)"})
		henge.InitOracleTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "NVARCHAR2(32)",
			session.SqlCol_Session_AppId:       "NVARCHAR2(32)",
			session.SqlCol_Session_UserId:      "NVARCHAR2(32)",
			session.SqlCol_Session_SessionType: "NVARCHAR2(32)",
			session.SqlCol_Session_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
	}

	if sqlc != nil {
		henge.CreateIndexSql(sqlc, app.TableApp, false, []string{app.SqlCol_App_UserId})
		henge.CreateIndexSql(sqlc, session.TableSession, false, []string{session.SqlCol_Session_IdSource})
		henge.CreateIndexSql(sqlc, session.TableSession, false, []string{session.SqlCol_Session_AppId})
		henge.CreateIndexSql(sqlc, session.TableSession, false, []string{session.SqlCol_Session_Expiry})

		userDao = _createUserDaoSql(sqlc)
		appDao = _createAppDaoSql(sqlc)
		sessionDao = _createSessionDaoSql(sqlc)
	}
	if dync != nil {
		henge.InitDynamodbTable(dync, user.TableUser, 2, 1)
		henge.InitDynamodbTable(dync, app.TableApp, 2, 1)
		henge.InitDynamodbTable(dync, session.TableSession, 4, 2)

		userDao = _createUserDaoAwsDynamodb(dync)
		appDao = _createAppDaoAwsDynamodb(dync)
		sessionDao = _createSessionDaoAwsDynamodb(dync)
	}

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

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
	"main/src/gvabe/bo"
	"main/src/gvabe/bo/app"
	"main/src/gvabe/bo/session"
	"main/src/gvabe/bo/user"
	"main/src/utils"
)

var (
	dbTypeCosmosDb = []string{"cosmosdb", "azurecosmosdb", "azure_cosmosdb", "azure-cosmosdb"}
	dbTypeDynamoDb = []string{"dynamodb", "awsdynamodb", "aws_dynamodb", "aws-dynamodb"}
	dbTypeMssql    = []string{"mssql"}
	dbTypeMysql    = []string{"mysql"}
	dbTypeMongoDb  = []string{"mongo", "mongodb"}
	dbTypeOracle   = []string{"oracle"}
	dbTypePgsql    = []string{"pg", "pgsql", "postgres", "postgresql"}
	dbTypeSqlite   = []string{"sqlite", "sqlite3"}
)

func _createSqlConnect(dbtype string) *prom.SqlConnect {
	var poolOpts *prom.SqlPoolOptions = nil
	timezone := goapi.AppConfig.GetString("timezone")
	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	var dsn, driver string
	var dbflavor = prom.FlavorDefault
	switch {
	case utils.InSlideStr(dbtype, dbTypeSqlite):
		driver = "sqlite3"
		dbflavor = prom.FlavorSqlite
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		os.MkdirAll(dir, 0711)
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		dsn = dir + "/" + dbname + ".db"
	case utils.InSlideStr(dbtype, dbTypePgsql):
		driver = "pgx"
		dbflavor = prom.FlavorPgSql
		dsn = goapi.AppConfig.GetString("gvabe.db.pgsql.url")
	case utils.InSlideStr(dbtype, dbTypeCosmosDb):
		daoMultitenant = goapi.AppConfig.GetBoolean("gvabe.db.cosmosdb.multitenant", true)
		driver = "gocosmos"
		dbflavor = prom.FlavorCosmosDb
		dsn = goapi.AppConfig.GetString("gvabe.db.cosmosdb.url")
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
	if utils.InSlideStr(dbtype, dbTypeDynamoDb) {
		daoMultitenant = goapi.AppConfig.GetBoolean("gvabe.db.dynamodb.multitenant", true)
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

func _createMongoConnect(dbtype string) *prom.MongoConnect {
	if utils.InSlideStr(dbtype, dbTypeMongoDb) {
		db := goapi.AppConfig.GetString("gvabe.db.mongodb.db")
		url := goapi.AppConfig.GetString("gvabe.db.mongodb.url")
		mc, err := prom.NewMongoConnect(url, db, 2345)
		if err != nil {
			panic(err)
		}
		return mc
	}
	return nil
}

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype)
	dync := _createAwsDynamodbConnect(dbtype)
	mc := _createMongoConnect(dbtype)
	if sqlc == nil && dync == nil && mc == nil {
		panic(fmt.Sprintf("unsupported database type: %s", dbtype))
	}
	switch {
	case utils.InSlideStr(dbtype, dbTypeSqlite):
		// SQLite, for non-production only!
		henge.InitSqliteTable(sqlc, user.TableUser, nil)
		henge.InitSqliteTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "VARCHAR(32)"})
		henge.InitSqliteTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "VARCHAR(32)",
			session.SqlCol_Session_AppId:       "VARCHAR(32)",
			session.SqlCol_Session_UserId:      "VARCHAR(32)",
			session.SqlCol_Session_SessionType: "VARCHAR(32)",
			session.SqlCol_Session_Expiry:      "TIMESTAMP",
		})
	case utils.InSlideStr(dbtype, dbTypeMssql):
		// MSSQL
		henge.InitMssqlTable(sqlc, user.TableUser, nil)
		henge.InitMssqlTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "NVARCHAR(32)"})
		henge.InitMssqlTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "NVARCHAR(32)",
			session.SqlCol_Session_AppId:       "NVARCHAR(32)",
			session.SqlCol_Session_UserId:      "NVARCHAR(32)",
			session.SqlCol_Session_SessionType: "NVARCHAR(32)",
			session.SqlCol_Session_Expiry:      "DATETIMEOFFSET",
		})
	case utils.InSlideStr(dbtype, dbTypeMysql):
		// MySQL
		henge.InitMysqlTable(sqlc, user.TableUser, nil)
		henge.InitMysqlTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "VARCHAR(32)"})
		henge.InitMysqlTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "VARCHAR(32)",
			session.SqlCol_Session_AppId:       "VARCHAR(32)",
			session.SqlCol_Session_UserId:      "VARCHAR(32)",
			session.SqlCol_Session_SessionType: "VARCHAR(32)",
			session.SqlCol_Session_Expiry:      "DATETIME",
		})
	case utils.InSlideStr(dbtype, dbTypeOracle):
		henge.InitOracleTable(sqlc, user.TableUser, nil)
		henge.InitOracleTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "NVARCHAR2(32)"})
		henge.InitOracleTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "NVARCHAR2(32)",
			session.SqlCol_Session_AppId:       "NVARCHAR2(32)",
			session.SqlCol_Session_UserId:      "NVARCHAR2(32)",
			session.SqlCol_Session_SessionType: "NVARCHAR2(32)",
			session.SqlCol_Session_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
	case utils.InSlideStr(dbtype, dbTypePgsql):
		// PostgreSQL
		henge.InitPgsqlTable(sqlc, user.TableUser, nil)
		henge.InitPgsqlTable(sqlc, app.TableApp, map[string]string{app.SqlCol_App_UserId: "VARCHAR(32)"})
		henge.InitPgsqlTable(sqlc, session.TableSession, map[string]string{
			session.SqlCol_Session_IdSource:    "VARCHAR(32)",
			session.SqlCol_Session_AppId:       "VARCHAR(32)",
			session.SqlCol_Session_UserId:      "VARCHAR(32)",
			session.SqlCol_Session_SessionType: "VARCHAR(32)",
			session.SqlCol_Session_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
	}

	if dync != nil {
		// AWS DynamoDB
		spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: false}
		if daoMultitenant {
			spec.MainTablePkPrefix = bo.DynamodbMultitenantPkName
			spec.MainTableCustomAttrs = []prom.AwsDynamodbNameAndType{{Name: bo.DynamodbMultitenantPkName, Type: prom.AwsAttrTypeString}}
			henge.InitDynamodbTables(dync, bo.DynamodbMultitenantTableName, spec)

			appDao = app.NewAppDaoMultitenantAwsDynamodb(dync, bo.DynamodbMultitenantTableName)
			sessionDao = session.NewSessionDaoMultitenantAwsDynamodb(dync, bo.DynamodbMultitenantTableName)
			userDao = user.NewUserDaoMultitenantAwsDynamodb(dync, bo.DynamodbMultitenantTableName)
		} else {
			henge.InitDynamodbTables(dync, app.TableApp, spec)
			henge.InitDynamodbTables(dync, session.TableSession, spec)
			henge.InitDynamodbTables(dync, user.TableUser, spec)

			appDao = app.NewAppDaoAwsDynamodb(dync, app.TableApp)
			sessionDao = session.NewSessionDaoAwsDynamodb(dync, session.TableSession)
			userDao = user.NewUserDaoAwsDynamodb(dync, user.TableUser)
		}
	} else if mc != nil {
		// MongoDB
		henge.InitMongoCollection(mc, app.TableApp)
		henge.InitMongoCollection(mc, session.TableSession)
		henge.InitMongoCollection(mc, user.TableUser)

		mc.CreateCollectionIndexes(app.TableApp, []interface{}{
			map[string]interface{}{
				"key":  map[string]interface{}{app.FieldApp_OwnerId: 1},
				"name": "idx_ownerid",
			},
		})
		mc.CreateCollectionIndexes(session.TableSession, []interface{}{
			map[string]interface{}{
				"key":  map[string]interface{}{session.FieldSession_IdSource: 1},
				"name": "idx_idsource",
			},
			map[string]interface{}{
				"key":  map[string]interface{}{session.FieldSession_AppId: 1},
				"name": "idx_appid",
			},
			map[string]interface{}{
				"key":  map[string]interface{}{session.FieldSession_Expiry: 1},
				"name": "idx_expiry",
			},
		})

		appDao = app.NewAppDaoMongo(mc, app.TableApp)
		sessionDao = session.NewSessionDaoMongo(mc, session.TableSession)
		userDao = user.NewUserDaoMongo(mc, user.TableUser)
	} else if sqlc != nil && utils.InSlideStr(dbtype, dbTypeCosmosDb) {
		// Azure Cosmos DB
		spec := &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbPkName}
		if daoMultitenant {
			spec.Pk = bo.CosmosdbMultitenantPkName
			henge.InitCosmosdbCollection(sqlc, bo.CosmosdbMultitenantTableName, spec)

			appDao = app.NewAppDaoMultitenantCosmosdb(sqlc, bo.CosmosdbMultitenantTableName)
			sessionDao = session.NewSessionDaoMultitenantCosmosdb(sqlc, bo.CosmosdbMultitenantTableName)
			userDao = user.NewUserDaoMultitenantCosmosdb(sqlc, bo.CosmosdbMultitenantTableName)
		} else {
			henge.InitCosmosdbCollection(sqlc, app.TableApp, spec)
			henge.InitCosmosdbCollection(sqlc, session.TableSession, spec)
			henge.InitCosmosdbCollection(sqlc, user.TableUser, spec)

			appDao = app.NewAppDaoCosmosdb(sqlc, app.TableApp)
			sessionDao = session.NewSessionDaoCosmosdb(sqlc, session.TableSession)
			userDao = user.NewUserDaoCosmosdb(sqlc, user.TableUser)
		}
	} else if sqlc != nil {
		// other RDBMS
		henge.CreateIndexSql(sqlc, app.TableApp, false, []string{app.SqlCol_App_UserId})
		henge.CreateIndexSql(sqlc, session.TableSession, false, []string{session.SqlCol_Session_IdSource})
		henge.CreateIndexSql(sqlc, session.TableSession, false, []string{session.SqlCol_Session_AppId})
		henge.CreateIndexSql(sqlc, session.TableSession, false, []string{session.SqlCol_Session_Expiry})

		appDao = app.NewAppDaoSql(sqlc, app.TableApp)
		sessionDao = session.NewSessionDaoSql(sqlc, session.TableSession)
		userDao = user.NewUserDaoSql(sqlc, user.TableUser)
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

package app

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func newSqlConnectSqlite(driver, url, timezone string, timeoutMs int, poolOptions *prom.SqlPoolOptions) (*prom.SqlConnect, error) {
	os.Remove(url)
	sqlc, err := prom.NewSqlConnectWithFlavor(driver, url, timeoutMs, poolOptions, prom.FlavorSqlite)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}

func newSqlConnectMssql(driver, url, timezone string, timeoutMs int, poolOptions *prom.SqlPoolOptions) (*prom.SqlConnect, error) {
	sqlc, err := prom.NewSqlConnectWithFlavor(driver, url, timeoutMs, poolOptions, prom.FlavorMsSql)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}

func newSqlConnectMysql(driver, url, timezone string, timeoutMs int, poolOptions *prom.SqlPoolOptions) (*prom.SqlConnect, error) {
	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	sqlc, err := prom.NewSqlConnectWithFlavor(driver, url, timeoutMs, poolOptions, prom.FlavorMySql)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}

func newSqlConnectOracle(driver, url, timezone string, timeoutMs int, poolOptions *prom.SqlPoolOptions) (*prom.SqlConnect, error) {
	sqlc, err := prom.NewSqlConnectWithFlavor(driver, url, timeoutMs, poolOptions, prom.FlavorOracle)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}

func newSqlConnectPgsql(driver, url, timezone string, timeoutMs int, poolOptions *prom.SqlPoolOptions) (*prom.SqlConnect, error) {
	sqlc, err := prom.NewSqlConnectWithFlavor(driver, url, timeoutMs, poolOptions, prom.FlavorPgSql)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}

const (
	envSqliteDriver = "SQLITE_DRIVER"
	envSqliteUrl    = "SQLITE_URL"
	envMssqlDriver  = "MSSQL_DRIVER"
	envMssqlUrl     = "MSSQL_URL"
	envMysqlDriver  = "MYSQL_DRIVER"
	envMysqlUrl     = "MYSQL_URL"
	envOracleDriver = "ORACLE_DRIVER"
	envOracleUrl    = "ORACLE_URL"
	envPgsqlDriver  = "PGSQL_DRIVER"
	envPgsqlUrl     = "PGSQL_URL"
	tableNameSql    = "exter_test_app"
	timezoneSql     = "Asia/Ho_Chi_Minh"
)

type sqlDriverAndUrl struct {
	driver, url string
}

func newSqlDriverAndUrl(driver, url string) sqlDriverAndUrl {
	return sqlDriverAndUrl{driver: strings.Trim(driver, `"`), url: strings.Trim(url, `"`)}
}

func sqlGetUrlFromEnv() map[string]sqlDriverAndUrl {
	urlMap := make(map[string]sqlDriverAndUrl)
	if os.Getenv(envSqliteDriver) != "" && os.Getenv(envSqliteUrl) != "" {
		urlMap["sqlite"] = newSqlDriverAndUrl(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl))
	}
	if os.Getenv(envMssqlDriver) != "" && os.Getenv(envMssqlUrl) != "" {
		urlMap["mssql"] = newSqlDriverAndUrl(os.Getenv(envMssqlDriver), os.Getenv(envMssqlUrl))
	}
	if os.Getenv(envMysqlDriver) != "" && os.Getenv(envMysqlUrl) != "" {
		urlMap["mysql"] = newSqlDriverAndUrl(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl))
	}
	if os.Getenv(envOracleDriver) != "" && os.Getenv(envOracleUrl) != "" {
		urlMap["oracle"] = newSqlDriverAndUrl(os.Getenv(envOracleDriver), os.Getenv(envOracleUrl))
	}
	if os.Getenv(envPgsqlDriver) != "" && os.Getenv(envPgsqlUrl) != "" {
		urlMap["pgsql"] = newSqlDriverAndUrl(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl))
	}
	return urlMap
}

func TestNewAppDaoSql(t *testing.T) {
	name := "TestNewAppDaoSql"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for k, info := range urlMap {
		var sqlc *prom.SqlConnect
		var err error
		switch k {
		case "sqlite", "sqlite3":
			sqlc, err = newSqlConnectSqlite(info.driver, info.url, timezoneSql, 10000, nil)
		case "mssql":
			sqlc, err = newSqlConnectMssql(info.driver, info.url, timezoneSql, 10000, nil)
		case "mysql":
			sqlc, err = newSqlConnectMysql(info.driver, info.url, timezoneSql, 10000, nil)
		case "oracle":
			sqlc, err = newSqlConnectOracle(info.driver, info.url, timezoneSql, 10000, nil)
		case "pgsql":
			sqlc, err = newSqlConnectPgsql(info.driver, info.url, timezoneSql, 10000, nil)
		default:
			t.Fatalf("%s failed: unknown database type [%s]", name, k)
		}
		if err != nil {
			t.Fatalf("%s failed: error [%e]", name+"/"+k, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+k)
		}
		appDao := NewAppDaoSql(sqlc, tableNameSql)
		if appDao == nil {
			t.Fatalf("%s failed: nil", name)
		}
	}
}

func _initAppDaoSql(t *testing.T, testName string, sqlc *prom.SqlConnect) AppDao {
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", tableNameSql))
	switch sqlc.GetDbFlavor() {
	case prom.FlavorPgSql:
		henge.InitPgsqlTable(sqlc, tableNameSql, map[string]string{SqlCol_App_UserId: "VARCHAR(32)"})
	case prom.FlavorMsSql:
		henge.InitMssqlTable(sqlc, tableNameSql, map[string]string{SqlCol_App_UserId: "NVARCHAR(32)"})
	case prom.FlavorMySql:
		henge.InitMysqlTable(sqlc, tableNameSql, map[string]string{SqlCol_App_UserId: "VARCHAR(32)"})
	case prom.FlavorOracle:
		henge.InitOracleTable(sqlc, tableNameSql, map[string]string{SqlCol_App_UserId: "NVARCHAR2(32)"})
	case prom.FlavorSqlite:
		henge.InitSqliteTable(sqlc, tableNameSql, map[string]string{SqlCol_App_UserId: "VARCHAR(32)"})
	default:
		t.Fatalf("%s failed: unknown database type %#v", testName, sqlc.GetDbFlavor())
	}
	return NewAppDaoSql(sqlc, tableNameSql)
}

func TestAppDaosql_Create(t *testing.T) {
	name := "TestAppDaosql_Create"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for k, info := range urlMap {
		var sqlc *prom.SqlConnect
		var err error
		switch k {
		case "sqlite", "sqlite3":
			sqlc, err = newSqlConnectSqlite(info.driver, info.url, timezoneSql, 10000, nil)
		case "mssql":
			sqlc, err = newSqlConnectMssql(info.driver, info.url, timezoneSql, 10000, nil)
		case "mysql":
			sqlc, err = newSqlConnectMysql(info.driver, info.url, timezoneSql, 10000, nil)
		case "oracle":
			sqlc, err = newSqlConnectOracle(info.driver, info.url, timezoneSql, 10000, nil)
		case "pgsql":
			sqlc, err = newSqlConnectPgsql(info.driver, info.url, timezoneSql, 10000, nil)
		default:
			t.Fatalf("%s failed: unknown database type [%s]", name, k)
		}
		appDao := _initAppDaoSql(t, name, sqlc)
		app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
		ok, err := appDao.Create(app)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}
	}
}

// func TestAppDaoAwsDynamodb_Get(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_Get"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
// 	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
// 	if app, err := appDao.Get("not_found"); err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app != nil {
// 		t.Fatalf("%s failed: app %s should not exist", name, "not_found")
// 	}
//
// 	if app, err := appDao.Get("exter"); err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app == nil {
// 		t.Fatalf("%s failed: nil", name)
// 	} else {
// 		if v := app.GetId(); v != "exter" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
// 		}
// 		if v := app.GetTagVersion(); v != 1357 {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
// 		}
// 		if v := app.GetOwnerId(); v != "btnguyen2k" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
// 		}
// 		if v := app.GetAttrsPublic().Description; v != "System application (do not delete)" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "System application (do not delete)", v)
// 		}
// 	}
// }
//
// func TestAppDaoAwsDynamodb_Delete(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_Delete"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
//
// 	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
// 	app, err := appDao.Get("exter")
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app == nil {
// 		t.Fatalf("%s failed: nil", name)
// 	}
//
// 	ok, err := appDao.Delete(app)
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if !ok {
// 		t.Fatalf("%s failed: cannot delete app [%s]", name, app.GetId())
// 	}
//
// 	app, err = appDao.Get("exter")
// 	if app, err := appDao.Get("exter"); err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app != nil {
// 		t.Fatalf("%s failed: app %s should not exist", name, "exter")
// 	}
// }
//
// func TestAppDaoAwsDynamodb_Update(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_Update"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
//
// 	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
// 	appDao.Create(app)
//
// 	app.SetOwnerId("nbthanh")
// 	app.SetTagVersion(2468)
// 	app.attrsPublic.Description = "App description"
// 	ok, err := appDao.Update(app)
// 	if err != nil || !ok {
// 		t.Fatalf("%s failed: %#v / %s", name, ok, err)
// 	}
//
// 	if app, err := appDao.Get("exter"); err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app == nil {
// 		t.Fatalf("%s failed: nil", name)
// 	} else {
// 		if v := app.GetId(); v != "exter" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
// 		}
// 		if v := app.GetTagVersion(); v != 2468 {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
// 		}
// 		if v := app.GetOwnerId(); v != "nbthanh" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
// 		}
// 		if v := app.GetAttrsPublic().Description; v != "App description" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "App description", v)
// 		}
// 	}
// }
//
// func TestAppDaoAwsDynamodb_GetUserApps(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_GetUserApps"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
//
// 	for i := 0; i < 10; i++ {
// 		app := NewApp(uint64(i), strconv.Itoa(i), strconv.Itoa(i%3), "App #"+strconv.Itoa(i))
// 		appDao.Create(app)
// 	}
//
// 	u := user.NewUser(123, "2")
// 	appList, err := appDao.GetUserApps(u)
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	}
// 	if len(appList) != 3 {
// 		t.Fatalf("%s failed: expected %#v apps but received %#v", name, 3, len(appList))
// 	}
// 	for _, app := range appList {
// 		if app.GetOwnerId() != "2" {
// 			t.Fatalf("%s failed: app %#v does not belong to user %#v", name, app.GetId(), "2")
// 		}
// 	}
// }

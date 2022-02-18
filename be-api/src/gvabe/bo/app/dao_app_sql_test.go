package app

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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

var (
	testSqlDbtype   string
	testSqlConnInfo sqlDriverAndUrl
)

func _createSqlConnect(t *testing.T, testName string, dbtype string, connInfo sqlDriverAndUrl) *prom.SqlConnect {
	var sqlc *prom.SqlConnect
	var err error
	switch dbtype {
	case "sqlite", "sqlite3":
		sqlc, err = newSqlConnectSqlite(connInfo.driver, connInfo.url, timezoneSql, 10000, nil)
	case "mssql":
		sqlc, err = newSqlConnectMssql(connInfo.driver, connInfo.url, timezoneSql, 10000, nil)
	case "mysql":
		sqlc, err = newSqlConnectMysql(connInfo.driver, connInfo.url, timezoneSql, 10000, nil)
	case "oracle":
		sqlc, err = newSqlConnectOracle(connInfo.driver, connInfo.url, timezoneSql, 10000, nil)
	case "pgsql":
		sqlc, err = newSqlConnectPgsql(connInfo.driver, connInfo.url, timezoneSql, 10000, nil)
	default:
		t.Fatalf("%s failed: unknown database type [%s]", testName, dbtype)
	}
	if err != nil {
		t.Fatalf("%s failed: error [%e]", testName+"/"+dbtype, err)
	} else if sqlc == nil {
		t.Fatalf("%s failed: nil", testName+"/"+dbtype)
	}
	return sqlc
}

var setupTestSql = func(t *testing.T, testName string) {
	testSqlc = _createSqlConnect(t, testName, testSqlDbtype, testSqlConnInfo)
	testSqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", tableNameSql))
	err := InitAppTableSql(testSqlc, tableNameSql)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestSql = func(t *testing.T, testName string) {
	if testSqlc != nil {
		defer func() {
			defer func() { testSqlc = nil }()
			testSqlc.Close()
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewAppDaoSql(t *testing.T) {
	testName := "TestNewAppDaoSql"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			appDao := NewAppDaoSql(testSqlc, tableNameSql)
			if appDao == nil {
				t.Fatalf("%s failed: nil", testName+"/"+testSqlDbtype)
			}
		})
	}
}

func TestAppDaosql_Create(t *testing.T) {
	testName := "TestAppDaosql_Create"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			appDao := NewAppDaoSql(testSqlc, tableNameSql)
			doTestAppDao_Create(t, testName, appDao)
		})
	}
}

func TestAppDaoSql_Get(t *testing.T) {
	testName := "TestAppDaoSql_Get"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			appDao := NewAppDaoSql(testSqlc, tableNameSql)
			doTestAppDao_Get(t, testName, appDao)
		})
	}
}

func TestAppDaoSql_Delete(t *testing.T) {
	testName := "TestAppDaoSql_Delete"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			appDao := NewAppDaoSql(testSqlc, tableNameSql)
			doTestAppDao_Delete(t, testName, appDao)
		})
	}
}

func TestAppDaoSql_Update(t *testing.T) {
	testName := "TestAppDaoSql_Update"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			appDao := NewAppDaoSql(testSqlc, tableNameSql)
			doTestAppDao_Update(t, testName, appDao)
		})
	}
}

func TestAppDaoSql_GetUserApps(t *testing.T) {
	testName := "TestAppDaoSql_GetUserApps"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			appDao := NewAppDaoSql(testSqlc, tableNameSql)
			doTestAppDao_GetUserApps(t, testName, appDao)
		})
	}
}

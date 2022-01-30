package user

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
	tableNameSql    = "exter_test_user"
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
	err := InitUserTableSql(testSqlc, tableNameSql)
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

func TestNewUserDaoSql(t *testing.T) {
	testName := "TestNewUserDaoSql"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			userDao := NewUserDaoSql(testSqlc, tableNameSql)
			if userDao == nil {
				t.Fatalf("%s failed: nil", testName+"/"+testSqlDbtype)
			}
		})
	}
}

func TestUserDaoSql_Create(t *testing.T) {
	testName := "TestUserDaoSql_Create"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			userDao := NewUserDaoSql(testSqlc, tableNameSql)
			u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
			ok, err := userDao.Create(u)
			if err != nil || !ok {
				t.Fatalf("%s failed: %#v / %s", testName, ok, err)
			}
		})
	}
}

func TestUserDaoSql_Get(t *testing.T) {
	testName := "TestUserDaoSql_Get"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			userDao := NewUserDaoSql(testSqlc, tableNameSql)
			u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
			ok, err := userDao.Create(u)
			if err != nil || !ok {
				t.Fatalf("%s failed: %#v / %s", testName, ok, err)
			}
			if u, err := userDao.Get("not_found"); err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if u != nil {
				t.Fatalf("%s failed: user %s should not exist", testName, "not_found")
			}

			if u, err := userDao.Get("btnguyen2k"); err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if u == nil {
				t.Fatalf("%s failed: nil", testName)
			} else {
				if v := u.GetId(); v != "btnguyen2k" {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "btnguyen2k", v)
				}
				if v := u.GetTagVersion(); v != 1357 {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, 1357, v)
				}
				if v := u.GetDisplayName(); v != "Thanh Nguyen" {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "Thanh Nguyen", v)
				}
				if v := u.GetAesKey(); v != "aeskey" {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "aeskey", v)
				}
			}
		})
	}
}

func TestUserDaoSql_Delete(t *testing.T) {
	testName := "TestUserDaoSql_Delete"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			userDao := NewUserDaoSql(testSqlc, tableNameSql)
			u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
			ok, err := userDao.Create(u)
			if err != nil || !ok {
				t.Fatalf("%s failed: %#v / %s", testName, ok, err)
			}

			ok, err = userDao.Delete(u)
			if err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if !ok {
				t.Fatalf("%s failed: cannot delete user [%s]", testName, u.GetId())
			}

			u, err = userDao.Get("btnguyen2k")
			if app, err := userDao.Get("exter"); err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if app != nil {
				t.Fatalf("%s failed: user %s should not exist", testName, "userDao")
			}
		})
	}
}

func TestUserDaoUser_Update(t *testing.T) {
	testName := "TestUserDaoUser_Update"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", testName)
	}
	for testSqlDbtype, testSqlConnInfo = range urlMap {
		t.Run(testSqlDbtype, func(t *testing.T) {
			teardownTest := setupTest(t, testName, setupTestSql, teardownTestSql)
			defer teardownTest(t)
			userDao := NewUserDaoSql(testSqlc, tableNameSql)

			u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
			userDao.Create(u)

			u.SetDisplayName("nbthanh")
			u.SetAesKey("newaeskey")
			ok, err := userDao.Update(u)
			if err != nil || !ok {
				t.Fatalf("%s failed: %#v / %s", testName, ok, err)
			}

			if u, err := userDao.Get("btnguyen2k"); err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if u == nil {
				t.Fatalf("%s failed: nil", testName)
			} else {
				if v := u.GetId(); v != "btnguyen2k" {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "btnguyen2k", v)
				}
				if v := u.GetTagVersion(); v != 1357 {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, 1357, v)
				}
				if v := u.GetDisplayName(); v != "nbthanh" {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "nbthanh", v)
				}
				if v := u.GetAesKey(); v != "newaeskey" {
					t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "newaeskey", v)
				}
			}
		})
	}
}

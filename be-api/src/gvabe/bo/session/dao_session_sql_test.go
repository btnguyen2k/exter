package session

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
	tableNameSql    = "exter_test_sesion"
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

func TestNewSessionDaoSql(t *testing.T) {
	name := "TestNewSessionDaoSql"
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
		appDao := NewSessionDaoSql(sqlc, tableNameSql)
		if appDao == nil {
			t.Fatalf("%s failed: nil", name)
		}
	}
}

func _initSessionDaoSql(t *testing.T, testName string, sqlc *prom.SqlConnect) SessionDao {
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", tableNameSql))
	switch sqlc.GetDbFlavor() {
	case prom.FlavorPgSql:
		henge.InitPgsqlTable(sqlc, tableNameSql, map[string]string{
			SqlCol_Session_IdSource:    "VARCHAR(32)",
			SqlCol_Session_AppId:       "VARCHAR(32)",
			SqlCol_Session_UserId:      "VARCHAR(32)",
			SqlCol_Session_SessionType: "VARCHAR(32)",
			SqlCol_Session_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
	case prom.FlavorMsSql:
		henge.InitMssqlTable(sqlc, tableNameSql, map[string]string{
			SqlCol_Session_IdSource:    "NVARCHAR(32)",
			SqlCol_Session_AppId:       "NVARCHAR(32)",
			SqlCol_Session_UserId:      "NVARCHAR(32)",
			SqlCol_Session_SessionType: "NVARCHAR(32)",
			SqlCol_Session_Expiry:      "DATETIMEOFFSET",
		})
	case prom.FlavorMySql:
		henge.InitMysqlTable(sqlc, tableNameSql, map[string]string{
			SqlCol_Session_IdSource:    "VARCHAR(32)",
			SqlCol_Session_AppId:       "VARCHAR(32)",
			SqlCol_Session_UserId:      "VARCHAR(32)",
			SqlCol_Session_SessionType: "VARCHAR(32)",
			SqlCol_Session_Expiry:      "TIMESTAMP",
		})
	case prom.FlavorOracle:
		henge.InitOracleTable(sqlc, tableNameSql, map[string]string{
			SqlCol_Session_IdSource:    "NVARCHAR2(32)",
			SqlCol_Session_AppId:       "NVARCHAR2(32)",
			SqlCol_Session_UserId:      "NVARCHAR2(32)",
			SqlCol_Session_SessionType: "NVARCHAR2(32)",
			SqlCol_Session_Expiry:      "TIMESTAMP WITH TIME ZONE",
		})
	case prom.FlavorSqlite:
		henge.InitSqliteTable(sqlc, tableNameSql, map[string]string{
			SqlCol_Session_IdSource:    "VARCHAR(32)",
			SqlCol_Session_AppId:       "VARCHAR(32)",
			SqlCol_Session_UserId:      "VARCHAR(32)",
			SqlCol_Session_SessionType: "VARCHAR(32)",
			SqlCol_Session_Expiry:      "TIMESTAMP",
		})
	default:
		t.Fatalf("%s failed: unknown database type %#v", testName, sqlc.GetDbFlavor())
	}
	return NewSessionDaoSql(sqlc, tableNameSql)
}

func TestSessionDaoSql_Save(t *testing.T) {
	name := "TestSessionDaoSql_Save"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for k, info := range urlMap {
		var sqlc *prom.SqlConnect
		switch k {
		case "sqlite", "sqlite3":
			sqlc, _ = newSqlConnectSqlite(info.driver, info.url, timezoneSql, 10000, nil)
		case "mssql":
			sqlc, _ = newSqlConnectMssql(info.driver, info.url, timezoneSql, 10000, nil)
		case "mysql":
			sqlc, _ = newSqlConnectMysql(info.driver, info.url, timezoneSql, 10000, nil)
		case "oracle":
			sqlc, _ = newSqlConnectOracle(info.driver, info.url, timezoneSql, 10000, nil)
		case "pgsql":
			sqlc, _ = newSqlConnectPgsql(info.driver, info.url, timezoneSql, 10000, nil)
		default:
			t.Fatalf("%s failed: unknown database type [%s]", name, k)
		}
		sessDao := _initSessionDaoSql(t, name, sqlc)
		expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
		sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
		ok, err := sessDao.Save(sess)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}
	}
}

func TestSessionDaoSql_Get(t *testing.T) {
	name := "TestSessionDaoSql_Get"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for k, info := range urlMap {
		var sqlc *prom.SqlConnect
		switch k {
		case "sqlite", "sqlite3":
			sqlc, _ = newSqlConnectSqlite(info.driver, info.url, timezoneSql, 10000, nil)
		case "mssql":
			sqlc, _ = newSqlConnectMssql(info.driver, info.url, timezoneSql, 10000, nil)
		case "mysql":
			sqlc, _ = newSqlConnectMysql(info.driver, info.url, timezoneSql, 10000, nil)
		case "oracle":
			sqlc, _ = newSqlConnectOracle(info.driver, info.url, timezoneSql, 10000, nil)
		case "pgsql":
			sqlc, _ = newSqlConnectPgsql(info.driver, info.url, timezoneSql, 10000, nil)
		default:
			t.Fatalf("%s failed: unknown database type [%s]", name, k)
		}
		sessDao := _initSessionDaoSql(t, name, sqlc)
		expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
		sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
		ok, err := sessDao.Save(sess)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}

		if sess, err := sessDao.Get("not_found"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if sess != nil {
			t.Fatalf("%s failed: session %s should not exist", name, "not_found")
		}

		if sess, err := sessDao.Get("1"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if sess == nil {
			t.Fatalf("%s failed: nil", name)
		} else {
			if v := sess.GetId(); v != "1" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "1", v)
			}
			if v := sess.GetTagVersion(); v != 1357 {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
			}
			if v := sess.GetSessionType(); v != "login" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "login", v)
			}
			if v := sess.GetIdSource(); v != "local" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "local", v)
			}
			if v := sess.GetAppId(); v != "exter" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
			}
			if v := sess.GetUserId(); v != "btnguyen2k" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
			}
			if v := sess.GetSessionData(); v != "session-data" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "session-data", v)
			}
			if v := sess.GetExpiry(); v.Unix() != expiry.Unix() {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, expiry, v)
			}
		}
	}
}

func TestSessionDaoSql_Delete(t *testing.T) {
	name := "TestSessionDaoSql_Delete"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for k, info := range urlMap {
		var sqlc *prom.SqlConnect
		switch k {
		case "sqlite", "sqlite3":
			sqlc, _ = newSqlConnectSqlite(info.driver, info.url, timezoneSql, 10000, nil)
		case "mssql":
			sqlc, _ = newSqlConnectMssql(info.driver, info.url, timezoneSql, 10000, nil)
		case "mysql":
			sqlc, _ = newSqlConnectMysql(info.driver, info.url, timezoneSql, 10000, nil)
		case "oracle":
			sqlc, _ = newSqlConnectOracle(info.driver, info.url, timezoneSql, 10000, nil)
		case "pgsql":
			sqlc, _ = newSqlConnectPgsql(info.driver, info.url, timezoneSql, 10000, nil)
		default:
			t.Fatalf("%s failed: unknown database type [%s]", name, k)
		}
		sessDao := _initSessionDaoSql(t, name, sqlc)
		expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
		sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
		ok, err := sessDao.Save(sess)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}
		if sess, err := sessDao.Get("1"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if sess == nil {
			t.Fatalf("%s failed: nill", name)
		}

		ok, err = sessDao.Delete(sess)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}

		if sess, err := sessDao.Get("1"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if sess != nil {
			t.Fatalf("%s failed: session %s should not exist", name, "not_found")
		}
	}
}

func TestSessionDaoSql_Update(t *testing.T) {
	name := "TestSessionDaoSql_Update"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for k, info := range urlMap {
		var sqlc *prom.SqlConnect
		switch k {
		case "sqlite", "sqlite3":
			sqlc, _ = newSqlConnectSqlite(info.driver, info.url, timezoneSql, 10000, nil)
		case "mssql":
			sqlc, _ = newSqlConnectMssql(info.driver, info.url, timezoneSql, 10000, nil)
		case "mysql":
			sqlc, _ = newSqlConnectMysql(info.driver, info.url, timezoneSql, 10000, nil)
		case "oracle":
			sqlc, _ = newSqlConnectOracle(info.driver, info.url, timezoneSql, 10000, nil)
		case "pgsql":
			sqlc, _ = newSqlConnectPgsql(info.driver, info.url, timezoneSql, 10000, nil)
		default:
			t.Fatalf("%s failed: unknown database type [%s]", name, k)
		}
		sessDao := _initSessionDaoSql(t, name, sqlc)
		expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)

		sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
		ok, err := sessDao.Save(sess)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}

		sess.SetTagVersion(2468)
		sess.SetSessionType("pre-login")
		sess.SetIdSource("external")
		sess.SetAppId("myapp")
		sess.SetUserId("nbthanh")
		sess.SetSessionData("data")
		sess.SetExpiry(expiry.Add(1 * time.Hour))
		ok, err = sessDao.Save(sess)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}

		if sess, err := sessDao.Get("1"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if sess == nil {
			t.Fatalf("%s failed: nil", name)
		} else {
			if v := sess.GetId(); v != "1" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "1", v)
			}
			if v := sess.GetTagVersion(); v != 2468 {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
			}
			if v := sess.GetSessionType(); v != "pre-login" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "pre-login", v)
			}
			if v := sess.GetIdSource(); v != "external" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "external", v)
			}
			if v := sess.GetAppId(); v != "myapp" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "myapp", v)
			}
			if v := sess.GetUserId(); v != "nbthanh" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
			}
			if v := sess.GetSessionData(); v != "data" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "data", v)
			}
			if v := sess.GetExpiry(); v.Unix() != expiry.Add(1*time.Hour).Unix() {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, expiry.Add(1*time.Hour), v)
			}
		}
	}
}

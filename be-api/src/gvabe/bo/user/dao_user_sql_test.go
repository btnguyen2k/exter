package user

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

func TestNewUserDaoSql(t *testing.T) {
	name := "TestNewUserDaoSql"
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
		userDao := NewUserDaoSql(sqlc, tableNameSql)
		if userDao == nil {
			t.Fatalf("%s failed: nil", name+"/"+k)
		}
	}
}

func _initUserDaoSql(t *testing.T, testName string, sqlc *prom.SqlConnect) UserDao {
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", tableNameSql))
	switch sqlc.GetDbFlavor() {
	case prom.FlavorPgSql:
		henge.InitPgsqlTable(sqlc, tableNameSql, nil)
	case prom.FlavorMsSql:
		henge.InitMssqlTable(sqlc, tableNameSql, nil)
	case prom.FlavorMySql:
		henge.InitMysqlTable(sqlc, tableNameSql, nil)
	case prom.FlavorOracle:
		henge.InitOracleTable(sqlc, tableNameSql, nil)
	case prom.FlavorSqlite:
		henge.InitSqliteTable(sqlc, tableNameSql, nil)
	default:
		t.Fatalf("%s failed: unknown database type %#v", testName, sqlc.GetDbFlavor())
	}
	return NewUserDaoSql(sqlc, tableNameSql)
}

func TestUserDaoSql_Create(t *testing.T) {
	name := "TestUserDaoSql_Create"
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
		userDao := _initUserDaoSql(t, name, sqlc)
		u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
		ok, err := userDao.Create(u)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}
	}
}

func TestUserDaoSql_Get(t *testing.T) {
	name := "TestUserDaoSql_Get"
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
		userDao := _initUserDaoSql(t, name, sqlc)
		u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
		ok, err := userDao.Create(u)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}
		if u, err := userDao.Get("not_found"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if u != nil {
			t.Fatalf("%s failed: user %s should not exist", name, "not_found")
		}

		if u, err := userDao.Get("btnguyen2k"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if u == nil {
			t.Fatalf("%s failed: nil", name)
		} else {
			if v := u.GetId(); v != "btnguyen2k" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
			}
			if v := u.GetTagVersion(); v != 1357 {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
			}
			if v := u.GetDisplayName(); v != "Thanh Nguyen" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "Thanh Nguyen", v)
			}
			if v := u.GetAesKey(); v != "aeskey" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "aeskey", v)
			}
		}
	}
}

func TestUserDaoSql_Delete(t *testing.T) {
	name := "TestUserDaoSql_Delete"
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
		userDao := _initUserDaoSql(t, name, sqlc)

		u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
		ok, err := userDao.Create(u)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}

		ok, err = userDao.Delete(u)
		if err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if !ok {
			t.Fatalf("%s failed: cannot delete user [%s]", name, u.GetId())
		}

		u, err = userDao.Get("btnguyen2k")
		if app, err := userDao.Get("exter"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if app != nil {
			t.Fatalf("%s failed: user %s should not exist", name, "userDao")
		}
	}
}

func TestUserDaoUser_Update(t *testing.T) {
	name := "TestUserDaoUser_Update"
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
		userDao := _initUserDaoSql(t, name, sqlc)

		u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
		userDao.Create(u)

		u.SetDisplayName("nbthanh")
		u.SetAesKey("newaeskey")
		ok, err := userDao.Update(u)
		if err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", name, ok, err)
		}

		if u, err := userDao.Get("btnguyen2k"); err != nil {
			t.Fatalf("%s failed: %s", name, err)
		} else if u == nil {
			t.Fatalf("%s failed: nil", name)
		} else {
			if v := u.GetId(); v != "btnguyen2k" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
			}
			if v := u.GetTagVersion(); v != 1357 {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
			}
			if v := u.GetDisplayName(); v != "nbthanh" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
			}
			if v := u.GetAesKey(); v != "newaeskey" {
				t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "newaeskey", v)
			}
		}
	}
}

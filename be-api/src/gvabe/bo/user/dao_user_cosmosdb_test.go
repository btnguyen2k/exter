package user

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

func _createCosmosdbConnect(t *testing.T, testName string) *prom.SqlConnect {
	driver := strings.ReplaceAll(os.Getenv("COSMOSDB_DRIVER"), `"`, "")
	url := strings.ReplaceAll(os.Getenv("COSMOSDB_URL"), `"`, "")
	if driver == "" || url == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	timezone := strings.ReplaceAll(os.Getenv("TIMEZONE"), `"`, "")
	if timezone == "" {
		timezone = "UTC"
	}
	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	url += ";Db=exter"
	sqlc, err := henge.NewCosmosdbConnection(url, timezone, driver, 10000, nil)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewCosmosdbConnection", err)
	}
	sqlc.GetDB().Exec("CREATE DATABASE exter WITH maxru=10000")
	return sqlc
}

const tableNameCosmosdb = "exter_test_user"

var setupTestCosmosdb = func(t *testing.T, testName string) {
	testSqlc = _createCosmosdbConnect(t, testName)
	testSqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameCosmosdb))
	err := InitUserTableCosmosdb(testSqlc, tableNameCosmosdb)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestCosmosdb = func(t *testing.T, testName string) {
	if testSqlc != nil {
		defer func() {
			defer func() { testSqlc = nil }()
			testSqlc.Close()
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewUserDaoCosmosdb(t *testing.T) {
	testName := "TestNewUserDaoCosmosdb"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)
	if userDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

// func _initUserDaoCosmosdb(t *testing.T, testName string, sqlc *prom.SqlConnect) UserDao {
// 	if _, err := sqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameCosmosdb)); err != nil {
// 		t.Fatalf("%s failed: %s", testName+"/DROP COLLECTION", err)
// 	}
// 	err := henge.InitCosmosdbCollection(sqlc, tableNameCosmosdb, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbPkName})
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", testName+"/InitCosmosdbCollection", err)
// 	}
// 	return NewUserDaoCosmosdb(sqlc, tableNameCosmosdb)
// }

func _ensureCosmosdbNumRows(t *testing.T, testName string, sqlc *prom.SqlConnect, numRows int) {
	if dbRows, err := sqlc.GetDB().Query(fmt.Sprintf("SELECT COUNT(1) FROM %s c WITH cross_partition=true", tableNameCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if rows, err := sqlc.FetchRows(dbRows); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if value := rows[0]["$1"]; int(value.(float64)) != numRows {
		t.Fatalf("%s failed: expected collection to have %#v rows but received %#v", testName, numRows, value)
	}
}

func TestUserDaoCosmosdb_Create(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Create"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}

	_ensureCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestUserDaoCosmosdb_Get(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Get"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)

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
}

func TestUserDaoCosmosdb_Delete(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Delete"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)

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

	_ensureCosmosdbNumRows(t, testName, testSqlc, 0)
}

func TestUserDaoCosmosdb_Update(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Update"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)

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

	_ensureCosmosdbNumRows(t, testName, testSqlc, 1)
}

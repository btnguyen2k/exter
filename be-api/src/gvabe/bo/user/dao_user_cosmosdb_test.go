package user

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
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

func TestNewUserDaoCosmosdb(t *testing.T) {
	name := "TestNewUserDaoCosmosdb"
	sqlc := _createCosmosdbConnect(t, name)
	userDao := NewUserDaoCosmosdb(sqlc, tableNameCosmosdb)
	if userDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initUserDaoCosmosdb(t *testing.T, testName string, sqlc *prom.SqlConnect) UserDao {
	if _, err := sqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName+"/DROP COLLECTION", err)
	}
	err := henge.InitCosmosdbCollection(sqlc, tableNameCosmosdb, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbPkName})
	if err != nil {
		t.Fatalf("%s failed: %s", testName+"/InitCosmosdbCollection", err)
	}
	return NewUserDaoCosmosdb(sqlc, tableNameCosmosdb)
}

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
	name := "TestUserDaoCosmosdb_Create"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoCosmosdb(t, name, sqlc)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	_ensureCosmosdbNumRows(t, name, sqlc, 1)
}

func TestUserDaoCosmosdb_Get(t *testing.T) {
	name := "TestUserDaoCosmosdb_Get"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoCosmosdb(t, name, sqlc)

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

func TestUserDaoCosmosdb_Delete(t *testing.T) {
	name := "TestUserDaoCosmosdb_Delete"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoCosmosdb(t, name, sqlc)

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

	_ensureCosmosdbNumRows(t, name, sqlc, 0)
}

func TestUserDaoCosmosdb_Update(t *testing.T) {
	name := "TestUserDaoCosmosdb_Update"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoCosmosdb(t, name, sqlc)

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

	_ensureCosmosdbNumRows(t, name, sqlc, 1)
}

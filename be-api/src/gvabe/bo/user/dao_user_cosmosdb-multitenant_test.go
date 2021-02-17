package user

import (
	"fmt"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

const tableNameMultitenantCosmosdb = "exter_test"

func TestNewUserDaoMultitenantCosmosdb(t *testing.T) {
	name := "TestNewUserDaoMultitenantCosmosdb"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := NewUserDaoMultitenantCosmosdb(sqlc, tableNameMultitenantCosmosdb)
	if userDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initUserDaoMultitenantCosmosdb(t *testing.T, testName string, sqlc *prom.SqlConnect) UserDao {
	if _, err := sqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameMultitenantCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName+"/DROP COLLECTION", err)
	}
	err := henge.InitCosmosdbCollection(sqlc, tableNameMultitenantCosmosdb, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbMultitenantPkName})
	if err != nil {
		t.Fatalf("%s failed: %s", testName+"/InitCosmosdbCollection", err)
	}
	return NewUserDaoMultitenantCosmosdb(sqlc, tableNameMultitenantCosmosdb)
}

func _ensureMultitenantCosmosdbNumRows(t *testing.T, testName string, sqlc *prom.SqlConnect, numRows int) {
	if dbRows, err := sqlc.GetDB().Query(fmt.Sprintf("SELECT COUNT(1) FROM %s c WITH cross_partition=true", tableNameMultitenantCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if rows, err := sqlc.FetchRows(dbRows); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if value := rows[0]["$1"]; int(value.(float64)) != numRows {
		t.Fatalf("%s failed: expected collection to have %#v rows but received %#v", testName, numRows, value)
	}
}

func TestUserDaoMultitenantCosmosdb_Create(t *testing.T) {
	name := "TestUserDaoMultitenantCosmosdb_Create"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoMultitenantCosmosdb(t, name, sqlc)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 1)
}

func TestUserDaoMultitenantCosmosdb_Get(t *testing.T) {
	name := "TestUserDaoMultitenantCosmosdb_Get"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoMultitenantCosmosdb(t, name, sqlc)

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

func TestUserDaoMultitenantCosmosdb_Delete(t *testing.T) {
	name := "TestUserDaoMultitenantCosmosdb_Delete"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoMultitenantCosmosdb(t, name, sqlc)

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

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 0)
}

func TestUserDaoMultitenantCosmosdb_Update(t *testing.T) {
	name := "TestUserDaoMultitenantCosmosdb_Update"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	userDao := _initUserDaoMultitenantCosmosdb(t, name, sqlc)

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

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 1)
}

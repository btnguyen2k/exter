package app

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
	"main/src/gvabe/bo/user"
)

const tableNameMultitenantCosmosdb = "exter_test"

func TestNewAppDaoMultitenantCosmosdb(t *testing.T) {
	name := "tableNameMultitenantCosmosdb"
	sqlc := _createCosmosdbConnect(t, name)
	appDao := NewAppDaoMultitenantCosmosdb(sqlc, tableNameMultitenantCosmosdb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initAppDaoMultitenantCosmosdb(t *testing.T, testName string, sqlc *prom.SqlConnect) AppDao {
	if _, err := sqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameMultitenantCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName+"/DROP COLLECTION", err)
	}
	err := henge.InitCosmosdbCollection(sqlc, tableNameMultitenantCosmosdb, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbMultitenantPkName})
	if err != nil {
		t.Fatalf("%s failed: %s", testName+"/InitCosmosdbCollection", err)
	}
	return NewAppDaoMultitenantCosmosdb(sqlc, tableNameMultitenantCosmosdb)
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

func TestAppDaoMultitenantCosmosdb_Create(t *testing.T) {
	name := "TestAppDaoMultitenantCosmosdb_Create"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	appDao := _initAppDaoMultitenantCosmosdb(t, name, sqlc)

	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	ok, err := appDao.Create(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 1)
}

func TestAppDaoMultitenantCosmosdb_Get(t *testing.T) {
	name := "TestAppDaoMultitenantCosmosdb_Get"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	appDao := _initAppDaoMultitenantCosmosdb(t, name, sqlc)

	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	if app, err := appDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", name, "not_found")
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := app.GetTagVersion(); v != 1357 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
		}
		if v := app.GetOwnerId(); v != "btnguyen2k" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
		}
		if v := app.GetAttrsPublic().Description; v != "System application (do not delete)" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "System application (do not delete)", v)
		}
	}
}

func TestAppDaoMultitenantCosmosdb_Delete(t *testing.T) {
	name := "TestAppDaoMultitenantCosmosdb_Delete"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	appDao := _initAppDaoMultitenantCosmosdb(t, name, sqlc)

	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	app, err := appDao.Get("exter")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	}

	ok, err := appDao.Delete(app)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if !ok {
		t.Fatalf("%s failed: cannot delete app [%s]", name, app.GetId())
	}

	app, err = appDao.Get("exter")
	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", name, "exter")
	}

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 0)
}

func TestAppDaoMultitenantCosmosdb_Update(t *testing.T) {
	name := "TestAppDaoMultitenantCosmosdb_Update"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	appDao := _initAppDaoMultitenantCosmosdb(t, name, sqlc)

	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	appDao.Create(app)

	app.SetOwnerId("nbthanh")
	app.SetTagVersion(2468)
	app.attrsPublic.Description = "App description"
	ok, err := appDao.Update(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := app.GetTagVersion(); v != 2468 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
		}
		if v := app.GetOwnerId(); v != "nbthanh" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
		}
		if v := app.GetAttrsPublic().Description; v != "App description" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "App description", v)
		}
	}

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 1)
}

func TestAppDaoMultitenantCosmosdb_GetUserApps(t *testing.T) {
	name := "TestAppDaoMultitenantCosmosdb_GetUserApps"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	appDao := _initAppDaoMultitenantCosmosdb(t, name, sqlc)

	for i := 0; i < 10; i++ {
		app := NewApp(uint64(i), strconv.Itoa(i), strconv.Itoa(i%3), "App #"+strconv.Itoa(i))
		appDao.Create(app)
	}

	u := user.NewUser(123, "2")
	appList, err := appDao.GetUserApps(u)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(appList) != 3 {
		t.Fatalf("%s failed: expected %#v apps but received %#v", name, 3, len(appList))
	}
	for _, app := range appList {
		if app.GetOwnerId() != "2" {
			t.Fatalf("%s failed: app %#v does not belong to user %#v", name, app.GetId(), "2")
		}
	}

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 10)
}

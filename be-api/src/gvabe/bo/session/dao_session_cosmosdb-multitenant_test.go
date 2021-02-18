package session

import (
	"fmt"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

const tableNameMultitenantCosmosdb = "exter_test"

func TestNewSessionDaoMultitenantCosmosdb(t *testing.T) {
	name := "TestNewSessionDaoMultitenantCosmosdb"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	sessDao := NewSessionDaoMultitenantCosmosdb(sqlc, tableNameMultitenantCosmosdb)
	if sessDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initSessionDaoMultitenantCosmosdb(t *testing.T, testName string, sqlc *prom.SqlConnect) SessionDao {
	if _, err := sqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameMultitenantCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName+"/DROP COLLECTION", err)
	}
	err := henge.InitCosmosdbCollection(sqlc, tableNameMultitenantCosmosdb, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbMultitenantPkName})
	if err != nil {
		t.Fatalf("%s failed: %s", testName+"/InitCosmosdbCollection", err)
	}
	return NewSessionDaoMultitenantCosmosdb(sqlc, tableNameMultitenantCosmosdb)
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

func TestSessionDaoMultitenantCosmosdb_Save(t *testing.T) {
	name := "TestSessionDaoMultitenantCosmosdb_Save"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	sessDao := _initSessionDaoMultitenantCosmosdb(t, name, sqlc)

	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 1)
}

func TestSessionDaoMultitenantCosmosdb_Get(t *testing.T) {
	name := "TestSessionDaoMultitenantCosmosdb_Get"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	sessDao := _initSessionDaoMultitenantCosmosdb(t, name, sqlc)

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

func TestSessionDaoMultitenantCosmosdb_Delete(t *testing.T) {
	name := "TestSessionDaoMultitenantCosmosdb_Delete"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	sessDao := _initSessionDaoMultitenantCosmosdb(t, name, sqlc)

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

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 0)
}

func TestSessionDaoMultitenantCosmosdb_Update(t *testing.T) {
	name := "TestSessionDaoMultitenantCosmosdb_Update"
	sqlc := _createCosmosdbConnect(t, name)
	defer sqlc.Close()
	sessDao := _initSessionDaoMultitenantCosmosdb(t, name, sqlc)

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

	_ensureMultitenantCosmosdbNumRows(t, name, sqlc, 1)
}

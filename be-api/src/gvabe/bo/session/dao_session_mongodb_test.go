package session

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

func _createMongoConnect(t *testing.T, testName string) *prom.MongoConnect {
	mongoDb := strings.ReplaceAll(os.Getenv("MONGO_DB"), `"`, "")
	mongoUrl := strings.ReplaceAll(os.Getenv("MONGO_URL"), `"`, "")
	if mongoDb == "" || mongoUrl == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	mc, err := prom.NewMongoConnect(mongoUrl, mongoDb, 10000)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewMongoConnect", err)
	}
	return mc
}

const collectionNameMongo = "exter_test_session"

func TestNewSessionDaoMongo(t *testing.T) {
	name := "TestNewSessionDaoMongo"
	mc := _createMongoConnect(t, name)
	appDao := NewSessionDaoMongo(mc, collectionNameMongo)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initSessionDaoMongo(t *testing.T, testName string, mc *prom.MongoConnect) SessionDao {
	mc.GetCollection(collectionNameMongo).Drop(nil)
	henge.InitMongoCollection(mc, collectionNameMongo)
	return NewSessionDaoMongo(mc, collectionNameMongo)
}

func TestSessionDaoMongo_Save(t *testing.T) {
	name := "TestSessionDaoMongo_Save"
	mc := _createMongoConnect(t, name)
	sessDao := _initSessionDaoMongo(t, name, mc)
	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
}

func TestSessionDaoMongo_Get(t *testing.T) {
	name := "TestSessionDaoMongo_Get"
	mc := _createMongoConnect(t, name)
	sessDao := _initSessionDaoMongo(t, name, mc)
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

func TestSessionDaoMongo_Delete(t *testing.T) {
	name := "TestSessionDaoMongo_Delete"
	mc := _createMongoConnect(t, name)
	sessDao := _initSessionDaoMongo(t, name, mc)
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

func TestSessionDaoMongo_Update(t *testing.T) {
	name := "TestSessionDaoMongo_Update"
	mc := _createMongoConnect(t, name)
	sessDao := _initSessionDaoMongo(t, name, mc)
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

package user

import (
	"os"
	"strings"
	"testing"

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

const collectionNameMongo = "exter_test_user"

var setupTestMongo = func(t *testing.T, testName string) {
	testMc = _createMongoConnect(t, testName)
	testMc.GetCollection(collectionNameMongo).Drop(nil)
	err := InitUserTableMongo(testMc, collectionNameMongo)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestMongo = func(t *testing.T, testName string) {
	if testMc != nil {
		defer func() {
			defer func() { testMc = nil }()
			testMc.Close(nil)
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewUserDaoMongo(t *testing.T) {
	testName := "TestNewUserDaoMongo"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)
	if userDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestUserDaoMongo_Create(t *testing.T) {
	testName := "TestUserDaoMongo_Create"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}
}

func TestUserDaoMongo_Get(t *testing.T) {
	testName := "TestUserDaoMongo_Get"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)

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

func TestUserDaoMongo_Delete(t *testing.T) {
	testName := "TestUserDaoMongo_Delete"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)

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
}

func TestUserDaoMongo_Update(t *testing.T) {
	testName := "TestUserDaoMongo_Update"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)

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
}

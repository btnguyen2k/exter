package user

import (
	"os"
	"strings"
	"testing"

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

const collectionNameMongo = "exter_test_user"

func TestNewUserDaoMongo(t *testing.T) {
	name := "TestNewUserDaoMongo"
	mc := _createMongoConnect(t, name)
	userDao := NewUserDaoMongo(mc, collectionNameMongo)
	if userDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initUserDaoMongo(t *testing.T, testName string, mc *prom.MongoConnect) UserDao {
	mc.GetCollection(collectionNameMongo).Drop(nil)
	henge.InitMongoCollection(mc, collectionNameMongo)
	return NewUserDaoMongo(mc, collectionNameMongo)
}

func TestUserDaoMongo_Create(t *testing.T) {
	name := "TestUserDaoMongo_Create"
	mc := _createMongoConnect(t, name)
	userDao := _initUserDaoMongo(t, name, mc)
	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
}

func TestUserDaoMongo_Get(t *testing.T) {
	name := "TestUserDaoMongo_Get"
	mc := _createMongoConnect(t, name)
	userDao := _initUserDaoMongo(t, name, mc)
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

func TestUserDaoMongo_Delete(t *testing.T) {
	name := "TestUserDaoMongo_Delete"
	mc := _createMongoConnect(t, name)
	userDao := _initUserDaoMongo(t, name, mc)

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

func TestUserDaoMongo_Update(t *testing.T) {
	name := "TestUserDaoMongo_Update"
	mc := _createMongoConnect(t, name)
	userDao := _initUserDaoMongo(t, name, mc)

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

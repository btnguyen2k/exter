package app

import (
	"strconv"
	"testing"

	"github.com/btnguyen2k/prom"
	"main/src/gvabe/bo/user"
)

type TestSetupOrTeardownFunc func(t *testing.T, testName string)

func setupTest(t *testing.T, testName string, extraSetupFunc, extraTeardownFunc TestSetupOrTeardownFunc) func(t *testing.T) {
	if extraSetupFunc != nil {
		extraSetupFunc(t, testName)
	}
	return func(t *testing.T) {
		if extraTeardownFunc != nil {
			extraTeardownFunc(t, testName)
		}
	}
}

var (
	testAdc  *prom.AwsDynamodbConnect
	testMc   *prom.MongoConnect
	testSqlc *prom.SqlConnect
)

/*----------------------------------------------------------------------*/

func doTestAppDao_Create(t *testing.T, testName string, appDao AppDao) {
	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	ok, err := appDao.Create(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}
}

func doTestAppDao_Get(t *testing.T, testName string, appDao AppDao) {
	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))

	if app, err := appDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", testName, "not_found")
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", testName)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "exter", v)
		}
		if v := app.GetTagVersion(); v != 1357 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, 1357, v)
		}
		if v := app.GetOwnerId(); v != "btnguyen2k" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "btnguyen2k", v)
		}
		if v := app.GetAttrsPublic().Description; v != "System application (do not delete)" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "System application (do not delete)", v)
		}
	}
}

func doTestAppDao_Delete(t *testing.T, testName string, appDao AppDao) {
	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	app, err := appDao.Get("exter")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", testName)
	}

	ok, err := appDao.Delete(app)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if !ok {
		t.Fatalf("%s failed: cannot delete app [%s]", testName, app.GetId())
	}

	app, err = appDao.Get("exter")
	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", testName, "exter")
	}
}

func doTestAppDao_Update(t *testing.T, testName string, appDao AppDao) {
	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	appDao.Create(app)

	app.SetOwnerId("nbthanh")
	app.SetTagVersion(2468)
	app.attrsPublic.Description = "App description"
	ok, err := appDao.Update(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", testName)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "exter", v)
		}
		if v := app.GetTagVersion(); v != 2468 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, 2468, v)
		}
		if v := app.GetOwnerId(); v != "nbthanh" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "nbthanh", v)
		}
		if v := app.GetAttrsPublic().Description; v != "App description" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", testName, "App description", v)
		}
	}
}

func doTestAppDao_GetUserApps(t *testing.T, testName string, appDao AppDao) {
	for i := 0; i < 10; i++ {
		app := NewApp(uint64(i), strconv.Itoa(i), strconv.Itoa(i%3), "App #"+strconv.Itoa(i))
		appDao.Create(app)
	}

	u := user.NewUser(123, "2")
	appList, err := appDao.GetUserApps(u)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(appList) != 3 {
		t.Fatalf("%s failed: expected %#v apps but received %#v", testName, 3, len(appList))
	}
	for _, app := range appList {
		if app.GetOwnerId() != "2" {
			t.Fatalf("%s failed: app %#v does not belong to user %#v", testName, app.GetId(), "2")
		}
	}
}

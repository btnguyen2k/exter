package app

import (
	"strconv"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

	"main/src/gvabe/bo"
	"main/src/gvabe/bo/user"
)

const tableNameMultitenantDynamodb = "exter_test"

func TestNewAppDaoMultitenantAwsDynamodb(t *testing.T) {
	name := "TestNewAppDaoMultitenantAwsDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := NewAppDaoMultitenantAwsDynamodb(adc, tableNameMultitenantDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initAppDaoMultitenantDynamodb(t *testing.T, testName string, adc *prom.AwsDynamodbConnect) AppDao {
	err := adc.DeleteTable(nil, tableNameMultitenantDynamodb)
	if err = prom.AwsIgnoreErrorIfMatched(err, awsdynamodb.ErrCodeTableNotFoundException); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	_waitForTable(adc, tableNameMultitenantDynamodb, []string{""}, 1)
	err = henge.InitDynamodbTables(adc, tableNameMultitenantDynamodb, &henge.DynamodbTablesSpec{
		MainTableRcu:         2,
		MainTableWcu:         2,
		MainTableCustomAttrs: []prom.AwsDynamodbNameAndType{{Name: bo.DynamodbMultitenantPkName, Type: prom.AwsAttrTypeString}},
		MainTablePkPrefix:    bo.DynamodbMultitenantPkName,
		CreateUidxTable:      true,
		UidxTableRcu:         2,
		UidxTableWcu:         2,
	})
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	_waitForTable(adc, tableNameMultitenantDynamodb, []string{"ACTIVE"}, 1)
	return NewAppDaoMultitenantAwsDynamodb(adc, tableNameMultitenantDynamodb)
}

func TestAppDaoMultitenantAwsDynamodb_Create(t *testing.T) {
	name := "TestAppDaoMultitenantAwsDynamodb_Create"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoMultitenantDynamodb(t, name, adc)

	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	ok, err := appDao.Create(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueApp {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueApp, items[0])
	}
}

func TestAppDaoMultitenantAwsDynamodb_Get(t *testing.T) {
	name := "TestAppDaoMultitenantAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoMultitenantDynamodb(t, name, adc)

	ok, err := appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}
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

func TestAppDaoMultitenantAwsDynamodb_Delete(t *testing.T) {
	name := "TestAppDaoMultitenantAwsDynamodb_Delete"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoMultitenantDynamodb(t, name, adc)

	ok, err := appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	app, err := appDao.Get("exter")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	}

	ok, err = appDao.Delete(app)
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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 0 item inserted but received %#v", name, len(items))
	}
}

func TestAppDaoMultitenantAwsDynamodb_Update(t *testing.T) {
	name := "TestAppDaoMultitenantAwsDynamodb_Update"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoMultitenantDynamodb(t, name, adc)

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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueApp {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueApp, items[0])
	}
}

func TestAppDaoMultitenantAwsDynamodb_GetUserApps(t *testing.T) {
	name := "TestAppDaoMultitenantAwsDynamodb_GetUserApps"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoMultitenantDynamodb(t, name, adc)

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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 10 {
		t.Fatalf("%s failed: expected 10 items inserted but received %#v", name, len(items))
	}
	for _, item := range items {
		if v, _ := item[bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueApp {
			t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueApp, items[0])
		}
	}
}

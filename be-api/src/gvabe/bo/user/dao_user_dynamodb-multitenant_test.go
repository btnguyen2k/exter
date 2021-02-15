package user

import (
	"testing"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

const tableNameMultitenantDynamodb = "exter_test"

func TestNewUserDaoMultitenantAwsDynamodb(t *testing.T) {
	name := "TestNewUserDaoMultitenantAwsDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := NewUserDaoMultitenantAwsDynamodb(adc, tableNameMultitenantDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initUserDaoMultitenantDynamodb(t *testing.T, testName string, adc *prom.AwsDynamodbConnect) UserDao {
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
	return NewUserDaoMultitenantAwsDynamodb(adc, tableNameMultitenantDynamodb)
}

func TestUserDaoMultitenantAwsDynamodb_Create(t *testing.T) {
	name := "TestUserDaoMultitenantAwsDynamodb_Create"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoMultitenantDynamodb(t, name, adc)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
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
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueUser {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueUser, items[0])
	}
}

func TestUserDaoMultitenantAwsDynamodb_Get(t *testing.T) {
	name := "TestUserDaoMultitenantAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoMultitenantDynamodb(t, name, adc)

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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueUser {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueUser, items[0])
	}
}

func TestUserDaoMultitenantAwsDynamodb_Delete(t *testing.T) {
	name := "TestUserDaoMultitenantAwsDynamodb_Delete"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoMultitenantDynamodb(t, name, adc)

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

func TestUserDaoMultitenantAwsDynamodb_Update(t *testing.T) {
	name := "TestUserDaoMultitenantAwsDynamodb_Update"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoMultitenantDynamodb(t, name, adc)

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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueUser {
		t.Fatalf("%s failed: expected item has field %s with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueUser, items[0])
	}
}

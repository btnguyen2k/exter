package user

import (
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

const tableNameMultitenantDynamodb = "exter_test"

var setupTestDynamodbMultitenant = func(t *testing.T, testName string) {
	testAdc = _createAwsDynamodbConnect(t, testName)
	for _, tableName := range []string{tableNameMultitenantDynamodb, tableNameMultitenantDynamodb + henge.AwsDynamodbUidxTableSuffix} {
		testAdc.DeleteTable(nil, tableName)
		err := prom.AwsDynamodbWaitForTableStatus(testAdc, tableName, []string{""}, 1*time.Second, 10*time.Second)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}
	err := bo.InitMultitenantTableAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestDynamodbMultitenant = func(t *testing.T, testName string) {
	if testAdc != nil {
		defer func() {
			defer func() { testAdc = nil }()
			testAdc.Close()
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewUserDaoMultitenantAwsDynamodb(t *testing.T) {
	testName := "TestNewUserDaoMultitenantAwsDynamodb"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	if userDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestUserDaoMultitenantAwsDynamodb_Create(t *testing.T) {
	testName := "TestUserDaoMultitenantAwsDynamodb_Create"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}

	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueUser {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueUser, items[0])
	}
}

func TestUserDaoMultitenantAwsDynamodb_Get(t *testing.T) {
	testName := "TestUserDaoMultitenantAwsDynamodb_Get"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)

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

	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueUser {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueUser, items[0])
	}
}

func TestUserDaoMultitenantAwsDynamodb_Delete(t *testing.T) {
	testName := "TestUserDaoMultitenantAwsDynamodb_Delete"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)

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

	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
}

func TestUserDaoMultitenantAwsDynamodb_Update(t *testing.T) {
	testName := "TestUserDaoMultitenantAwsDynamodb_Update"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)

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

	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueUser {
		t.Fatalf("%s failed: expected item has field %s with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueUser, items[0])
	}
}

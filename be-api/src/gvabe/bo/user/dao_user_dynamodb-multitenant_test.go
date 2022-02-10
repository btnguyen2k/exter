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
	doTestUserDao_Create(t, testName, userDao)
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
	doTestUserDao_Get(t, testName, userDao)
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
	doTestUserDao_Delete(t, testName, userDao)
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
	doTestUserDao_Update(t, testName, userDao)
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

package app

import (
	"fmt"
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

func TestNewAppDaoMultitenantAwsDynamodb(t *testing.T) {
	testName := "TestNewAppDaoMultitenantAwsDynamodb"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestAppDaoMultitenantAwsDynamodb_Create(t *testing.T) {
	testName := "TestAppDaoMultitenantAwsDynamodb_Create"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestAppDao_Create(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueApp {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueApp, items[0])
	}
}

func TestAppDaoMultitenantAwsDynamodb_Get(t *testing.T) {
	testName := "TestAppDaoMultitenantAwsDynamodb_Get"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestAppDao_Get(t, testName, appDao)
}

func TestAppDaoMultitenantAwsDynamodb_Delete(t *testing.T) {
	testName := "TestAppDaoMultitenantAwsDynamodb_Delete"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestAppDao_Delete(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 0 {
		for _, item := range items {
			fmt.Printf("\tDEBUG: %#v\n", item)
		}
		t.Fatalf("%s failed: expected 0 item inserted but received %#v", testName, len(items))
	}
}

func TestAppDaoMultitenantAwsDynamodb_Update(t *testing.T) {
	testName := "TestAppDaoMultitenantAwsDynamodb_Update"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestAppDao_Update(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		for _, item := range items {
			fmt.Printf("\tDEBUG: %#v\n", item)
		}
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueApp {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueApp, items[0])
	}
}

func TestAppDaoMultitenantAwsDynamodb_GetUserApps(t *testing.T) {
	testName := "TestAppDaoMultitenantAwsDynamodb_GetUserApps"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestAppDao_GetUserApps(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 10 {
		for _, item := range items {
			fmt.Printf("\tDEBUG: %#v\n", item)
		}
		t.Fatalf("%s failed: expected 10 items inserted but received %#v", testName, len(items))
	}
	for _, item := range items {
		if v, _ := item[bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueApp {
			t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueApp, items[0])
		}
	}
}

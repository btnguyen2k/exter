package session

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

func TestNewSessionDaoMultitenantAwsDynamodb(t *testing.T) {
	testName := "TestNewSessionDaoMultitenantAwsDynamodb"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	if sessDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestSessionDaoMultitenantAwsDynamodb_Save(t *testing.T) {
	testName := "TestSessionDaoMultitenantAwsDynamodb_Save"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestSessionDao_Save(t, testName, sessDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueSession {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueSession, items[0])
	}
}

func TestSessionDaoMultitenantAwsDynamodb_Get(t *testing.T) {
	testName := "TestSessionDaoMultitenantAwsDynamodb_Get"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestSessionDao_Get(t, testName, sessDao)
}

func TestSessionDaoMultitenantAwsDynamodb_Delete(t *testing.T) {
	testName := "TestSessionDaoMultitenantAwsDynamodb_Delete"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestSessionDao_Delete(t, testName, sessDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 0 item inserted but received %#v", testName, len(items))
	}
}

func TestSessionDaoMultitenantAwsDynamodb_Update(t *testing.T) {
	testName := "TestSessionDaoMultitenantAwsDynamodb_Update"
	teardownTest := setupTest(t, testName, setupTestDynamodbMultitenant, teardownTestDynamodbMultitenant)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantAwsDynamodb(testAdc, tableNameMultitenantDynamodb)
	doTestSessionDao_Update(t, testName, sessDao)
	items, err := testAdc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueSession {
		t.Fatalf("%s failed: expected item has field %s with value %s but received %#v", testName, bo.DynamodbMultitenantPkName, dynamodbPkValueSession, items[0])
	}
}

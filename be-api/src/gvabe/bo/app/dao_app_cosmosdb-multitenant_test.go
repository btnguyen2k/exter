package app

import (
	"fmt"
	"testing"

	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

const tableNameMultitenantCosmosdb = "exter_test"

var setupTestMultitenantCosmosdb = func(t *testing.T, testName string) {
	testSqlc = _createCosmosdbConnect(t, testName)
	testSqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameMultitenantCosmosdb))
	err := bo.InitMultitenantTableCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestMultitenantCosmosdb = func(t *testing.T, testName string) {
	if testSqlc != nil {
		defer func() {
			defer func() { testSqlc = nil }()
			testSqlc.Close()
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewAppDaoMultitenantCosmosdb(t *testing.T) {
	testName := "tableNameMultitenantCosmosdb"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func _ensureMultitenantCosmosdbNumRows(t *testing.T, testName string, sqlc *prom.SqlConnect, numRows int) {
	if dbRows, err := sqlc.GetDB().Query(fmt.Sprintf("SELECT COUNT(1) FROM %s c WITH cross_partition=true", tableNameMultitenantCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if rows, err := sqlc.FetchRows(dbRows); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if value := rows[0]["$1"]; int(value.(float64)) != numRows {
		t.Fatalf("%s failed: expected collection to have %#v rows but received %#v", testName, numRows, value)
	}
}

func TestAppDaoMultitenantCosmosdb_Create(t *testing.T) {
	testName := "TestAppDaoMultitenantCosmosdb_Create"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestAppDao_Create(t, testName, appDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestAppDaoMultitenantCosmosdb_Get(t *testing.T) {
	testName := "TestAppDaoMultitenantCosmosdb_Get"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestAppDao_Get(t, testName, appDao)
}

func TestAppDaoMultitenantCosmosdb_Delete(t *testing.T) {
	testName := "TestAppDaoMultitenantCosmosdb_Delete"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestAppDao_Delete(t, testName, appDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 0)
}

func TestAppDaoMultitenantCosmosdb_Update(t *testing.T) {
	testName := "TestAppDaoMultitenantCosmosdb_Update"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestAppDao_Update(t, testName, appDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestAppDaoMultitenantCosmosdb_GetUserApps(t *testing.T) {
	testName := "TestAppDaoMultitenantCosmosdb_GetUserApps"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	appDao := NewAppDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestAppDao_GetUserApps(t, testName, appDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 10)
}

package session

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

func TestNewSessionDaoMultitenantCosmosdb(t *testing.T) {
	testName := "TestNewSessionDaoMultitenantCosmosdb"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	if sessDao == nil {
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

func TestSessionDaoMultitenantCosmosdb_Save(t *testing.T) {
	testName := "TestSessionDaoMultitenantCosmosdb_Save"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestSessionDao_Save(t, testName, sessDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestSessionDaoMultitenantCosmosdb_Get(t *testing.T) {
	testName := "TestSessionDaoMultitenantCosmosdb_Get"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestSessionDao_Get(t, testName, sessDao)
}

func TestSessionDaoMultitenantCosmosdb_Delete(t *testing.T) {
	testName := "TestSessionDaoMultitenantCosmosdb_Delete"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestSessionDao_Delete(t, testName, sessDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 0)
}

func TestSessionDaoMultitenantCosmosdb_Update(t *testing.T) {
	testName := "TestSessionDaoMultitenantCosmosdb_Update"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestSessionDao_Update(t, testName, sessDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 1)
}

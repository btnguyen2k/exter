package user

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

func TestNewUserDaoMultitenantCosmosdb(t *testing.T) {
	testName := "TestNewUserDaoMultitenantCosmosdb"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	if userDao == nil {
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

func TestUserDaoMultitenantCosmosdb_Create(t *testing.T) {
	testName := "TestUserDaoMultitenantCosmosdb_Create"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestUserDao_Create(t, testName, userDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestUserDaoMultitenantCosmosdb_Get(t *testing.T) {
	testName := "TestUserDaoMultitenantCosmosdb_Get"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestUserDao_Get(t, testName, userDao)
}

func TestUserDaoMultitenantCosmosdb_Delete(t *testing.T) {
	testName := "TestUserDaoMultitenantCosmosdb_Delete"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestUserDao_Delete(t, testName, userDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 0)
}

func TestUserDaoMultitenantCosmosdb_Update(t *testing.T) {
	testName := "TestUserDaoMultitenantCosmosdb_Update"
	teardownTest := setupTest(t, testName, setupTestMultitenantCosmosdb, teardownTestMultitenantCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoMultitenantCosmosdb(testSqlc, tableNameMultitenantCosmosdb)
	doTestUserDao_Update(t, testName, userDao)
	_ensureMultitenantCosmosdbNumRows(t, testName, testSqlc, 1)
}

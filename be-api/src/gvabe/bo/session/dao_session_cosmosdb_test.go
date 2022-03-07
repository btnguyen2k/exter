package session

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

func _createCosmosdbConnect(t *testing.T, testName string) *prom.SqlConnect {
	driver := strings.ReplaceAll(os.Getenv("COSMOSDB_DRIVER"), `"`, "")
	url := strings.ReplaceAll(os.Getenv("COSMOSDB_URL"), `"`, "")
	if driver == "" || url == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	timezone := strings.ReplaceAll(os.Getenv("TIMEZONE"), `"`, "")
	if timezone == "" {
		timezone = "UTC"
	}
	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	url += ";Db=exter"
	sqlc, err := henge.NewCosmosdbConnection(url, timezone, driver, 10000, nil)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewCosmosdbConnection", err)
	}
	sqlc.GetDB().Exec("CREATE DATABASE exter WITH maxru=10000")
	return sqlc
}

const tableNameCosmosdb = "exter_test_session"

var setupTestCosmosdb = func(t *testing.T, testName string) {
	testSqlc = _createCosmosdbConnect(t, testName)
	testSqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameCosmosdb))
	err := InitSessionTableCosmosdb(testSqlc, tableNameCosmosdb)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestCosmosdb = func(t *testing.T, testName string) {
	if testSqlc != nil {
		defer func() {
			defer func() { testSqlc = nil }()
			testSqlc.Close()
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewSessionDaoCosmosdb(t *testing.T) {
	testName := "TestNewSessionDaoCosmosdb"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoCosmosdb(testSqlc, tableNameCosmosdb)
	if sessDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func _ensureCosmosdbNumRows(t *testing.T, testName string, sqlc *prom.SqlConnect, numRows int) {
	if dbRows, err := sqlc.GetDB().Query(fmt.Sprintf("SELECT COUNT(1) FROM %s c WITH cross_partition=true", tableNameCosmosdb)); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if rows, err := sqlc.FetchRows(dbRows); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if value := rows[0]["$1"]; int(value.(float64)) != numRows {
		t.Fatalf("%s failed: expected collection to have %#v rows but received %#v", testName, numRows, value)
	}
}

func TestSessionDaoCosmosdb_Save(t *testing.T) {
	testName := "TestSessionDaoCosmosdb_Save"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestSessionDao_Save(t, testName, sessDao)
	_ensureCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestSessionDaoCosmosdb_Get(t *testing.T) {
	testName := "TestSessionDaoCosmosdb_Get"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestSessionDao_Get(t, testName, sessDao)
}

func TestSessionDaoCosmosdb_Delete(t *testing.T) {
	testName := "TestSessionDaoCosmosdb_Delete"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestSessionDao_Delete(t, testName, sessDao)
	_ensureCosmosdbNumRows(t, testName, testSqlc, 0)
}

func TestSessionDaoCosmosdb_Update(t *testing.T) {
	testName := "TestSessionDaoCosmosdb_Update"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	sessDao := NewSessionDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestSessionDao_Update(t, testName, sessDao)
	_ensureCosmosdbNumRows(t, testName, testSqlc, 1)
}

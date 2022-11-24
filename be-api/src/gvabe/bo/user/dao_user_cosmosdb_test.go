package user

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

func _createCosmosdbConnect(t *testing.T, testName string) *prom.SqlConnect {
	driver := strings.Trim(os.Getenv("COSMOSDB_DRIVER"), `'"`)
	url := strings.Trim(os.Getenv("COSMOSDB_URL"), `'"`)
	if driver == "" || url == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	timezone := strings.Trim(os.Getenv("TIMEZONE"), `'"`)
	if timezone == "" {
		timezone = "UTC"
	}
	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	db := "exter"
	dbre := regexp.MustCompile(`(?i);db=(\w+)`)
	findResult := dbre.FindAllStringSubmatch(url, -1)
	if findResult == nil {
		url += ";Db=" + db
	} else {
		db = findResult[0][1]
	}
	sqlc, err := henge.NewCosmosdbConnection(url, timezone, driver, 10000, nil)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewCosmosdbConnection", err)
	}
	sqlc.GetDB().Exec("CREATE DATABASE IF NOT EXISTS " + db + " WITH maxru=10000")
	return sqlc
}

const tableNameCosmosdb = "exter_test_user"

var setupTestCosmosdb = func(t *testing.T, testName string) {
	testSqlc = _createCosmosdbConnect(t, testName)
	testSqlc.GetDB().Exec(fmt.Sprintf("DROP COLLECTION IF EXISTS %s", tableNameCosmosdb))
	err := InitUserTableCosmosdb(testSqlc, tableNameCosmosdb)
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

func TestNewUserDaoCosmosdb(t *testing.T) {
	testName := "TestNewUserDaoCosmosdb"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)
	if userDao == nil {
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

func TestUserDaoCosmosdb_Create(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Create"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestUserDao_Create(t, testName, userDao)
	_ensureCosmosdbNumRows(t, testName, testSqlc, 1)
}

func TestUserDaoCosmosdb_Get(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Get"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestUserDao_Get(t, testName, userDao)
}

func TestUserDaoCosmosdb_Delete(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Delete"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestUserDao_Delete(t, testName, userDao)
	_ensureCosmosdbNumRows(t, testName, testSqlc, 0)
}

func TestUserDaoCosmosdb_Update(t *testing.T) {
	testName := "TestUserDaoCosmosdb_Update"
	teardownTest := setupTest(t, testName, setupTestCosmosdb, teardownTestCosmosdb)
	defer teardownTest(t)
	userDao := NewUserDaoCosmosdb(testSqlc, tableNameCosmosdb)
	doTestUserDao_Update(t, testName, userDao)
	_ensureCosmosdbNumRows(t, testName, testSqlc, 1)
}

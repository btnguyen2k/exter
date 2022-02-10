package session

import (
	"os"
	"strings"
	"testing"

	"github.com/btnguyen2k/prom"
)

func _createMongoConnect(t *testing.T, testName string) *prom.MongoConnect {
	mongoDb := strings.ReplaceAll(os.Getenv("MONGO_DB"), `"`, "")
	mongoUrl := strings.ReplaceAll(os.Getenv("MONGO_URL"), `"`, "")
	if mongoDb == "" || mongoUrl == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	mc, err := prom.NewMongoConnect(mongoUrl, mongoDb, 10000)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewMongoConnect", err)
	}
	return mc
}

const collectionNameMongo = "exter_test_session"

var setupTestMongo = func(t *testing.T, testName string) {
	testMc = _createMongoConnect(t, testName)
	testMc.GetCollection(collectionNameMongo).Drop(nil)
	err := InitSessionTableMongo(testMc, collectionNameMongo)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestMongo = func(t *testing.T, testName string) {
	if testMc != nil {
		defer func() {
			defer func() { testMc = nil }()
			testMc.Close(nil)
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewSessionDaoMongo(t *testing.T) {
	testName := "TestNewSessionDaoMongo"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	sessDao := NewSessionDaoMongo(testMc, collectionNameMongo)
	if sessDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestSessionDaoMongo_Save(t *testing.T) {
	testName := "TestSessionDaoMongo_Save"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	sessDao := NewSessionDaoMongo(testMc, collectionNameMongo)
	doTestSessionDao_Save(t, testName, sessDao)
}

func TestSessionDaoMongo_Get(t *testing.T) {
	testName := "TestSessionDaoMongo_Get"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	sessDao := NewSessionDaoMongo(testMc, collectionNameMongo)
	doTestSessionDao_Get(t, testName, sessDao)
}

func TestSessionDaoMongo_Delete(t *testing.T) {
	testName := "TestSessionDaoMongo_Delete"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	sessDao := NewSessionDaoMongo(testMc, collectionNameMongo)
	doTestSessionDao_Delete(t, testName, sessDao)
}

func TestSessionDaoMongo_Update(t *testing.T) {
	testName := "TestSessionDaoMongo_Update"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	sessDao := NewSessionDaoMongo(testMc, collectionNameMongo)
	doTestSessionDao_Update(t, testName, sessDao)
}

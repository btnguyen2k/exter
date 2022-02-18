package app

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/btnguyen2k/prom"
)

func _createMongoConnect(t *testing.T, testName string) *prom.MongoConnect {
	mongoDb := strings.ReplaceAll(os.Getenv("MONGO_DB"), `"`, "")
	mongoUrl := strings.ReplaceAll(os.Getenv("MONGO_URL"), `"`, "")
	if mongoDb == "" || mongoUrl == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	mongoPoolOpts := &prom.MongoPoolOpts{
		ConnectTimeout:         5 * time.Second,
		SocketTimeout:          7 * time.Second,
		ServerSelectionTimeout: 11 * time.Second,
	}
	mc, err := prom.NewMongoConnectWithPoolOptions(mongoUrl, mongoDb, 10000, mongoPoolOpts)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewMongoConnect", err)
	}
	return mc
}

const collectionNameMongo = "exter_test_app"

var setupTestMongo = func(t *testing.T, testName string) {
	testMc = _createMongoConnect(t, testName)
	testMc.GetCollection(collectionNameMongo).Drop(nil)
	err := InitAppTableMongo(testMc, collectionNameMongo)
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

func TestNewAppDaoMongo(t *testing.T) {
	testName := "TestNewAppDaoMongo"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	appDao := NewAppDaoMongo(testMc, collectionNameMongo)
	if appDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestAppDaoMongo_Create(t *testing.T) {
	testName := "TestAppDaoMongo_Create"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	appDao := NewAppDaoMongo(testMc, collectionNameMongo)
	doTestAppDao_Create(t, testName, appDao)
}

func TestAppDaoMongo_Get(t *testing.T) {
	testName := "TestAppDaoMongo_Get"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	appDao := NewAppDaoMongo(testMc, collectionNameMongo)
	doTestAppDao_Get(t, testName, appDao)
}

func TestAppDaoMongo_Delete(t *testing.T) {
	testName := "TestAppDaoMongo_Delete"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	appDao := NewAppDaoMongo(testMc, collectionNameMongo)
	doTestAppDao_Delete(t, testName, appDao)
}

func TestAppDaoMongo_Update(t *testing.T) {
	testName := "TestAppDaoMongo_Update"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	appDao := NewAppDaoMongo(testMc, collectionNameMongo)
	doTestAppDao_Update(t, testName, appDao)
}

func TestAppDaoMongo_GetUserApps(t *testing.T) {
	testName := "TestAppDaoMongo_GetUserApps"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	appDao := NewAppDaoMongo(testMc, collectionNameMongo)
	doTestAppDao_GetUserApps(t, testName, appDao)
}

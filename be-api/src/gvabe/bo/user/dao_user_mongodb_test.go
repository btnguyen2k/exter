package user

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

const collectionNameMongo = "exter_test_user"

var setupTestMongo = func(t *testing.T, testName string) {
	testMc = _createMongoConnect(t, testName)
	testMc.GetCollection(collectionNameMongo).Drop(nil)
	err := InitUserTableMongo(testMc, collectionNameMongo)
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

func TestNewUserDaoMongo(t *testing.T) {
	testName := "TestNewUserDaoMongo"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)
	if userDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestUserDaoMongo_Create(t *testing.T) {
	testName := "TestUserDaoMongo_Create"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)
	doTestUserDao_Create(t, testName, userDao)
}

func TestUserDaoMongo_Get(t *testing.T) {
	testName := "TestUserDaoMongo_Get"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)
	doTestUserDao_Get(t, testName, userDao)
}

func TestUserDaoMongo_Delete(t *testing.T) {
	testName := "TestUserDaoMongo_Delete"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)
	doTestUserDao_Delete(t, testName, userDao)
}

func TestUserDaoMongo_Update(t *testing.T) {
	testName := "TestUserDaoMongo_Update"
	teardownTest := setupTest(t, testName, setupTestMongo, teardownTestMongo)
	defer teardownTest(t)
	userDao := NewUserDaoMongo(testMc, collectionNameMongo)
	doTestUserDao_Update(t, testName, userDao)
}

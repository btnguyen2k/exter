package user

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/btnguyen2k/prom"
)

func _createAwsDynamodbConnect(t *testing.T, testName string) *prom.AwsDynamodbConnect {
	awsRegion := strings.ReplaceAll(os.Getenv("AWS_REGION"), `"`, "")
	awsAccessKeyId := strings.ReplaceAll(os.Getenv("AWS_ACCESS_KEY_ID"), `"`, "")
	awsSecretAccessKey := strings.ReplaceAll(os.Getenv("AWS_SECRET_ACCESS_KEY"), `"`, "")
	if awsRegion == "" || awsAccessKeyId == "" || awsSecretAccessKey == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	cfg := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewEnvCredentials(),
	}
	if awsDynamodbEndpoint := strings.ReplaceAll(os.Getenv("AWS_DYNAMODB_ENDPOINT"), `"`, ""); awsDynamodbEndpoint != "" {
		cfg.Endpoint = aws.String(awsDynamodbEndpoint)
		if strings.HasPrefix(awsDynamodbEndpoint, "http://") {
			cfg.DisableSSL = aws.Bool(true)
		}
	}
	adc, err := prom.NewAwsDynamodbConnect(cfg, nil, nil, 10000)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewAwsDynamodbConnect", err)
	}
	return adc
}

const tableNameDynamodb = "exter_test_user"

var setupTestDynamodb = func(t *testing.T, testName string) {
	testAdc = _createAwsDynamodbConnect(t, testName)
	testAdc.DeleteTable(nil, tableNameDynamodb)
	err := prom.AwsDynamodbWaitForTableStatus(testAdc, tableNameDynamodb, []string{""}, 1*time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	err = InitUserTableAwsDynamodb(testAdc, tableNameDynamodb)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
}

var teardownTestDynamodb = func(t *testing.T, testName string) {
	if testAdc != nil {
		defer func() {
			defer func() { testAdc = nil }()
			testAdc.Close()
		}()
	}
}

/*----------------------------------------------------------------------*/

func TestNewUserDaoAwsDynamodb(t *testing.T) {
	testName := "TestNewUserDaoAwsDynamodb"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	userDao := NewUserDaoAwsDynamodb(testAdc, tableNameDynamodb)
	if userDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestUserDaoAwsDynamodb_Create(t *testing.T) {
	testName := "TestAppDaoAwsDynamodb_Create"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	userDao := NewUserDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestUserDao_Create(t, testName, userDao)
	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
}

func TestUserDaoAwsDynamodb_Get(t *testing.T) {
	testName := "TestUserDaoAwsDynamodb_Get"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	userDao := NewUserDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestUserDao_Get(t, testName, userDao)
}

func TestUserDaoAwsDynamodb_Delete(t *testing.T) {
	testName := "TestUserDaoAwsDynamodb_Delete"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	userDao := NewUserDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestUserDao_Delete(t, testName, userDao)
	if app, err := userDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if app != nil {
		t.Fatalf("%s failed: user %s should not exist", testName, "userDao")
	}

	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
}

func TestUserDaoAwsDynamodb_Update(t *testing.T) {
	testName := "TestUserDaoAwsDynamodb_Update"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	userDao := NewUserDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestUserDao_Update(t, testName, userDao)
	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
}

package app

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

const tableNameDynamodb = "exter_test_app"

var setupTestDynamodb = func(t *testing.T, testName string) {
	testAdc = _createAwsDynamodbConnect(t, testName)
	testAdc.DeleteTable(nil, tableNameDynamodb)
	err := prom.AwsDynamodbWaitForTableStatus(testAdc, tableNameDynamodb, []string{""}, 1*time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	err = InitAppTableAwsDynamodb(testAdc, tableNameDynamodb)
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

func TestNewAppDaoAwsDynamodb(t *testing.T) {
	testName := "TestNewAppDaoAwsDynamodb"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	appDao := NewAppDaoAwsDynamodb(testAdc, tableNameDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", testName)
	}
}

func TestAppDaoAwsDynamodb_Create(t *testing.T) {
	testName := "TestAppDaoAwsDynamodb_Create"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	appDao := NewAppDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestAppDao_Create(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
}

func TestAppDaoAwsDynamodb_Get(t *testing.T) {
	testName := "TestAppDaoAwsDynamodb_Get"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	appDao := NewAppDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestAppDao_Get(t, testName, appDao)
}

func TestAppDaoAwsDynamodb_Delete(t *testing.T) {
	testName := "TestAppDaoAwsDynamodb_Delete"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	appDao := NewAppDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestAppDao_Delete(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 0 item inserted but received %#v", testName, len(items))
	}
}

func TestAppDaoAwsDynamodb_Update(t *testing.T) {
	testName := "TestAppDaoAwsDynamodb_Update"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	appDao := NewAppDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestAppDao_Update(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", testName, len(items))
	}
}

func TestAppDaoAwsDynamodb_GetUserApps(t *testing.T) {
	testName := "TestAppDaoAwsDynamodb_GetUserApps"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	appDao := NewAppDaoAwsDynamodb(testAdc, tableNameDynamodb)
	doTestAppDao_GetUserApps(t, testName, appDao)
	items, err := testAdc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if len(items) != 10 {
		t.Fatalf("%s failed: expected 10 items inserted but received %#v", testName, len(items))
	}
}

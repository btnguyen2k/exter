package session

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/btnguyen2k/henge"
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

const tableNameDynamodb = "exter_test_session"

func TestNewSessionDaoAwsDynamodb(t *testing.T) {
	name := "TestNewSessionDaoAwsDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	appDao := NewSessionDaoAwsDynamodb(adc, tableNameDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initSessionDaoDynamodb(t *testing.T, testName string, adc *prom.AwsDynamodbConnect) SessionDao {
	adc.DeleteTable(nil, tableNameDynamodb)
	henge.InitDynamodbTable(adc, tableNameDynamodb, 2, 2)
	return NewSessionDaoAwsDynamodb(adc, tableNameDynamodb)
}

func TestSessionDaoAwsDynamodb_Save(t *testing.T) {
	name := "TestSessionDaoAwsDynamodb_Save"
	adc := _createAwsDynamodbConnect(t, name)
	sessDao := _initSessionDaoDynamodb(t, name, adc)
	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
}

func TestSessionDaoAwsDynamodb_Get(t *testing.T) {
	name := "TestSessionDaoAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	sessDao := _initSessionDaoDynamodb(t, name, adc)
	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	if sess, err := sessDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if sess != nil {
		t.Fatalf("%s failed: session %s should not exist", name, "not_found")
	}

	if sess, err := sessDao.Get("1"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if sess == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := sess.GetId(); v != "1" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "1", v)
		}
		if v := sess.GetTagVersion(); v != 1357 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
		}
		if v := sess.GetSessionType(); v != "login" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "login", v)
		}
		if v := sess.GetIdSource(); v != "local" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "local", v)
		}
		if v := sess.GetAppId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := sess.GetUserId(); v != "btnguyen2k" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
		}
		if v := sess.GetSessionData(); v != "session-data" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "session-data", v)
		}
		if v := sess.GetExpiry(); v.Unix() != expiry.Unix() {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, expiry, v)
		}
	}
}

func TestSessionDaoAwsDynamodb_Delete(t *testing.T) {
	name := "TestSessionDaoAwsDynamodb_Delete"
	adc := _createAwsDynamodbConnect(t, name)
	sessDao := _initSessionDaoDynamodb(t, name, adc)
	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
	if sess, err := sessDao.Get("1"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if sess == nil {
		t.Fatalf("%s failed: nill", name)
	}

	ok, err = sessDao.Delete(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	if sess, err := sessDao.Get("1"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if sess != nil {
		t.Fatalf("%s failed: session %s should not exist", name, "not_found")
	}
}

func TestSessionDaoAwsDynamodb_Update(t *testing.T) {
	name := "TestSessionDaoAwsDynamodb_Update"
	adc := _createAwsDynamodbConnect(t, name)
	sessDao := _initSessionDaoDynamodb(t, name, adc)
	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)

	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	sess.SetTagVersion(2468)
	sess.SetSessionType("pre-login")
	sess.SetIdSource("external")
	sess.SetAppId("myapp")
	sess.SetUserId("nbthanh")
	sess.SetSessionData("data")
	sess.SetExpiry(expiry.Add(1 * time.Hour))
	ok, err = sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	if sess, err := sessDao.Get("1"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if sess == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := sess.GetId(); v != "1" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "1", v)
		}
		if v := sess.GetTagVersion(); v != 2468 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
		}
		if v := sess.GetSessionType(); v != "pre-login" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "pre-login", v)
		}
		if v := sess.GetIdSource(); v != "external" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "external", v)
		}
		if v := sess.GetAppId(); v != "myapp" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "myapp", v)
		}
		if v := sess.GetUserId(); v != "nbthanh" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
		}
		if v := sess.GetSessionData(); v != "data" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "data", v)
		}
		if v := sess.GetExpiry(); v.Unix() != expiry.Add(1*time.Hour).Unix() {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, expiry.Add(1*time.Hour), v)
		}
	}
}

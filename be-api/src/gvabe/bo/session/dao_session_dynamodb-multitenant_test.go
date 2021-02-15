package session

import (
	"testing"
	"time"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

const tableNameMultitenantDynamodb = "exter_test"

func TestNewSessionDaoMultitenantAwsDynamodb(t *testing.T) {
	name := "TestNewSessionDaoMultitenantAwsDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := NewSessionDaoMultitenantAwsDynamodb(adc, tableNameMultitenantDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initSessionDaoMultitenantDynamodb(t *testing.T, testName string, adc *prom.AwsDynamodbConnect) SessionDao {
	err := adc.DeleteTable(nil, tableNameMultitenantDynamodb)
	if err = prom.AwsIgnoreErrorIfMatched(err, awsdynamodb.ErrCodeTableNotFoundException); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	_waitForTable(adc, tableNameMultitenantDynamodb, []string{""}, 1)
	err = henge.InitDynamodbTables(adc, tableNameMultitenantDynamodb, &henge.DynamodbTablesSpec{
		MainTableRcu:         2,
		MainTableWcu:         2,
		MainTableCustomAttrs: []prom.AwsDynamodbNameAndType{{Name: bo.DynamodbMultitenantPkName, Type: prom.AwsAttrTypeString}},
		MainTablePkPrefix:    bo.DynamodbMultitenantPkName,
		CreateUidxTable:      true,
		UidxTableRcu:         2,
		UidxTableWcu:         2,
	})
	if err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	_waitForTable(adc, tableNameMultitenantDynamodb, []string{"ACTIVE"}, 1)
	return NewSessionDaoMultitenantAwsDynamodb(adc, tableNameMultitenantDynamodb)
}

func TestSessionDaoMultitenantAwsDynamodb_Save(t *testing.T) {
	name := "TestSessionDaoMultitenantAwsDynamodb_Save"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	sessDao := _initSessionDaoMultitenantDynamodb(t, name, adc)

	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueSession {
		t.Fatalf("%s failed: expected item has field '%s' with value '%s' but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueSession, items[0])
	}
}

func TestSessionDaoMultitenantAwsDynamodb_Get(t *testing.T) {
	name := "TestSessionDaoMultitenantAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	sessDao := _initSessionDaoMultitenantDynamodb(t, name, adc)

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

func TestSessionDaoMultitenantAwsDynamodb_Delete(t *testing.T) {
	name := "TestSessionDaoMultitenantAwsDynamodb_Delete"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	sessDao := _initSessionDaoMultitenantDynamodb(t, name, adc)

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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 0 item inserted but received %#v", name, len(items))
	}
}

func TestSessionDaoMultitenantAwsDynamodb_Update(t *testing.T) {
	name := "TestSessionDaoMultitenantAwsDynamodb_Update"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	sessDao := _initSessionDaoMultitenantDynamodb(t, name, adc)

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

	items, err := adc.ScanItems(nil, tableNameMultitenantDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
	if v, _ := items[0][bo.DynamodbMultitenantPkName].(string); v != dynamodbPkValueSession {
		t.Fatalf("%s failed: expected item has field %s with value %s but received %#v", name, bo.DynamodbMultitenantPkName, dynamodbPkValueSession, items[0])
	}
}

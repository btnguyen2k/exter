package session

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
	"main/src/gvabe/bo"
)

type TestSetupOrTeardownFunc func(t *testing.T, testName string)

func setupTest(t *testing.T, testName string, extraSetupFunc, extraTeardownFunc TestSetupOrTeardownFunc) func(t *testing.T) {
	if extraSetupFunc != nil {
		extraSetupFunc(t, testName)
	}
	return func(t *testing.T) {
		if extraTeardownFunc != nil {
			extraTeardownFunc(t, testName)
		}
	}
}

var (
	testAdc  *prom.AwsDynamodbConnect
	testMc   *prom.MongoConnect
	testSqlc *prom.SqlConnect
)

/*----------------------------------------------------------------------*/

func doTestSessionDao_Save(t *testing.T, testName string, sessDao SessionDao) {
	expiry := time.Now().Add(5 * time.Minute).Round(time.Millisecond)
	sess := NewSession(1357, "1", "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}
}

func doTestSessionDao_Get(t *testing.T, testName string, sessDao SessionDao) {
	bo.UboTimestampRounding = henge.TimestampRoundingSettingNone
	_appVersion := uint64(1337)
	_sid := "1"
	_stype := "login"
	_idSrc := "google"
	_appId := "test"
	_userId := "btnguyen2k"
	_sdata := "My session data"
	_sexpiry := time.Now().Add(5 * time.Minute)
	sess := NewSession(_appVersion, _sid, _stype, _idSrc, _appId, _userId, _sdata, _sexpiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}

	if sess, err := sessDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if sess != nil {
		t.Fatalf("%s failed: session %s should not exist", testName, "not_found")
	}

	for _, timeRoundingSetting := range allTimeRoundingSettings {
		t.Run(fmt.Sprintf("%#v", timeRoundingSetting), func(t *testing.T) {
			if sess, err := sessDao.Get("1"); err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if sess == nil {
				t.Fatalf("%s failed: nil", testName)
			} else {
				if f, v, expected := "app-version", sess.GetTagVersion(), _appVersion; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "id", sess.GetId(), _sid; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "sess-type", sess.GetSessionType(), _stype; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "id-source", sess.GetIdSource(), _idSrc; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "app-id", sess.GetAppId(), _appId; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "user-id", sess.GetUserId(), _userId; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "expiry", sess.GetExpiry(), sess.RoundTimestamp(_sexpiry); v.UnixNano() != expected.UnixNano() {
					t.Fatalf("%s failed: expected %s to be %s but received %s", testName, f, expected, v)
				}
				if f, v, expected := "sess-data", sess.GetSessionData(), _sdata; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
			}
		})
	}
}

func doTestSessionDao_Delete(t *testing.T, testName string, sessDao SessionDao) {
	_sid := "1"
	expiry := time.Now().Add(5 * time.Minute)
	sess := NewSession(1357, _sid, "login", "local", "exter", "btnguyen2k", "session-data", expiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}
	if sess, err := sessDao.Get(_sid); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if sess == nil {
		t.Fatalf("%s failed: nill", testName)
	}

	ok, err = sessDao.Delete(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}

	if sess, err := sessDao.Get(_sid); err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	} else if sess != nil {
		t.Fatalf("%s failed: session %s should not exist", testName, "not_found")
	}
}

func doTestSessionDao_Update(t *testing.T, testName string, sessDao SessionDao) {
	bo.UboTimestampRounding = henge.TimestampRoundingSettingNone
	_appVersion := uint64(1337)
	_sid := "1"
	_stype := "login"
	_idSrc := "google"
	_appId := "test"
	_userId := "btnguyen2k"
	_sdata := "My session data"
	_sexpiry := time.Now().Add(5 * time.Minute)
	sess := NewSession(_appVersion, _sid, _stype, _idSrc, _appId, _userId, _sdata, _sexpiry)
	ok, err := sessDao.Save(sess)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", testName, ok, err)
	}

	for runIndex, timeRoundingSetting := range allTimeRoundingSettings {
		t.Run(fmt.Sprintf("%#v", timeRoundingSetting), func(t *testing.T) {
			sess.SetTagVersion(_appVersion + uint64(runIndex))
			sess.SetSessionType(_stype + "-ex" + strconv.Itoa(runIndex))
			sess.SetIdSource(_idSrc + "-ex" + strconv.Itoa(runIndex))
			sess.SetAppId(_appId + "-ex" + strconv.Itoa(runIndex))
			sess.SetUserId(_userId + "-ex" + strconv.Itoa(runIndex))
			sess.SetSessionData(_sdata + "-ex" + strconv.Itoa(runIndex))
			sess.SetExpiry(_sexpiry.Add(time.Duration(runIndex) * time.Hour))
			ok, err = sessDao.Save(sess)
			if err != nil || !ok {
				t.Fatalf("%s failed: %#v / %s", testName, ok, err)
			}

			if sess, err := sessDao.Get("1"); err != nil {
				t.Fatalf("%s failed: %s", testName, err)
			} else if sess == nil {
				t.Fatalf("%s failed: nil", testName)
			} else {
				if f, v, expected := "app-version", sess.GetTagVersion(), _appVersion+uint64(runIndex); v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "id", sess.GetId(), _sid; v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "sess-type", sess.GetSessionType(), _stype+"-ex"+strconv.Itoa(runIndex); v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "id-source", sess.GetIdSource(), _idSrc+"-ex"+strconv.Itoa(runIndex); v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "app-id", sess.GetAppId(), _appId+"-ex"+strconv.Itoa(runIndex); v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "user-id", sess.GetUserId(), _userId+"-ex"+strconv.Itoa(runIndex); v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
				if f, v, expected := "expiry", sess.GetExpiry(), sess.RoundTimestamp(_sexpiry.Add(time.Duration(runIndex)*time.Hour)); v.UnixNano() != expected.UnixNano() {
					t.Fatalf("%s failed: expected %s to be %s but received %s", testName, f, expected, v)
				}
				if f, v, expected := "sess-data", sess.GetSessionData(), _sdata+"-ex"+strconv.Itoa(runIndex); v != expected {
					t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
				}
			}
		})
	}
}

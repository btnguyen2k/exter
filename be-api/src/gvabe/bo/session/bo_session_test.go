package session

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"main/src/gvabe/bo"
)

var allTimeRoundingSettings = []henge.TimestampRoundingSetting{
	henge.TimestampRoundingSettingNone,
	henge.TimestampRoundingSettingNanosecond,
	henge.TimestampRoundingSettingMicrosecond,
	henge.TimestampRoundingSettingMillisecond,
	henge.TimestampRoundingSettingSecond,
}

func TestNewSession(t *testing.T) {
	testName := "TestNewSession"

	_appVersion := uint64(1337)
	_sid := "1"
	_stype := "login"
	_idSrc := "google"
	_appId := "test"
	_userId := "btnguyen2k"
	_sdata := "My sess data"
	_delta := 5 * time.Minute
	sess := NewSession(_appVersion, "", _stype, _idSrc, _appId, _userId, _sdata, time.Now().Add(_delta))
	if sess == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	if sess.GetId() == "" {
		t.Fatalf("%s failed: empty id", testName)
	}

	for _, timeRoundingSetting := range allTimeRoundingSettings {
		t.Run(fmt.Sprintf("%#v", timeRoundingSetting), func(t *testing.T) {
			teardownTest := setupTest(t, testName, func(t *testing.T, testName string) {
				bo.UboTimestampRounding = timeRoundingSetting
			}, nil)
			defer teardownTest(t)

			_now := time.Now()
			sess := NewSession(_appVersion, _sid, _stype, _idSrc, _appId, _userId, _sdata, _now.Add(_delta))
			if sess == nil {
				t.Fatalf("%s failed: nil", testName)
			}
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
			if f, v, expected := "expiry", sess.GetExpiry(), sess.RoundTimestamp(_now.Add(_delta)); v.UnixNano() != expected.UnixNano() {
				t.Fatalf("%s failed: expected %s to be %s but received %s", testName, f, expected, v)
			}
			if f, v, expected := "sess-data", sess.GetSessionData(), _sdata; v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
		})
	}
}

func TestNewSessionFromUbo(t *testing.T) {
	testName := "TestNewSessionFromUbo"
	if sess := NewSessionFromUbo(nil); sess != nil {
		t.Fatalf("%s failed: expected nil but received %#v", testName, sess)
	}

	_appVersion := uint64(1337)
	_sid := "1"
	_stype := "login"
	_idSrc := "google"
	_appId := "test"
	_userId := "btnguyen2k"
	_sdata := "My session data"
	_delta := 5 * time.Minute
	ubo := henge.NewUniversalBo(_sid, _appVersion)
	if sess := NewSessionFromUbo(ubo); sess == nil {
		t.Fatalf("%s failed: nil", testName)
	}

	ubo.SetDataJson("invalid json string")
	if sess := NewSessionFromUbo(ubo); sess == nil {
		t.Fatalf("%s failed: nil", testName)
	}

	for _, timeRoundingSetting := range allTimeRoundingSettings {
		t.Run(fmt.Sprintf("%#v", timeRoundingSetting), func(t *testing.T) {
			teardownTest := setupTest(t, testName, func(t *testing.T, testName string) {
				bo.UboTimestampRounding = timeRoundingSetting
			}, nil)
			defer teardownTest(t)

			_now := time.Now()
			ubo.SetExtraAttr(FieldSessionExpiry, _now.Add(_delta))
			ubo.SetExtraAttr(FieldSessionSessionType, _stype)
			ubo.SetExtraAttr(FieldSessionIdSource, _idSrc)
			ubo.SetExtraAttr(FieldSessionAppId, _appId)
			ubo.SetExtraAttr(FieldSessionUserId, _userId)
			ubo.SetDataAttr(AttrSessionData, _sdata)
			sess := NewSessionFromUbo(ubo)
			if sess == nil {
				t.Fatalf("%s failed: nil", testName)
			}
			if f, v, expected := "app-version", sess.GetTagVersion(), _appVersion; v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "id", sess.GetId(), _sid; v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "session-type", sess.GetSessionType(), _stype; v != expected {
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
			if f, v, expected := "expiry", sess.GetExpiry(), sess.RoundTimestamp(_now.Add(_delta)); v.UnixNano() != expected.UnixNano() {
				t.Fatalf("%s failed: expected %s to be %s but received %s", testName, f, expected, v)
			}
			if f, v, expected := "session-data", sess.GetSessionData(), _sdata; v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
		})
	}
}

func TestSession_json(t *testing.T) {
	testName := "TestSession_json"

	for _, timeRoundingSetting := range allTimeRoundingSettings {
		t.Run(fmt.Sprintf("%#v", timeRoundingSetting), func(t *testing.T) {
			teardownTest := setupTest(t, testName, func(t *testing.T, testName string) {
				bo.UboTimestampRounding = timeRoundingSetting
			}, nil)
			defer teardownTest(t)

			_now := time.Now()
			_delta := 5 * time.Minute
			_appVersion := uint64(1337)
			_sid := "1"
			_stype := "login"
			_idSrc := "google"
			_appId := "test"
			_userId := "btnguyen2k"
			_sdata := "My session data"
			sess1 := NewSession(_appVersion, _sid, _stype, _idSrc, _appId, _userId, _sdata, _now.Add(_delta))

			js1, _ := json.Marshal(sess1)

			var sess2 *Session
			err := json.Unmarshal(js1, &sess2)
			if err != nil {
				t.Fatalf("%s failed: %e", testName, err)
			}
			if f, v, expected := "app-version", sess2.GetTagVersion(), sess1.GetTagVersion(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "id", sess2.GetId(), sess1.GetId(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "id-source", sess2.GetIdSource(), sess1.GetIdSource(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "user-id", sess2.GetUserId(), sess1.GetUserId(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "app-id", sess2.GetAppId(), sess1.GetAppId(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "session-type", sess2.GetSessionType(), sess1.GetSessionType(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "session-data", sess2.GetSessionData(), sess1.GetSessionData(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
			if f, v, expected := "expiry", sess2.GetExpiry(), sess1.GetExpiry(); v.UnixNano() != expected.UnixNano() {
				t.Fatalf("%s failed: expected %s to be %s but received %s", testName, f, expected, v)
			}
			if f, v, expected := "checksum", sess2.GetChecksum(), sess1.GetChecksum(); v != expected {
				t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
			}
		})
	}
}

func TestSession_IsExpired(t *testing.T) {
	testName := "TestSession_IsExpired"
	now := time.Now()
	sess := NewSession(1357, "1", "login", "google", "test", "btnguyen2k", "My session data", now.Add(-5*time.Minute))
	if !sess.IsExpired() {
		t.Fatalf("%s failed: session should already expire", testName)
	}
}

package session

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
)

func TestNewSession(t *testing.T) {
	name := "TestNewSession"
	now := time.Now()
	session := NewSession(1357, "1", "login", "google", "test", "btnguyen2k", "My session data", now.Add(5*time.Minute))
	if session == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := session.GetTagVersion(); tagVersion != 1357 {
		t.Fatalf("%s failed: expected app-version to be %#v but received %#v", name, 1357, tagVersion)
	}
	if id := session.GetId(); id != "1" {
		t.Fatalf("%s failed: id to be %#v but received %#v", name, "id", id)
	}
	if sessionType := session.GetSessionType(); sessionType != "login" {
		t.Fatalf("%s failed: expected session-type to be %#v but received %#v", name, "login", sessionType)
	}
	if idSource := session.GetIdSource(); idSource != "google" {
		t.Fatalf("%s failed: expected id-source to be %#v but received %#v", name, "test", idSource)
	}
	if appId := session.GetAppId(); appId != "test" {
		t.Fatalf("%s failed: expected app-id to be %#v but received %#v", name, "test", appId)
	}
	if userId := session.GetUserId(); userId != "btnguyen2k" {
		t.Fatalf("%s failed: expected user-id to be %#v but received %#v", name, "btnguyen2k", userId)
	}
	if data := session.GetSessionData(); data != "My session data" {
		t.Fatalf("%s failed: expected session-data to be %#v but received %#v", name, "My session data", data)
	}
}

func TestSession_json(t *testing.T) {
	name := "TestSession_json"
	now := time.Now()
	sess1 := NewSession(1357, "1", "login", "google", "test", "btnguyen2k", "My session data", now.Add(5*time.Minute))

	js1, _ := json.Marshal(sess1)

	var sess2 *Session
	err := json.Unmarshal(js1, &sess2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if sess1.GetId() != sess2.GetId() {
		t.Fatalf("%s failed [Id]: expected %#v but received %#v", name, sess1.GetId(), sess2.GetId())
	}
	if sess1.GetTagVersion() != sess2.GetTagVersion() {
		t.Fatalf("%s failed [AppVersion]: expected %#v but received %#v", name, sess1.GetTagVersion(), sess2.GetTagVersion())
	}
	if sess1.GetIdSource() != sess2.GetIdSource() {
		t.Fatalf("%s failed [IdSource]: expected %#v but received %#v", name, sess1.GetIdSource(), sess2.GetIdSource())
	}
	if sess1.GetUserId() != sess2.GetUserId() {
		t.Fatalf("%s failed [UserId]: expected %#v but received %#v", name, sess1.GetUserId(), sess2.GetUserId())
	}
	if sess1.GetAppId() != sess2.GetAppId() {
		t.Fatalf("%s failed [AppId]: expected %#v but received %#v", name, sess1.GetAppId(), sess2.GetAppId())
	}
	if sess1.GetSessionType() != sess2.GetSessionType() {
		t.Fatalf("%s failed [SessionType]: expected %#v but received %#v", name, sess1.GetSessionType(), sess2.GetSessionType())
	}
	if sess1.GetSessionData() != sess2.GetSessionData() {
		t.Fatalf("%s failed [SessionData]: expected %#v but received %#v", name, sess1.GetSessionData(), sess2.GetSessionData())
	}
	if sess1.GetExpiry().Format(henge.TimeLayout) != sess2.GetExpiry().Format(henge.TimeLayout) {
		t.Fatalf("%s failed [Expiry]: expected %#v but received %#v", name, sess1.GetExpiry(), sess2.GetExpiry())
	}
	if sess1.GetChecksum() != sess2.GetChecksum() {
		t.Fatalf("%s failed [Checksum]: expected %#v but received %#v", name, sess1.GetChecksum(), sess2.GetChecksum())
	}
}

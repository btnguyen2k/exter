package user

import (
	"encoding/json"
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "TestNewUser"
	user := NewUser(1357, "btnguyen2k")
	if user == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := user.GetTagVersion(); tagVersion != 1357 {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, 1357, tagVersion)
	}
	if id := user.GetId(); id != "btnguyen2k" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "id", id)
	}
}

func TestUser_json(t *testing.T) {
	name := "TestUser_json"

	user1 := NewUser(1357, "btnguyen2k")
	js1, _ := json.Marshal(user1)

	var user2 *User
	err := json.Unmarshal(js1, &user2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if user1.GetId() != user2.GetId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetId(), user2.GetId())
	}
	if user1.GetTagVersion() != user2.GetTagVersion() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetTagVersion(), user2.GetTagVersion())
	}
	if user1.GetAesKey() != user2.GetAesKey() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetAesKey(), user1.GetAesKey())
	}
	if user1.GetChecksum() != user2.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetChecksum(), user2.GetChecksum())
	}
}

package user

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/btnguyen2k/henge"
)

func TestNewUser(t *testing.T) {
	name := "TestNewUser"
	user := NewUser(1357, "myid")
	if user == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := user.GetTagVersion(); tagVersion != 1357 {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, 1357, tagVersion)
	}
	if id := user.GetId(); id != "myid" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "id", id)
	}
	if user.GetAesKey() == "" {
		t.Fatalf("%s failed: empty AES key", name)
	}
	for _, newAesKey := range []string{"  0123456789abcdef ", " abcdef0123456789   "} {
		user.SetAesKey(newAesKey)
		if aesKey, expected := user.GetAesKey(), strings.TrimSpace(newAesKey); aesKey != expected {
			t.Fatalf("%s failed: expected AES-key to be %#v but received %#v", name, expected, aesKey)
		}
	}
	for _, newDisplayName := range []string{"  My   name   ", "   Display name   "} {
		user.SetDisplayName(newDisplayName)
		if displayName, expected := user.GetDisplayName(), strings.TrimSpace(newDisplayName); displayName != expected {
			t.Fatalf("%s failed: expected AES-key to be %#v but received %#v", name, expected, displayName)
		}
	}
}

func TestNewUserFromUbo(t *testing.T) {
	name := "TestNewUserFromUbo"
	if user := NewUserFromUbo(nil); user != nil {
		t.Fatalf("%s failed: expected nil but received %#v", name, user)
	}

	ubo := henge.NewUniversalBo("myid", 1337)
	if user := NewUserFromUbo(ubo); user == nil {
		t.Fatalf("%s failed: nil", name)
	}

	ubo.SetDataJson("invalid json string")
	for _, newAesKey := range []string{"  0123456789abcdef ", " abcdef0123456789   "} {
		ubo.SetDataAttr(AttrUserAesKey, newAesKey)
		for _, newDisplayName := range []string{"  My   name   ", "   Display name   "} {
			ubo.SetDataAttr(AttrUserDisplayName, newDisplayName)

			if user := NewUserFromUbo(ubo); user == nil {
				t.Fatalf("%s failed: nil", name)
			} else {
				if displayName, expected := user.GetDisplayName(), strings.TrimSpace(newDisplayName); displayName != expected {
					t.Fatalf("%s failed: expected AES-key to be %#v but received %#v", name, expected, displayName)
				}

				if aesKey, expected := user.GetAesKey(), strings.TrimSpace(newAesKey); aesKey != expected {
					t.Fatalf("%s failed: expected AES-key to be %#v but received %#v", name, expected, aesKey)
				}
			}
		}
	}
}

func TestUser_json(t *testing.T) {
	name := "TestUser_json"

	user1 := NewUser(1357, "myid")
	for _, newAesKey := range []string{"  0123456789abcdef ", " abcdef0123456789   "} {
		user1.SetAesKey(newAesKey)
		for _, newDisplayName := range []string{"  My   name   ", "   Display name   "} {
			user1.SetDisplayName(newDisplayName)

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
	}
}

package app

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewApp(t *testing.T) {
	name := "TestNewApp"
	app := NewApp(1357, "test", "btnguyen2k", "My test application")
	if app == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := app.GetTagVersion(); tagVersion != 1357 {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, 1357, tagVersion)
	}
	if id := app.GetId(); id != "test" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "id", id)
	}
	if ownerId := app.GetOwnerId(); ownerId != "btnguyen2k" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "btnguyen2k", ownerId)
	}
	if desc := app.GetAttrsPublic().Description; desc != "My test application" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "My test application", desc)
	}
}

func TestApp_json(t *testing.T) {
	name := "TestApp_json"
	app1 := NewApp(1357, "test", "btnguyen2k", "My test application")
	attrs := app1.GetAttrsPublic()
	attrs.DefaultReturnUrl = "http://localhost/login?token="
	attrs.Tags = []string{"social", "internal"}
	attrs.IdentitySources = map[string]bool{"facebook": true, "google": false}
	attrs.RsaPublicKey = "RSA Public Key"
	app1.SetAttrsPublic(attrs)

	js1, _ := json.Marshal(app1)

	var app2 *App
	err := json.Unmarshal(js1, &app2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if app1.GetId() != app2.GetId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, app1.GetId(), app2.GetId())
	}
	if app1.GetTagVersion() != app2.GetTagVersion() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, app1.GetTagVersion(), app2.GetTagVersion())
	}
	if app1.GetOwnerId() != app2.GetOwnerId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, app1.ownerId, app2.ownerId)
	}
	if app1.GetChecksum() != app2.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, app1.GetChecksum(), app2.GetChecksum())
	}
	if !reflect.DeepEqual(app1.attrsPublic, app2.attrsPublic) {
		t.Fatalf("%s failed:\nexpected %#v\nbut received %#v", name, app1.attrsPublic, app2.attrsPublic)
	}
}

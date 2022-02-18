package app

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/btnguyen2k/henge"
)

func TestNewApp(t *testing.T) {
	testName := "TestNewApp"
	_appVersion := uint64(1337)
	_aid := "test"
	_oid := "system"
	_desc := "My test application"
	app := NewApp(_appVersion, _aid, _oid, _desc)
	if app == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	if f, v, expected := "app-version", app.GetTagVersion(), _appVersion; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "id", app.GetId(), _aid; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "owner-id", app.GetOwnerId(), _oid; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/description", app.GetAttrsPublic().Description, _desc; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
}

func TestNewAppFromUbo(t *testing.T) {
	testName := "TestNewAppFromUbo"
	if app := NewAppFromUbo(nil); app != nil {
		t.Fatalf("%s failed: expected nil but received %#v", testName, app)
	}

	_appVersion := uint64(1337)
	_aid := "test"
	_oid := "system"
	_desc := "My test application"
	_isAtive := true
	_rurl := "http://default_return_url"
	_curl := "http://default_cancel_url"
	_rsaPubKey := "rsa pub key"
	_idstr := map[string]bool{"src1": true, "src2": false}
	_tags := []string{"tag1", "tag2", "tag3"}
	ubo := henge.NewUniversalBo(_aid, _appVersion)
	ubo.SetDataJson("invalid json string")
	if app := NewAppFromUbo(ubo); app == nil {
		t.Fatalf("%s failed: nil", testName)
	}

	ubo.SetExtraAttr(FieldAppOwnerId, _oid)
	ubo.SetDataAttr(AttrAppPublicAttrs, AppAttrsPublic{
		IsActive:         _isAtive,
		Description:      _desc,
		DefaultReturnUrl: _rurl,
		DefaultCancelUrl: _curl,
		IdentitySources:  _idstr,
		Tags:             _tags,
		RsaPublicKey:     _rsaPubKey,
	})
	app := NewAppFromUbo(ubo)
	if app == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	if f, v, expected := "app-version", app.GetTagVersion(), _appVersion; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "id", app.GetId(), _aid; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "owner-id", app.GetOwnerId(), _oid; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/is-active", app.GetAttrsPublic().IsActive, _isAtive; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/description", app.GetAttrsPublic().Description, _desc; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/default-return-url", app.GetAttrsPublic().DefaultReturnUrl, _rurl; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/default-cancel-url", app.GetAttrsPublic().DefaultCancelUrl, _curl; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/id-sources", app.GetAttrsPublic().IdentitySources, _idstr; !reflect.DeepEqual(v, expected) {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/tags", app.GetAttrsPublic().Tags, _tags; !reflect.DeepEqual(v, expected) {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
}

func TestApp_json(t *testing.T) {
	testName := "TestApp_json"

	_appVersion := uint64(1337)
	_aid := "test"
	_oid := "system"
	_desc := "My test application"
	_isAtive := true
	_rurl := "http://default_return_url"
	_curl := "http://default_cancel_url"
	_rsaPubKey := "rsa pub key"
	_idstr := map[string]bool{"src1": true, "src2": false}
	_tags := []string{"tag1", "tag2", "tag3"}
	app1 := NewApp(_appVersion, _aid, _oid, _desc)
	attrs := app1.GetAttrsPublic()
	attrs.IsActive = _isAtive
	attrs.DefaultReturnUrl = _rurl
	attrs.DefaultCancelUrl = _curl
	attrs.Tags = _tags
	attrs.IdentitySources = _idstr
	attrs.RsaPublicKey = _rsaPubKey
	app1.SetAttrsPublic(attrs)

	js1, _ := json.Marshal(app1)

	var app2 *App
	err := json.Unmarshal(js1, &app2)
	if err != nil {
		t.Fatalf("%s failed: %e", testName, err)
	}

	if f, v, expected := "app-version", app2.GetTagVersion(), _appVersion; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "id", app2.GetId(), _aid; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "owner-id", app2.GetOwnerId(), _oid; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/is-active", app2.GetAttrsPublic().IsActive, _isAtive; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/description", app2.GetAttrsPublic().Description, _desc; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/default-return-url", app2.GetAttrsPublic().DefaultReturnUrl, _rurl; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/default-cancel-url", app2.GetAttrsPublic().DefaultCancelUrl, _curl; v != expected {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/id-sources", app2.GetAttrsPublic().IdentitySources, _idstr; !reflect.DeepEqual(v, expected) {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}
	if f, v, expected := "public-attrs/tags", app2.GetAttrsPublic().Tags, _tags; !reflect.DeepEqual(v, expected) {
		t.Fatalf("%s failed: expected %s to be %#v but received %#v", testName, f, expected, v)
	}

	if app1.GetChecksum() != app2.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, app1.GetChecksum(), app2.GetChecksum())
	}
	if !reflect.DeepEqual(app1.attrsPublic, app2.attrsPublic) {
		t.Fatalf("%s failed:\nexpected %#v\nbut received %#v", testName, app1.attrsPublic, app2.attrsPublic)
	}
}

func TestApp_GenerateReturnUrl(t *testing.T) {
	testName := "TestApp_GenerateReturnUrl"
	app := NewApp(0, "appid", "ownerid", "test app")

	if url := app.GenerateReturnUrl(""); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}

	if url := app.GenerateReturnUrl("in%20valid://invalid"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}

	app.SetAttrsPublic(AppAttrsPublic{DefaultReturnUrl: "in%20valid://invalid"})
	if url := app.GenerateReturnUrl("url://whatever"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}

	app.SetAttrsPublic(AppAttrsPublic{DefaultReturnUrl: "/login?src=exter"})
	if url := app.GenerateReturnUrl("url://absolute/url"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}
	if url := app.GenerateReturnUrl("/relative/url?param=value"); url != "/relative/url?param=value" {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, "/relative/url?param=value", url)
	}

	app.SetAttrsPublic(AppAttrsPublic{DefaultReturnUrl: "url://domain/path1/subpath1?src=exter"})
	if url, e := app.GenerateReturnUrl("/another/path?param=value"), "url://domain/another/path?param=value"; url != e {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, e, url)
	}
	if url := app.GenerateReturnUrl("url://another_domain/url"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}
	if url, e := app.GenerateReturnUrl("url://domain/path2/subpath2?param=value"), "url://domain/path2/subpath2?param=value"; url != e {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, e, url)
	}
}

func TestApp_GenerateCancelUrl(t *testing.T) {
	testName := "TestApp_GenerateCancelUrl"
	app := NewApp(0, "appid", "ownerid", "test app")

	if url := app.GenerateCancelUrl(""); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}

	if url := app.GenerateCancelUrl("in%20valid://invalid"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}

	app.SetAttrsPublic(AppAttrsPublic{DefaultCancelUrl: "in%20valid://invalid"})
	if url := app.GenerateCancelUrl("url://whatever"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}

	app.SetAttrsPublic(AppAttrsPublic{DefaultCancelUrl: "/login?src=exter"})
	if url := app.GenerateCancelUrl("url://absolute/url"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}
	if url, e := app.GenerateCancelUrl("/relative/url?param=value"), "/relative/url?param=value"; url != e {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, e, url)
	}

	app.SetAttrsPublic(AppAttrsPublic{DefaultCancelUrl: "url://domain/path1/subpath1?src=exter"})
	if url, e := app.GenerateCancelUrl("/path2/subpath2?param=value"), "url://domain/path2/subpath2?param=value"; url != e {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, e, url)
	}
	if url := app.GenerateCancelUrl("url://another_domain/url"); url != "" {
		t.Fatalf("%s failed: expected empty but received %#v", testName, url)
	}
	if url, e := app.GenerateCancelUrl("url://domain/path2/subpath2?param=value"), "url://domain/path2/subpath2?param=value"; url != e {
		t.Fatalf("%s failed: expected %#v but received %#v", testName, e, url)
	}
}

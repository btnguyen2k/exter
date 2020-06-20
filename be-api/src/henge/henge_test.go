package henge

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewUniversalBo(t *testing.T) {
	name := "TestNewUniversalBo"
	ubo := NewUniversalBo("id", 1357)
	if ubo == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if id := ubo.GetId(); id != "id" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "id", id)
	}
	if appVersion := ubo.GetAppVersion(); appVersion != 1357 {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, 1357, appVersion)
	}
}

func TestRowMapper(t *testing.T) {
	name := "TestRowMapper"
	tableName := "test_user"
	extraColNameToFieldMappings := map[string]string{"zuid": "owner_id"}
	rowMapper := buildRowMapper(tableName, extraColNameToFieldMappings)

	myColList := rowMapper.ColumnsList(tableName)
	expectedColList := append(columnNames, "zuid")
	if !reflect.DeepEqual(myColList, expectedColList) {
		t.Fatalf("%s failed: expected column list %#v but received %#v", name, expectedColList, myColList)
	}
}

func TestUniversalBo_json(t *testing.T) {
	name := "TestUniversalBo_json"
	ubo1 := NewUniversalBo("id1", 1357)
	ubo1.SetDataAttr("key", "value")
	ubo1.SetExtraAttr("exkey", "value")
	js1, _ := json.Marshal(ubo1)

	var ubo2 *UniversalBo
	err := json.Unmarshal(js1, &ubo2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if ubo1.id != ubo2.id {
		t.Fatalf("%s failed: expected %#v but received %#v", name, ubo1, ubo2)
	}
	if ubo1.appVersion != ubo2.appVersion {
		t.Fatalf("%s failed: expected %#v but received %#v", name, ubo1, ubo2)
	}
	if !reflect.DeepEqual(ubo1._data, ubo2._data) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, ubo1, ubo2)
	}
	if ubo1.checksum != ubo2.checksum {
		t.Fatalf("%s failed: expected %#v but received %#v", name, ubo1, ubo2)
	}
}

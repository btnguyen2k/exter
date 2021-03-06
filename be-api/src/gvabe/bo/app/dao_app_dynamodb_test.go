package app

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo/user"
)

func _inSlide(item string, slide []string) bool {
	for _, s := range slide {
		if item == s {
			return true
		}
	}
	return false
}

func _deleteTableWithWait(t *testing.T, testName string, adc *prom.AwsDynamodbConnect, table string) {
	statusList := []string{""}
	for {
		err := adc.DeleteTable(nil, table)
		if err = prom.AwsIgnoreErrorIfMatched(err, awsdynamodb.ErrCodeTableNotFoundException); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if status, err := adc.GetTableStatus(nil, table); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		} else if _inSlide(status, statusList) {
			return
		}
	}
}

func _createTableWithWait(t *testing.T, testName string, adc *prom.AwsDynamodbConnect, table string, spec *henge.DynamodbTablesSpec) {
	statusList := []string{"ACTIVE"}
	for {
		if status, err := adc.GetTableStatus(nil, table); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		} else if _inSlide(status, statusList) {
			return
		}
		err := henge.InitDynamodbTables(adc, table, spec)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}
}

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

const tableNameDynamodb = "exter_test_app"

func TestNewAppDaoAwsDynamodb(t *testing.T) {
	name := "TestNewAppDaoAwsDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	appDao := NewAppDaoAwsDynamodb(adc, tableNameDynamodb)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initAppDaoDynamodb(t *testing.T, testName string, adc *prom.AwsDynamodbConnect) AppDao {
	_deleteTableWithWait(t, testName, adc, tableNameDynamodb)
	_createTableWithWait(t, testName, adc, tableNameDynamodb, &henge.DynamodbTablesSpec{
		MainTableRcu:    2,
		MainTableWcu:    2,
		CreateUidxTable: true,
		UidxTableRcu:    2,
		UidxTableWcu:    2,
	})
	return NewAppDaoAwsDynamodb(adc, tableNameDynamodb)
}

func TestAppDaoAwsDynamodb_Create(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_Create"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoDynamodb(t, name, adc)

	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	ok, err := appDao.Create(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	items, err := adc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
}

func TestAppDaoAwsDynamodb_Get(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoDynamodb(t, name, adc)

	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	if app, err := appDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", name, "not_found")
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := app.GetTagVersion(); v != 1357 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
		}
		if v := app.GetOwnerId(); v != "btnguyen2k" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
		}
		if v := app.GetAttrsPublic().Description; v != "System application (do not delete)" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "System application (do not delete)", v)
		}
	}
}

func TestAppDaoAwsDynamodb_Delete(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_Delete"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoDynamodb(t, name, adc)

	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	app, err := appDao.Get("exter")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	}

	ok, err := appDao.Delete(app)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if !ok {
		t.Fatalf("%s failed: cannot delete app [%s]", name, app.GetId())
	}

	app, err = appDao.Get("exter")
	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", name, "exter")
	}

	items, err := adc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 0 item inserted but received %#v", name, len(items))
	}
}

func TestAppDaoAwsDynamodb_Update(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_Update"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoDynamodb(t, name, adc)

	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	appDao.Create(app)

	app.SetOwnerId("nbthanh")
	app.SetTagVersion(2468)
	app.attrsPublic.Description = "App description"
	ok, err := appDao.Update(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := app.GetTagVersion(); v != 2468 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
		}
		if v := app.GetOwnerId(); v != "nbthanh" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
		}
		if v := app.GetAttrsPublic().Description; v != "App description" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "App description", v)
		}
	}

	items, err := adc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 1 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
}

func TestAppDaoAwsDynamodb_GetUserApps(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_GetUserApps"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	appDao := _initAppDaoDynamodb(t, name, adc)

	for i := 0; i < 10; i++ {
		app := NewApp(uint64(i), strconv.Itoa(i), strconv.Itoa(i%3), "App #"+strconv.Itoa(i))
		appDao.Create(app)
	}

	u := user.NewUser(123, "2")
	appList, err := appDao.GetUserApps(u)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(appList) != 3 {
		t.Fatalf("%s failed: expected %#v apps but received %#v", name, 3, len(appList))
	}
	for _, app := range appList {
		if app.GetOwnerId() != "2" {
			t.Fatalf("%s failed: app %#v does not belong to user %#v", name, app.GetId(), "2")
		}
	}

	items, err := adc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 10 {
		t.Fatalf("%s failed: expected 10 items inserted but received %#v", name, len(items))
	}
}

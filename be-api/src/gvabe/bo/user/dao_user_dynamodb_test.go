package user

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

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

const tableNameDynamodb = "exter_test_user"

func TestNewUserDaoAwsDynamodb(t *testing.T) {
	name := "TestNewUserDaoAwsDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	userDao := NewUserDaoAwsDynamodb(adc, tableNameDynamodb)
	if userDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initUserDaoDynamodb(t *testing.T, testName string, adc *prom.AwsDynamodbConnect) UserDao {
	adc.DeleteTable(nil, tableNameDynamodb)
	henge.InitDynamodbTable(adc, tableNameDynamodb, 2, 2)
	return NewUserDaoAwsDynamodb(adc, tableNameDynamodb)
}

func TestUserDaoAwsDynamodb_Create(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_Create"
	adc := _createAwsDynamodbConnect(t, name)
	userDao := _initUserDaoDynamodb(t, name, adc)
	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
}

func TestUserDaoAwsDynamodb_Get(t *testing.T) {
	name := "TestUserDaoAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	userDao := _initUserDaoDynamodb(t, name, adc)
	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
	if u, err := userDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if u != nil {
		t.Fatalf("%s failed: user %s should not exist", name, "not_found")
	}

	if u, err := userDao.Get("btnguyen2k"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if u == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := u.GetId(); v != "btnguyen2k" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
		}
		if v := u.GetTagVersion(); v != 1357 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
		}
		if v := u.GetDisplayName(); v != "Thanh Nguyen" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "Thanh Nguyen", v)
		}
		if v := u.GetAesKey(); v != "aeskey" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "aeskey", v)
		}
	}
}

// func TestAppDaoAwsDynamodb_Delete(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_Delete"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
//
// 	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
// 	app, err := appDao.Get("exter")
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app == nil {
// 		t.Fatalf("%s failed: nil", name)
// 	}
//
// 	ok, err := appDao.Delete(app)
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if !ok {
// 		t.Fatalf("%s failed: cannot delete app [%s]", name, app.GetId())
// 	}
//
// 	app, err = appDao.Get("exter")
// 	if app, err := appDao.Get("exter"); err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app != nil {
// 		t.Fatalf("%s failed: app %s should not exist", name, "exter")
// 	}
// }
//
// func TestAppDaoAwsDynamodb_Update(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_Update"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
//
// 	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
// 	appDao.Create(app)
//
// 	app.SetOwnerId("nbthanh")
// 	app.SetTagVersion(2468)
// 	app.attrsPublic.Description = "App description"
// 	ok, err := appDao.Update(app)
// 	if err != nil || !ok {
// 		t.Fatalf("%s failed: %#v / %s", name, ok, err)
// 	}
//
// 	if app, err := appDao.Get("exter"); err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	} else if app == nil {
// 		t.Fatalf("%s failed: nil", name)
// 	} else {
// 		if v := app.GetId(); v != "exter" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
// 		}
// 		if v := app.GetTagVersion(); v != 2468 {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
// 		}
// 		if v := app.GetOwnerId(); v != "nbthanh" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
// 		}
// 		if v := app.GetAttrsPublic().Description; v != "App description" {
// 			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "App description", v)
// 		}
// 	}
// }
//
// func TestAppDaoAwsDynamodb_GetUserApps(t *testing.T) {
// 	name := "TestAppDaoAwsDynamodb_GetUserApps"
// 	adc := _createAwsDynamodbConnect(t, name)
// 	appDao := _initAppDaoDynamodb(t, name, adc)
//
// 	for i := 0; i < 10; i++ {
// 		app := NewApp(uint64(i), strconv.Itoa(i), strconv.Itoa(i%3), "App #"+strconv.Itoa(i))
// 		appDao.Create(app)
// 	}
//
// 	u := user.NewUser(123, "2")
// 	appList, err := appDao.GetUserApps(u)
// 	if err != nil {
// 		t.Fatalf("%s failed: %s", name, err)
// 	}
// 	if len(appList) != 3 {
// 		t.Fatalf("%s failed: expected %#v apps but received %#v", name, 3, len(appList))
// 	}
// 	for _, app := range appList {
// 		if app.GetOwnerId() != "2" {
// 			t.Fatalf("%s failed: app %#v does not belong to user %#v", name, app.GetId(), "2")
// 		}
// 	}
// }

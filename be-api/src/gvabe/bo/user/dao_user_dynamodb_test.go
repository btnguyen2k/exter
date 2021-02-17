package user

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
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
		if status, err := adc.GetTableStatus(nil, table); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		} else if _inSlide(status, statusList) {
			return
		}
		err := adc.DeleteTable(nil, table)
		if err = prom.AwsIgnoreErrorIfMatched(err, awsdynamodb.ErrCodeTableNotFoundException); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
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
	_deleteTableWithWait(t, testName, adc, tableNameDynamodb)
	_createTableWithWait(t, testName, adc, tableNameDynamodb, &henge.DynamodbTablesSpec{
		MainTableRcu:    2,
		MainTableWcu:    2,
		CreateUidxTable: true,
		UidxTableRcu:    2,
		UidxTableWcu:    2,
	})
	return NewUserDaoAwsDynamodb(adc, tableNameDynamodb)
}

func TestUserDaoAwsDynamodb_Create(t *testing.T) {
	name := "TestAppDaoAwsDynamodb_Create"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoDynamodb(t, name, adc)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
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

func TestUserDaoAwsDynamodb_Get(t *testing.T) {
	name := "TestUserDaoAwsDynamodb_Get"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
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

func TestUserDaoAwsDynamodb_Delete(t *testing.T) {
	name := "TestUserDaoAwsDynamodb_Delete"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoDynamodb(t, name, adc)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	ok, err := userDao.Create(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	ok, err = userDao.Delete(u)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if !ok {
		t.Fatalf("%s failed: cannot delete user [%s]", name, u.GetId())
	}

	u, err = userDao.Get("btnguyen2k")
	if app, err := userDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: user %s should not exist", name, "userDao")
	}

	items, err := adc.ScanItems(nil, tableNameDynamodb, nil, "")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(items) != 0 {
		t.Fatalf("%s failed: expected 1 item inserted but received %#v", name, len(items))
	}
}

func TestUserDaoAwsDynamodb_Update(t *testing.T) {
	name := "TestUserDaoAwsDynamodb_Update"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	userDao := _initUserDaoDynamodb(t, name, adc)

	u := NewUser(1357, "btnguyen2k").SetDisplayName("Thanh Nguyen").SetAesKey("aeskey")
	userDao.Create(u)

	u.SetDisplayName("nbthanh")
	u.SetAesKey("newaeskey")
	ok, err := userDao.Update(u)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
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
		if v := u.GetDisplayName(); v != "nbthanh" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
		}
		if v := u.GetAesKey(); v != "newaeskey" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "newaeskey", v)
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

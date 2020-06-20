package main

import (
	"fmt"
	"reflect"

	_ "github.com/jackc/pgx/v4/stdlib"

	"main/src/gvabe/bo/app"
	"main/src/henge"
)

func main() {
	url := "postgres://test:test@localhost:5432/test?sslmode=disable&client_encoding=UTF-8&application_name=govueadmin"
	timezone := "Asia/Ho_Chi_Minh"
	appId := "test"
	appOwnerId := "btnguyen2k"
	appDest := "My test application"
	sqlc := henge.NewPgsqlConnection(url, timezone)
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", app.TableApp))
	henge.InitPgsqlTable(sqlc, app.TableApp, map[string]string{app.ColApp_UserId: "VARCHAR(32)"})
	appDao := app.NewAppDaoSql(sqlc, app.TableApp)

	app1 := app.NewApp(1357, appId, appOwnerId, appDest)
	attrs := app1.GetAttrsPublic()
	attrs.DefaultReturnUrl = "http://localhost/login?token="
	attrs.Tags = []string{"social", "internal"}
	attrs.IdentitySources = map[string]bool{"facebook": true, "google": false}
	attrs.RsaPublicKey = "RSA Public Key"
	app1.SetAttrsPublic(attrs)

	result, err := appDao.Create(app1)
	fmt.Println("Create:", result, err)

	app2, err := appDao.Get(appId)
	fmt.Println("Get["+appId+"]:", err)

	if app1.GetId() != app2.GetId() {
		fmt.Printf("Failed: expected %#v but received %#v\n", app1.GetId(), app2.GetId())
	}
	if app1.GetAppVersion() != app2.GetAppVersion() {
		fmt.Printf("Failed: expected %#v but received %#v", app1.GetAppVersion(), app2.GetAppVersion())
	}
	if app1.GetOwnerId() != app2.GetOwnerId() {
		fmt.Printf("Failed: expected %#v but received %#v", app1.GetOwnerId(), app2.GetOwnerId())
	}
	if app1.GetChecksum() != app2.GetChecksum() {
		fmt.Printf("Failed:  expected %#v but received %#v", app1.GetChecksum(), app2.GetChecksum())
	}
	if !reflect.DeepEqual(app1.GetAttrsPublic(), app2.GetAttrsPublic()) {
		fmt.Printf("Failed: failed:\nexpected %#v\nbut received %#v", app1.GetAttrsPublic(), app2.GetAttrsPublic())
	}

	fmt.Println("DONE.")
}

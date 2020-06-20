package main

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"

	"main/src/gvabe/bo/user"
	"main/src/henge"
)

func main() {
	url := "postgres://test:test@localhost:5432/test?sslmode=disable&client_encoding=UTF-8&application_name=govueadmin"
	timezone := "Asia/Ho_Chi_Minh"
	userId := "btnguyen2k"
	sqlc := henge.NewPgsqlConnection(url, timezone)
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", user.TableUser))
	henge.InitPgsqlTable(sqlc, user.TableUser, nil)
	userDao := user.NewUserDaoSql(sqlc, user.TableUser)

	user1 := user.NewUser(2468, userId)

	result, err := userDao.Create(user1)
	fmt.Println("Create:", result, err)

	user2, err := userDao.Get(userId)
	fmt.Println("Get["+userId+"]:", err)

	if user1.GetId() != user2.GetId() {
		fmt.Printf("Failed: expected %#v but received %#v", user1.GetId(), user2.GetId())
	}
	if user1.GetAppVersion() != user2.GetAppVersion() {
		fmt.Printf("Failed: expected %#v but received %#v", user1.GetAppVersion(), user2.GetAppVersion())
	}
	if user1.GetAesKey() != user2.GetAesKey() {
		fmt.Printf("Failed: expected %#v but received %#v", user1.GetAesKey(), user1.GetAesKey())
	}
	if user1.GetChecksum() != user2.GetChecksum() {
		fmt.Printf("Failed: expected %#v but received %#v", user1.GetChecksum(), user2.GetChecksum())
	}

	fmt.Println("DONE.")
}

package main

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"main/src/gvabe/bo/session"
	"main/src/henge"
)

func main() {
	url := "postgres://test:test@localhost:5432/test?sslmode=disable&client_encoding=UTF-8&application_name=govueadmin"
	timezone := "Asia/Ho_Chi_Minh"
	sessId := "1"
	now := time.Now()
	sqlc := henge.NewPgsqlConnection(url, timezone)
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", session.TableSession))
	henge.InitPgsqlTable(sqlc, session.TableSession, map[string]string{
		session.ColSession_IdSource:    "VARCHAR(32)",
		session.ColSession_AppId:       "VARCHAR(32)",
		session.ColSession_UserId:      "VARCHAR(32)",
		session.ColSession_SessionType: "VARCHAR(32)",
		session.ColSession_Expiry:      "TIMESTAMP WITH TIME ZONE",
	})
	sessDao := session.NewSessionDaoSql(sqlc, session.TableSession)
	sess1 := session.NewSession(1357, "1", "login", "google", "test", "btnguyen2k", "My session data", now.Add(5*time.Minute))

	result, err := sessDao.Create(sess1)
	fmt.Println("Create:", result, err)

	sess2, err := sessDao.Get(sessId)
	fmt.Println("Get["+sessId+"]:", err)

	if sess1.GetId() != sess2.GetId() {
		fmt.Printf("Failed {GetId}: expected %#v but received %#v\n", sess1.GetId(), sess2.GetId())
	}
	if sess1.GetAppVersion() != sess2.GetAppVersion() {
		fmt.Printf("Failed {GetAppVersion}: expected %#v but received %#v\n", sess1.GetAppVersion(), sess2.GetAppVersion())
	}
	if sess1.GetIdSource() != sess2.GetIdSource() {
		fmt.Printf("Failed {GetIdSource}: expected %#v but received %#v\n", sess1.GetIdSource(), sess2.GetIdSource())
	}
	if sess1.GetUserId() != sess2.GetUserId() {
		fmt.Printf("Failed {GetUserId}: expected %#v but received %#v\n", sess1.GetUserId(), sess2.GetUserId())
	}
	if sess1.GetAppId() != sess2.GetAppId() {
		fmt.Printf("Failed {GetAppId}: expected %#v but received %#v\n", sess1.GetAppId(), sess2.GetAppId())
	}
	if sess1.GetSessionType() != sess2.GetSessionType() {
		fmt.Printf("Failed {GetSessionType}: expected %#v but received %#v\n", sess1.GetSessionType(), sess2.GetSessionType())
	}
	if sess1.GetSessionData() != sess2.GetSessionData() {
		fmt.Printf("Failed {GetSessionData}: expected %#v but received %#v\n", sess1.GetSessionData(), sess2.GetSessionData())
	}
	if sess1.GetExpiry().Format(henge.TimeLayout) != sess2.GetExpiry().Format(henge.TimeLayout) {
		fmt.Printf("Failed {GetExpiry}: expected %#v but received %#v\n", sess1.GetExpiry(), sess2.GetExpiry())
	}
	if sess1.GetChecksum() != sess2.GetChecksum() {
		fmt.Printf("Failed {GetChecksum}: expected %#v but received %#v\n", sess1.GetChecksum(), sess2.GetChecksum())
	}

	fmt.Println("DONE.")
}

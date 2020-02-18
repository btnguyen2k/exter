package bo

import (
	"fmt"
	"github.com/btnguyen2k/prom"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// NewSqliteConnection creates a new connection pool to SQLite3.
func NewSqliteConnection(dir, dbName string) *prom.SqlConnect {
	err := os.MkdirAll(dir, 0711)
	if err != nil {
		panic(err)
	}
	sqlc, err := prom.NewSqlConnect("sqlite3", dir+"/"+dbName+".db", 10000, nil)
	if err != nil {
		panic(err)
	}
	return sqlc
}

// InitSqliteTable initializes database table to store bo
func InitSqliteTable(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(255), %s TIMESTAMP, %s TIMESTAMP, %s BIGINT, PRIMARY KEY (%s))",
		tableName, ColId, ColData, ColTimeCreated, ColTimeUpdated, ColAppVersion, ColId)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

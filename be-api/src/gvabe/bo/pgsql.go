package bo

import (
	"fmt"
	"github.com/btnguyen2k/prom"
	_ "github.com/lib/pq"
	"time"
)

// NewPgsqlConnection creates a new connection pool to PostgreSQL.
func NewPgsqlConnection(url, timezone string) *prom.SqlConnect {
	driver := "postgres"
	sqlConnect, err := prom.NewSqlConnect(driver, url, 10000, nil)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation(timezone)
	sqlConnect.SetLocation(loc)
	return sqlConnect
}

// InitPgsqlTable initializes database table to store bo
func InitPgsqlTable(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s JSONB, %s TIMESTAMP WITH TIME ZONE, %s TIMESTAMP WITH TIME ZONE, %s BIGINT, PRIMARY KEY (%s))",
		tableName, ColId, ColData, ColTimeCreated, ColTimeUpdated, ColAppVersion, ColId)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

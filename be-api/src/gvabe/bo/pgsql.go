package bo

import (
	"time"

	"github.com/btnguyen2k/prom"
	_ "github.com/lib/pq"
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
func InitPgsqlTable(sqlc *prom.SqlConnect, tableName string, extraCols map[string]string) {
	colDef := map[string]string{
		ColId:          "VARCHAR(64)",
		ColData:        "JSONB",
		ColTimeCreated: "TIMESTAMP WITH TIME ZONE",
		ColTimeUpdated: "TIMESTAMP WITH TIME ZONE",
		ColAppVersion:  "BIGINT",
	}
	for k, v := range extraCols {
		colDef[k] = v
	}
	pk := []string{ColId}
	if err := CreateTable(sqlc, tableName, true, colDef, pk); err != nil {
		panic(err)
	}
}

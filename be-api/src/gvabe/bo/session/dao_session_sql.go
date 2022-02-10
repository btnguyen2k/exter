package session

import (
	"fmt"

	"github.com/btnguyen2k/prom"
	"main/src/gvabe/bo"

	"github.com/btnguyen2k/henge"
)

const (
	SqlColSessionSessionType = "ztype"
	SqlColSessionIdSource    = "zidsrc"
	SqlColSessionAppId       = "zappid"
	SqlColSessionUserId      = "zuid"
	SqlColSessionExpiry      = "zexpiry"
)

// NewSessionDaoSql is helper method to create SQL-implementation of SessionDao.
func NewSessionDaoSql(sqlc *prom.SqlConnect, tableName string) SessionDao {
	dao := &SessionDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, true, map[string]string{
		SqlColSessionSessionType: FieldSessionSessionType,
		SqlColSessionIdSource:    FieldSessionIdSource,
		SqlColSessionAppId:       FieldSessionAppId,
		SqlColSessionUserId:      FieldSessionUserId,
		SqlColSessionExpiry:      FieldSessionExpiry,
	})
	return dao
}

// InitSessionTableSql is helper function to initialize SQL-based table to store sessions.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitSessionTableSql(sqlc *prom.SqlConnect, tableName string) error {
	switch sqlc.GetDbFlavor() {
	case prom.FlavorPgSql:
		return henge.InitPgsqlTable(sqlc, tableName, map[string]string{
			SqlColSessionIdSource:    "VARCHAR(32)",
			SqlColSessionAppId:       "VARCHAR(32)",
			SqlColSessionUserId:      "VARCHAR(32)",
			SqlColSessionSessionType: "VARCHAR(32)",
			SqlColSessionExpiry:      "TIMESTAMP WITH TIME ZONE",
		})
	case prom.FlavorMsSql:
		return henge.InitMssqlTable(sqlc, tableName, map[string]string{
			SqlColSessionIdSource:    "NVARCHAR(32)",
			SqlColSessionAppId:       "NVARCHAR(32)",
			SqlColSessionUserId:      "NVARCHAR(32)",
			SqlColSessionSessionType: "NVARCHAR(32)",
			SqlColSessionExpiry:      "DATETIMEOFFSET",
		})
	case prom.FlavorMySql:
		return henge.InitMysqlTable(sqlc, tableName, map[string]string{
			SqlColSessionIdSource:    "VARCHAR(32)",
			SqlColSessionAppId:       "VARCHAR(32)",
			SqlColSessionUserId:      "VARCHAR(32)",
			SqlColSessionSessionType: "VARCHAR(32)",
			SqlColSessionExpiry:      "TIMESTAMP",
		})
	case prom.FlavorOracle:
		return henge.InitOracleTable(sqlc, tableName, map[string]string{
			SqlColSessionIdSource:    "NVARCHAR2(32)",
			SqlColSessionAppId:       "NVARCHAR2(32)",
			SqlColSessionUserId:      "NVARCHAR2(32)",
			SqlColSessionSessionType: "NVARCHAR2(32)",
			SqlColSessionExpiry:      "TIMESTAMP WITH TIME ZONE",
		})
	case prom.FlavorSqlite:
		return henge.InitSqliteTable(sqlc, tableName, map[string]string{
			SqlColSessionIdSource:    "VARCHAR(32)",
			SqlColSessionAppId:       "VARCHAR(32)",
			SqlColSessionUserId:      "VARCHAR(32)",
			SqlColSessionSessionType: "VARCHAR(32)",
			SqlColSessionExpiry:      "TIMESTAMP",
		})
	case prom.FlavorCosmosDb:
		return henge.InitCosmosdbCollection(sqlc, tableName, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbPkName})
	}
	return fmt.Errorf("unsupported database type %v", sqlc.GetDbFlavor())
}

// SessionDaoSql is SQL-implementation of SessionDao.
type SessionDaoSql struct {
	henge.UniversalDao
}

// Delete implements SessionDao.Delete.
func (dao *SessionDaoSql) Delete(sess *Session) (bool, error) {
	return dao.UniversalDao.Delete(sess.UniversalBo)
}

// Get implements SessionDao.Get.
func (dao *SessionDaoSql) Get(id string) (*Session, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewSessionFromUbo(ubo), err
}

// Update implements SessionDao.Save.
func (dao *SessionDaoSql) Save(sess *Session) (bool, error) {
	ok, _, err := dao.UniversalDao.Save(sess.sync().UniversalBo)
	return ok, err
}

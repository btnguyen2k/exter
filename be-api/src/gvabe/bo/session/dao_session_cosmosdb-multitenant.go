package session

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

// NewSessionDaoMultitenantCosmosdb is helper method to create CosmosDB-implementation (multi-tenant table) of SessionDao.
func NewSessionDaoMultitenantCosmosdb(sqlc *prom.SqlConnect, tableName string) SessionDao {
	spec := &henge.CosmosdbDaoSpec{PkName: bo.CosmosdbMultitenantPkName, PkValue: bo.CosmosdbMultitenantPkValueSession, TxModeOnWrite: true}
	innerDao := SessionDaoSql{UniversalDao: henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)}
	dao := &SessionDaoCosmosdb{SessionDaoSql: innerDao, spec: spec}
	return dao
}

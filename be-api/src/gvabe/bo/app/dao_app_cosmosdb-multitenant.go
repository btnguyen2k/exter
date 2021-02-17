package app

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo"
)

// NewAppDaoMultitenantCosmosdb is helper method to create CosmosDB-implementation (multi-tenant table) of AppDao.
func NewAppDaoMultitenantCosmosdb(sqlc *prom.SqlConnect, tableName string) AppDao {
	spec := &henge.CosmosdbDaoSpec{PkName: bo.CosmosdbMultitenantPkName, PkValue: bo.CosmosdbMultitenantPkValueApp, TxModeOnWrite: true}
	innerDao := AppDaoSql{UniversalDao: henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)}
	dao := &AppDaoCosmosdb{AppDaoSql: innerDao, spec: spec}
	return dao
}

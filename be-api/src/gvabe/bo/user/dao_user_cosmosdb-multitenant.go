package user

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

// NewUserDaoMultitenantCosmosdb is helper method to create CosmosDB-implementation (multi-tenant table) of UserDao.
func NewUserDaoMultitenantCosmosdb(sqlc *prom.SqlConnect, tableName string) UserDao {
	spec := &henge.CosmosdbDaoSpec{PkName: bo.CosmosdbMultitenantPkName, PkValue: bo.CosmosdbMultitenantPkValueUser, TxModeOnWrite: true}
	innerDao := UserDaoSql{UniversalDao: henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)}
	dao := &UserDaoCosmosdb{UserDaoSql: innerDao, spec: spec}
	return dao
}

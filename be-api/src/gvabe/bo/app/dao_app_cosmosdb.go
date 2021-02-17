package app

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo"
)

// NewAppDaoCosmosdb is helper method to create CosmosDB-implementation of AppDao.
func NewAppDaoCosmosdb(sqlc *prom.SqlConnect, tableName string) AppDao {
	spec := &henge.CosmosdbDaoSpec{PkName: bo.CosmosdbPkName, TxModeOnWrite: true}
	innerDao := AppDaoSql{UniversalDao: henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)}
	dao := &AppDaoCosmosdb{AppDaoSql: innerDao, spec: spec}
	return dao
}

// AppDaoCosmosdb is CosmosDB-implementation of AppDao.
type AppDaoCosmosdb struct {
	AppDaoSql
	spec *henge.CosmosdbDaoSpec
}

// Create implements AppDao.Create.
func (dao *AppDaoCosmosdb) Create(bo *App) (bool, error) {
	ubo := bo.sync().UniversalBo
	if dao.spec != nil && dao.spec.PkName != "" && dao.spec.PkValue != "" {
		ubo.SetExtraAttr(dao.spec.PkName, dao.spec.PkValue)
	}
	return dao.UniversalDao.Create(ubo)
}

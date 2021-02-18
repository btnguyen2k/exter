package user

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

// NewUserDaoCosmosdb is helper method to create CosmosDB-implementation of UserDao.
func NewUserDaoCosmosdb(sqlc *prom.SqlConnect, tableName string) UserDao {
	spec := &henge.CosmosdbDaoSpec{PkName: bo.CosmosdbPkName, TxModeOnWrite: true}
	innerDao := UserDaoSql{UniversalDao: henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)}
	dao := &UserDaoCosmosdb{UserDaoSql: innerDao, spec: spec}
	return dao
}

// UserDaoCosmosdb is CosmosDB-implementation of SessionDao.
type UserDaoCosmosdb struct {
	UserDaoSql
	spec *henge.CosmosdbDaoSpec
}

// Create implements UserDao.Create.
func (dao *UserDaoCosmosdb) Create(bo *User) (bool, error) {
	ubo := bo.sync().UniversalBo
	if dao.spec != nil && dao.spec.PkName != "" && dao.spec.PkValue != "" {
		ubo.SetExtraAttr(dao.spec.PkName, dao.spec.PkValue)
	}
	return dao.UniversalDao.Create(ubo)
}

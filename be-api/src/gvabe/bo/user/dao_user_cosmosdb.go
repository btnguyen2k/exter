package user

import (
	"fmt"

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

// InitUserTableCosmosdb is helper function to initialize CosmosDB-based table to store users.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitUserTableCosmosdb(sqlc *prom.SqlConnect, tableName string) error {
	switch sqlc.GetDbFlavor() {
	case prom.FlavorCosmosDb:
		return InitUserTableSql(sqlc, tableName)
	}
	return fmt.Errorf("unsupported database type %v", sqlc.GetDbFlavor())
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

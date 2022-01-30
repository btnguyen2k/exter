package user

import (
	"fmt"

	"github.com/btnguyen2k/prom"
	"main/src/gvabe/bo"

	"github.com/btnguyen2k/henge"
)

// NewUserDaoSql is helper method to create SQL-implementation of UserDao.
func NewUserDaoSql(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, true, nil)
	return dao
}

// InitUserTableSql is helper function to initialize SQL-based table (collection) to store users.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitUserTableSql(sqlc *prom.SqlConnect, tableName string) error {
	switch sqlc.GetDbFlavor() {
	case prom.FlavorPgSql:
		return henge.InitPgsqlTable(sqlc, tableName, nil)
	case prom.FlavorMsSql:
		return henge.InitMssqlTable(sqlc, tableName, nil)
	case prom.FlavorMySql:
		return henge.InitMysqlTable(sqlc, tableName, nil)
	case prom.FlavorOracle:
		return henge.InitOracleTable(sqlc, tableName, nil)
	case prom.FlavorSqlite:
		return henge.InitSqliteTable(sqlc, tableName, nil)
	case prom.FlavorCosmosDb:
		return henge.InitCosmosdbCollection(sqlc, tableName, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbPkName})
	}
	return fmt.Errorf("unsupported database type %v", sqlc.GetDbFlavor())
}

// UserDaoSql is SQL-implementation of UserDao.
type UserDaoSql struct {
	henge.UniversalDao
}

// Delete implements UserDao.Delete.
func (dao *UserDaoSql) Delete(bo *User) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements UserDao.Create.
func (dao *UserDaoSql) Create(bo *User) (bool, error) {
	return dao.UniversalDao.Create(bo.sync().UniversalBo)
}

// Get implements UserDao.Get.
func (dao *UserDaoSql) Get(id string) (*User, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewUserFromUbo(ubo), err
}

// Update implements UserDao.Update.
func (dao *UserDaoSql) Update(bo *User) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

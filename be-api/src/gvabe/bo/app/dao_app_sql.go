package app

import (
	"fmt"

	"github.com/btnguyen2k/prom"
	"main/src/gvabe/bo"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo/user"
)

const (
	SqlColAppUserId = "zuid"
)

// NewAppDaoSql is helper method to create SQL-implementation of AppDao.
func NewAppDaoSql(sqlc *prom.SqlConnect, tableName string) AppDao {
	dao := &AppDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, true, map[string]string{SqlColAppUserId: FieldAppOwnerId})
	return dao
}

// InitAppTableSql is helper function to initialize SQL-based table to store application data.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitAppTableSql(sqlc *prom.SqlConnect, tableName string) error {
	switch sqlc.GetDbFlavor() {
	case prom.FlavorPgSql:
		return henge.InitPgsqlTable(sqlc, tableName, map[string]string{SqlColAppUserId: "VARCHAR(32)"})
	case prom.FlavorMsSql:
		return henge.InitMssqlTable(sqlc, tableName, map[string]string{SqlColAppUserId: "NVARCHAR(32)"})
	case prom.FlavorMySql:
		return henge.InitMysqlTable(sqlc, tableName, map[string]string{SqlColAppUserId: "VARCHAR(32)"})
	case prom.FlavorOracle:
		return henge.InitOracleTable(sqlc, tableName, map[string]string{SqlColAppUserId: "NVARCHAR2(32)"})
	case prom.FlavorSqlite:
		return henge.InitSqliteTable(sqlc, tableName, map[string]string{SqlColAppUserId: "VARCHAR(32)"})
	case prom.FlavorCosmosDb:
		return henge.InitCosmosdbCollection(sqlc, tableName, &henge.CosmosdbCollectionSpec{Pk: bo.CosmosdbPkName})
	}
	return fmt.Errorf("unsupported database type %v", sqlc.GetDbFlavor())
}

// AppDaoSql is SQL-implementation of AppDao.
type AppDaoSql struct {
	henge.UniversalDao
}

// Delete implements AppDao.Delete.
func (dao *AppDaoSql) Delete(bo *App) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements AppDao.Create.
func (dao *AppDaoSql) Create(bo *App) (bool, error) {
	return dao.UniversalDao.Create(bo.sync().UniversalBo)
}

// Get implements AppDao.Get.
func (dao *AppDaoSql) Get(id string) (*App, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewAppFromUbo(ubo), err
}

// getN implements AppDao.getN.
func (dao *AppDaoSql) getN(fromOffset, maxNumRows int) ([]*App, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, nil, nil)
	if err != nil {
		return nil, err
	}
	result := make([]*App, 0)
	for _, ubo := range uboList {
		bo := NewAppFromUbo(ubo)
		result = append(result, bo)
	}
	return result, nil
}

// getAll implements AppDao.getAll.
func (dao *AppDaoSql) getAll() ([]*App, error) {
	return dao.getN(0, 0)
}

// GetUserApps implements AppDao.GetUserApps.
func (dao *AppDaoSql) GetUserApps(u *user.User) ([]*App, error) {
	if appList, err := dao.getAll(); err != nil {
		return nil, err
	} else {
		result := make([]*App, 0)
		for _, app := range appList {
			if app.ownerId == u.GetId() {
				result = append(result, app)
			}
		}
		return result, nil
	}
}

// Update implements AppDao.Update.
func (dao *AppDaoSql) Update(bo *App) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

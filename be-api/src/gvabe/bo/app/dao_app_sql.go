package app

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo/user"
)

const (
	SqlCol_App_UserId = "zuid"
)

// NewAppDaoSql is helper method to create SQL-implementation of AppDao.
func NewAppDaoSql(sqlc *prom.SqlConnect, tableName string) AppDao {
	dao := &AppDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, true, map[string]string{SqlCol_App_UserId: FieldApp_OwnerId})
	return dao
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

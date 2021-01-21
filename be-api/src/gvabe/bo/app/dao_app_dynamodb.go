package app

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo/user"
)

// NewAppDaoAwsDynamodb is helper method to create AWS DynamoDB-implementation of AppDao.
func NewAppDaoAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) AppDao {
	dao := &AppDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, nil)}
	return dao
}

// AppDaoAwsDynamodb is AWS DynamoDB-implementation of AppDao.
type AppDaoAwsDynamodb struct {
	henge.UniversalDao
}

// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter.
// func (dao *AppDaoAwsDynamodb) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
// 	return map[string]interface{}{henge.FieldId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
// }

// Delete implements AppDao.Delete.
func (dao *AppDaoAwsDynamodb) Delete(bo *App) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements AppDao.Create.
func (dao *AppDaoAwsDynamodb) Create(bo *App) (bool, error) {
	return dao.UniversalDao.Create(bo.sync().UniversalBo)
}

// Get implements AppDao.Get.
func (dao *AppDaoAwsDynamodb) Get(id string) (*App, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewAppFromUbo(ubo), err
}

// getN implements AppDao.getN.
func (dao *AppDaoAwsDynamodb) getN(fromOffset, maxNumRows int) ([]*App, error) {
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
func (dao *AppDaoAwsDynamodb) getAll() ([]*App, error) {
	return dao.getN(0, 0)
}

// GetUserApps implements AppDao.GetUserApps.
func (dao *AppDaoAwsDynamodb) GetUserApps(u *user.User) ([]*App, error) {
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
func (dao *AppDaoAwsDynamodb) Update(bo *App) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

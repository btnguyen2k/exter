package app

import (
	"strings"

	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo/user"
)

// NewAppDaoMongo is helper method to create MongoDB-implementation of AppDao.
//
// Available: since v0.6.0
func NewAppDaoMongo(mc *prom.MongoConnect, collectionName string) AppDao {
	txMode := strings.Index(strings.ToLower(mc.GetUrl()), "replicaset=") > 0
	dao := &AppDaoMongo{UniversalDao: henge.NewUniversalDaoMongo(mc, collectionName, txMode)}
	return dao
}

// InitAppTableMongo is helper function to initialize MongoDB table (collection) to store application data.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitAppTableMongo(mc *prom.MongoConnect, collectionName string) error {
	return henge.InitMongoCollection(mc, collectionName)
}

// AppDaoMongo is MongoDB-implementation of AppDao.
//
// Available: since v0.6.0
type AppDaoMongo struct {
	henge.UniversalDao
}

// Delete implements AppDao.Delete.
func (dao *AppDaoMongo) Delete(bo *App) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements AppDao.Create.
func (dao *AppDaoMongo) Create(bo *App) (bool, error) {
	return dao.UniversalDao.Create(bo.sync().UniversalBo)
}

// Get implements AppDao.Get.
func (dao *AppDaoMongo) Get(id string) (*App, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewAppFromUbo(ubo), err
}

// getN implements AppDao.getN.
func (dao *AppDaoMongo) getN(fromOffset, maxNumRows int) ([]*App, error) {
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
func (dao *AppDaoMongo) getAll() ([]*App, error) {
	return dao.getN(0, 0)
}

// GetUserApps implements AppDao.GetUserApps.
func (dao *AppDaoMongo) GetUserApps(u *user.User) ([]*App, error) {
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
func (dao *AppDaoMongo) Update(bo *App) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

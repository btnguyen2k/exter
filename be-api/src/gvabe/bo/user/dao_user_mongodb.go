package user

import (
	"strings"

	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewUserDaoMongo is helper method to create MongoDB-implementation of UserDao.
func NewUserDaoMongo(mc *prom.MongoConnect, collectionName string) UserDao {
	txMode := strings.Index(strings.ToLower(mc.GetUrl()), "replicaset=") > 0
	dao := &UserDaoMongo{UniversalDao: henge.NewUniversalDaoMongo(mc, collectionName, txMode)}
	return dao
}

// UserDaoMongo is MongoDB-implementation of UserDao.
type UserDaoMongo struct {
	henge.UniversalDao
}

// Delete implements UserDao.Delete.
func (dao *UserDaoMongo) Delete(bo *User) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements UserDao.Create.
func (dao *UserDaoMongo) Create(bo *User) (bool, error) {
	return dao.UniversalDao.Create(bo.sync().UniversalBo)
}

// Get implements UserDao.Get.
func (dao *UserDaoMongo) Get(id string) (*User, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewUserFromUbo(ubo), err
}

// Update implements UserDao.Update.
func (dao *UserDaoMongo) Update(bo *User) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

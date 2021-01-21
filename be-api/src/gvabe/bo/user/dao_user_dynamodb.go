package user

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewUserDaoAwsDynamodb is helper method to create AWS DynamoDB-implementation of UserDao.
func NewUserDaoAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) UserDao {
	dao := &UserDaoAwsDynamodb{}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(dync, tableName, nil)
	return dao
}

// UserDaoAwsDynamodb is AWS DynamoDB-implementation of UserDao.
type UserDaoAwsDynamodb struct {
	henge.UniversalDao
}

// Delete implements UserDao.Delete.
func (dao *UserDaoAwsDynamodb) Delete(bo *User) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements UserDao.Create.
func (dao *UserDaoAwsDynamodb) Create(bo *User) (bool, error) {
	return dao.UniversalDao.Create(bo.sync().UniversalBo)
}

// Get implements UserDao.Get.
func (dao *UserDaoAwsDynamodb) Get(id string) (*User, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewUserFromUbo(ubo), err
}

// Update implements UserDao.Update.
func (dao *UserDaoAwsDynamodb) Update(bo *User) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

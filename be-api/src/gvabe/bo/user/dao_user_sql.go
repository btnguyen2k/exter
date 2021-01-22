package user

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewUserDaoSql is helper method to create SQL-implementation of UserDao.
func NewUserDaoSql(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, true, nil)
	return dao
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

package user

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
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

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter.
func (dao *UserDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.SqlColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
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
	if err != nil {
		return nil, err
	}
	return NewUserFromUbo(ubo), nil
}

// GetN implements UserDao.getN.
func (dao *UserDaoSql) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, nil, nil)
	if err != nil {
		return nil, err
	}
	result := make([]*User, 0)
	for _, ubo := range uboList {
		bo := NewUserFromUbo(ubo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements UserDao.getAll.
func (dao *UserDaoSql) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update.
func (dao *UserDaoSql) Update(bo *User) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

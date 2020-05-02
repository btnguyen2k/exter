package user

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/prom"
	_ "github.com/go-sql-driver/mysql"

	"main/src/gvabe/bo"
)

// NewUserDaoSql is helper method to create SQL-implementation of UserDao
func NewUserDaoSql(sqlc *prom.SqlConnect, tableName string, dbFlavor prom.DbFlavor) UserDao {
	dao := &UserDaoSql{}
	dao.UniversalDao = bo.NewUniversalDaoSql(sqlc, tableName, dbFlavor, nil)
	return dao
}

// UserDaoSql is SQL-implementation of AppDao
type UserDaoSql struct {
	bo.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{bo.ColId: gbo.GboGetAttrUnsafe(bo.FieldId, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *UserDaoSql) toBo(gbo godal.IGenericBo) *User {
	return NewUserFromUniversal(dao.ToUniversalBo(gbo))
}

// toGbo transforms business object to godal.IGenericBo
func (dao *UserDaoSql) toGbo(user *User) godal.IGenericBo {
	if user == nil {
		return nil
	}
	return dao.ToGenericBo(user.UniversalBo)
}

// Delete implements UserDao.Delete
func (dao *UserDaoSql) Delete(user *User) (bool, error) {
	return dao.UniversalDao.Delete(user.UniversalBo)
}

// Create implements UserDao.Create
func (dao *UserDaoSql) Create(user *User) (bool, error) {
	return dao.UniversalDao.Create(user.sync().UniversalBo)
}

// Get implements UserDao.Get
func (dao *UserDaoSql) Get(id string) (*User, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewUserFromUniversal(ubo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoSql) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*User, 0)
	for _, ubo := range uboList {
		app := NewUserFromUniversal(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetAll implements UserDao.GetAll
func (dao *UserDaoSql) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoSql) Update(user *User) (bool, error) {
	return dao.UniversalDao.Update(user.sync().UniversalBo)
}

package user

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
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

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter.
func (dao *UserDaoAwsDynamodb) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.FieldId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
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
	if err != nil {
		return nil, err
	}
	return NewUserFromUbo(ubo), nil
}

// getN implements UserDao.getN.
func (dao *UserDaoAwsDynamodb) getN(fromOffset, maxNumRows int) ([]*User, error) {
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

// getAll implements UserDao.getAll.
func (dao *UserDaoAwsDynamodb) getAll() ([]*User, error) {
	return dao.getN(0, 0)
}

// Update implements UserDao.Update.
func (dao *UserDaoAwsDynamodb) Update(bo *User) (bool, error) {
	return dao.UniversalDao.Update(bo.sync().UniversalBo)
}

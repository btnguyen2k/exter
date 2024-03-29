package user

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewUserDaoAwsDynamodb is helper function to create AWS DynamoDB-implementation of UserDao.
func NewUserDaoAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) UserDao {
	var spec *henge.DynamodbDaoSpec = nil
	dao := &UserDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, spec)}
	dao.spec = spec
	return dao
}

// InitUserTableAwsDynamodb is helper function to initialize AWS DynamoDB table(s) to store users.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitUserTableAwsDynamodb(adc *prom.AwsDynamodbConnect, tableName string) error {
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 1, MainTableWcu: 1}
	return henge.InitDynamodbTables(adc, tableName, spec)
}

// UserDaoAwsDynamodb is AWS DynamoDB-implementation of UserDao.
type UserDaoAwsDynamodb struct {
	henge.UniversalDao
	spec *henge.DynamodbDaoSpec
}

// Delete implements UserDao.Delete.
func (dao *UserDaoAwsDynamodb) Delete(bo *User) (bool, error) {
	return dao.UniversalDao.Delete(bo.UniversalBo)
}

// Create implements UserDao.Create.
func (dao *UserDaoAwsDynamodb) Create(bo *User) (bool, error) {
	ubo := bo.sync().UniversalBo
	if dao.spec != nil && dao.spec.PkPrefix != "" {
		ubo.SetExtraAttr(dao.spec.PkPrefix, dao.spec.PkPrefixValue)
	}
	return dao.UniversalDao.Create(ubo)
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

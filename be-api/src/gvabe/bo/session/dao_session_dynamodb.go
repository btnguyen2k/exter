package session

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewSessionDaoAwsDynamodb is helper method to create AWS DynamoDB-implementation of SessionDao.
func NewSessionDaoAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) SessionDao {
	var spec *henge.DynamodbDaoSpec = nil
	dao := &SessionDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, spec)}
	dao.spec = spec
	return dao
}

// InitSessionTableAwsDynamodb is helper function to initialize AWS DynamoDB table(s) to store sessions.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitSessionTableAwsDynamodb(adc *prom.AwsDynamodbConnect, tableName string) error {
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 1, MainTableWcu: 1}
	return henge.InitDynamodbTables(adc, tableName, spec)
}

// SessionDaoAwsDynamodb is AWS DynamoDB-implementation of SessionDao.
type SessionDaoAwsDynamodb struct {
	henge.UniversalDao
	spec *henge.DynamodbDaoSpec
}

// Delete implements SessionDao.Delete.
func (dao *SessionDaoAwsDynamodb) Delete(sess *Session) (bool, error) {
	return dao.UniversalDao.Delete(sess.UniversalBo)
}

// Get implements SessionDao.Get.
func (dao *SessionDaoAwsDynamodb) Get(id string) (*Session, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewSessionFromUbo(ubo), err
}

// Save implements SessionDao.Save.
func (dao *SessionDaoAwsDynamodb) Save(sess *Session) (bool, error) {
	ubo := sess.sync().UniversalBo
	if dao.spec != nil && dao.spec.PkPrefix != "" {
		ubo.SetExtraAttr(dao.spec.PkPrefix, dao.spec.PkPrefixValue)
	}
	ok, _, err := dao.UniversalDao.Save(ubo)
	return ok, err
}

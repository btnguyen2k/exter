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

// Update implements SessionDao.Save.
func (dao *SessionDaoAwsDynamodb) Save(sess *Session) (bool, error) {
	ubo := sess.sync().UniversalBo
	if dao.spec != nil && dao.spec.PkPrefix != "" {
		ubo.SetExtraAttr(dao.spec.PkPrefix, dao.spec.PkPrefixValue)
	}
	ok, _, err := dao.UniversalDao.Save(ubo)
	return ok, err
}

package session

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewSessionDaoAwsDynamodb is helper method to create AWS DynamoDB-implementation of SessionDao.
func NewSessionDaoAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) SessionDao {
	dao := &SessionDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, nil)}
	return dao
}

// SessionDaoAwsDynamodb is AWS DynamoDB-implementation of SessionDao.
type SessionDaoAwsDynamodb struct {
	henge.UniversalDao
}

// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
// func (dao *SessionDaoAwsDynamodb) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
// 	return map[string]interface{}{henge.FieldId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
// }

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
	ok, _, err := dao.UniversalDao.Save(sess.sync().UniversalBo)
	return ok, err
}

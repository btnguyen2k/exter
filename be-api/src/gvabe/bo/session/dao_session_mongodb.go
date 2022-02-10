package session

import (
	"strings"

	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

// NewSessionDaoMongo is helper method to create MongoDB-implementation of SessionDao.
func NewSessionDaoMongo(mc *prom.MongoConnect, collectionName string) SessionDao {
	txMode := strings.Index(strings.ToLower(mc.GetUrl()), "replicaset=") > 0
	dao := &SessionDaoMongo{UniversalDao: henge.NewUniversalDaoMongo(mc, collectionName, txMode)}
	return dao
}

// InitSessionTableMongo is helper function to initialize MongoDB table (collection) to store sessions.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitSessionTableMongo(mc *prom.MongoConnect, collectionName string) error {
	return henge.InitMongoCollection(mc, collectionName)
}

// SessionDaoMongo is MongoDB-implementation of SessionDao.
type SessionDaoMongo struct {
	henge.UniversalDao
}

// Delete implements SessionDao.Delete.
func (dao *SessionDaoMongo) Delete(sess *Session) (bool, error) {
	return dao.UniversalDao.Delete(sess.UniversalBo)
}

// Get implements SessionDao.Get.
func (dao *SessionDaoMongo) Get(id string) (*Session, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewSessionFromUbo(ubo), err
}

// Update implements SessionDao.Save.
func (dao *SessionDaoMongo) Save(sess *Session) (bool, error) {
	ok, _, err := dao.UniversalDao.Save(sess.sync().UniversalBo)
	return ok, err
}

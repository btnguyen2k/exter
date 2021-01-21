package session

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"
)

const (
	SqlCol_Session_SessionType = "ztype"
	SqlCol_Session_IdSource    = "zidsrc"
	SqlCol_Session_AppId       = "zappid"
	SqlCol_Session_UserId      = "zuid"
	SqlCol_Session_Expiry      = "zexpiry"
)

// NewSessionDaoSql is helper method to create SQL-implementation of SessionDao.
func NewSessionDaoSql(sqlc *prom.SqlConnect, tableName string) SessionDao {
	dao := &SessionDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, true, map[string]string{
		SqlCol_Session_SessionType: FieldSession_SessionType,
		SqlCol_Session_IdSource:    FieldSession_IdSource,
		SqlCol_Session_AppId:       FieldSession_AppId,
		SqlCol_Session_UserId:      FieldSession_UserId,
		SqlCol_Session_Expiry:      FieldSession_Expiry,
	})
	return dao
}

// SessionDaoSql is SQL-implementation of SessionDao.
type SessionDaoSql struct {
	henge.UniversalDao
}

// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter.
// func (dao *SessionDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
// 	return map[string]interface{}{henge.SqlColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
// }

// Delete implements SessionDao.Delete.
func (dao *SessionDaoSql) Delete(sess *Session) (bool, error) {
	return dao.UniversalDao.Delete(sess.UniversalBo)
}

// Get implements SessionDao.Get.
func (dao *SessionDaoSql) Get(id string) (*Session, error) {
	ubo, err := dao.UniversalDao.Get(id)
	return NewSessionFromUbo(ubo), err
}

// Update implements SessionDao.Save.
func (dao *SessionDaoSql) Save(sess *Session) (bool, error) {
	ok, _, err := dao.UniversalDao.Save(sess.sync().UniversalBo)
	return ok, err
}

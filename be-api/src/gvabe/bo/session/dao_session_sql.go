package session

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/prom"

	"main/src/henge"
)

const (
	TableSession           = "exter_session"
	ColSession_SessionType = "ztype"
	ColSession_IdSource    = "zidsrc"
	ColSession_AppId       = "zappid"
	ColSession_UserId      = "zuid"
	ColSession_Expiry      = "zexpiry"
)

// NewSessionDaoSql is helper method to create SQL-implementation of SessionDao
func NewSessionDaoSql(sqlc *prom.SqlConnect, tableName string) SessionDao {
	dao := &SessionDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(sqlc, tableName, map[string]string{
		ColSession_SessionType: FieldSession_SessionType,
		ColSession_IdSource:    FieldSession_IdSource,
		ColSession_AppId:       FieldSession_AppId,
		ColSession_UserId:      FieldSession_UserId,
		ColSession_Expiry:      FieldSession_Expiry,
	})
	return dao
}

// SessionDaoSql is SQL-implementation of AppDao
type SessionDaoSql struct {
	henge.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *SessionDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.ColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
}

// Delete implements SessionDao.Delete
func (dao *SessionDaoSql) Delete(sess *Session) (bool, error) {
	return dao.UniversalDao.Delete(sess.UniversalBo.Clone())
}

// // Create implements SessionDao.Create
// func (dao *SessionDaoSql) Create(sess *Session) (bool, error) {
// 	return dao.UniversalDao.Create(sess.sync().UniversalBo.Clone())
// }

// Get implements SessionDao.Get
func (dao *SessionDaoSql) Get(id string) (*Session, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewSessionFromUbo(ubo), nil
}

// // Update implements SessionDao.Update
// func (dao *SessionDaoSql) Update(sess *Session) (bool, error) {
// 	return dao.UniversalDao.Update(sess.sync().UniversalBo.Clone())
// }

// Update implements SessionDao.Save
func (dao *SessionDaoSql) Save(sess *Session) (bool, error) {
	return dao.UniversalDao.Save(sess.sync().UniversalBo.Clone())
}

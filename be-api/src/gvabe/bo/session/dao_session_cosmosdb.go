package session

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

// NewSessionDaoCosmosdb is helper method to create CosmosDB-implementation of SessionDao.
func NewSessionDaoCosmosdb(sqlc *prom.SqlConnect, tableName string) SessionDao {
	spec := &henge.CosmosdbDaoSpec{PkName: bo.CosmosdbPkName, TxModeOnWrite: true}
	innerDao := SessionDaoSql{UniversalDao: henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)}
	dao := &SessionDaoCosmosdb{SessionDaoSql: innerDao, spec: spec}
	return dao
}

// SessionDaoCosmosdb is CosmosDB-implementation of SessionDao.
type SessionDaoCosmosdb struct {
	SessionDaoSql
	spec *henge.CosmosdbDaoSpec
}

// Update implements SessionDao.Save.
func (dao *SessionDaoCosmosdb) Save(sess *Session) (bool, error) {
	ubo := sess.sync().UniversalBo
	if dao.spec != nil && dao.spec.PkName != "" && dao.spec.PkValue != "" {
		ubo.SetExtraAttr(dao.spec.PkName, dao.spec.PkValue)
	}
	ok, _, err := dao.UniversalDao.Save(ubo)
	return ok, err
}

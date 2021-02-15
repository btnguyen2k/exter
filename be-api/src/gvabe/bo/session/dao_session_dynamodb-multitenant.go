package session

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo"
)

const (
	dynamodbPkValueSession = "session"
)

// NewSessionDaoMultitenantAwsDynamodb is helper method to create AWS DynamoDB-implementation (multi-tenant table) of SessionDao.
func NewSessionDaoMultitenantAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) SessionDao {
	spec := &henge.DynamodbDaoSpec{PkPrefix: bo.DynamodbMultitenantPkName, PkPrefixValue: dynamodbPkValueSession}
	dao := &SessionDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, spec)}
	dao.spec = spec
	return dao
}

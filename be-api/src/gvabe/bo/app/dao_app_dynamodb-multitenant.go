package app

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo"
)

const (
	dynamodbPkValueApp = "app"
)

// NewAppDaoMultitenantAwsDynamodb is helper method to create AWS DynamoDB-implementation (multi-tenant table) of AppDao.
func NewAppDaoMultitenantAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) AppDao {
	spec := &henge.DynamodbDaoSpec{PkPrefix: bo.DynamodbMultitenantPkName, PkPrefixValue: dynamodbPkValueApp}
	dao := &AppDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, spec)}
	dao.spec = spec
	return dao
}

package user

import (
	"github.com/btnguyen2k/prom"

	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bo"
)

const (
	dynamodbPkValueUser = "user"
)

// NewUserDaoMultitenantAwsDynamodb is helper method to create AWS DynamoDB-implementation (multi-tenant table) of UserDao.
func NewUserDaoMultitenantAwsDynamodb(dync *prom.AwsDynamodbConnect, tableName string) UserDao {
	spec := &henge.DynamodbDaoSpec{PkPrefix: bo.DynamodbMultitenantPkName, PkPrefixValue: dynamodbPkValueUser}
	dao := &UserDaoAwsDynamodb{UniversalDao: henge.NewUniversalDaoDynamodb(dync, tableName, spec)}
	dao.spec = spec
	return dao
}

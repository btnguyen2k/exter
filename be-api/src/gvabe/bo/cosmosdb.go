package bo

import (
	"github.com/btnguyen2k/henge"
)

const (
	CosmosdbPkName = henge.FieldId

	CosmosdbMultitenantTableName      = "exter_mt"
	CosmosdbMultitenantPkName         = "__mtpk"
	CosmosdbMultitenantPkValueApp     = "app"
	CosmosdbMultitenantPkValueSession = "session"
	CosmosdbMultitenantPkValueUser    = "user"
)

package bo

import (
	"fmt"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

const (
	CosmosdbPkName = henge.FieldId

	CosmosdbMultitenantTableName      = "exter_mt"
	CosmosdbMultitenantPkName         = "__mtpk"
	CosmosdbMultitenantPkValueApp     = "app"
	CosmosdbMultitenantPkValueSession = "session"
	CosmosdbMultitenantPkValueUser    = "user"
)

// InitMultitenantTableCosmosdb is helper function to initialize Cosmos DB multi-tenant table(s) to store BO.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitMultitenantTableCosmosdb(sqlc *prom.SqlConnect, tableName string) error {
	switch sqlc.GetDbFlavor() {
	case prom.FlavorCosmosDb:
		return henge.InitCosmosdbCollection(sqlc, tableName, &henge.CosmosdbCollectionSpec{Pk: CosmosdbMultitenantPkName})
	}
	return fmt.Errorf("unsupported database type %v", sqlc.GetDbFlavor())
}

package bo

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

const (
	DynamodbMultitenantTableName = "exter_mt"
	DynamodbMultitenantPkName    = "__mtpk"
)

// InitMultitenantTableAwsDynamodb is helper function to initialize AWS DynamoDB multi-tenant table(s) to store BO.
// This function also creates table indexes if needed.
//
// Available since v0.7.0.
func InitMultitenantTableAwsDynamodb(adc *prom.AwsDynamodbConnect, tableName string) error {
	spec := &henge.DynamodbTablesSpec{
		MainTableRcu: 1, MainTableWcu: 1,
		CreateUidxTable: true, UidxTableRcu: 1, UidxTableWcu: 1,
		MainTableCustomAttrs: []prom.AwsDynamodbNameAndType{{Name: DynamodbMultitenantPkName, Type: prom.AwsAttrTypeString}},
		MainTablePkPrefix:    DynamodbMultitenantPkName,
	}
	return henge.InitDynamodbTables(adc, tableName, spec)
}

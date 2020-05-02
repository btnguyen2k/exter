package bo

import (
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	_ "github.com/lib/pq"
)

func cloneMap(src map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range src {
		result[k] = v
	}
	return result
}

// NewUniversalDaoSql is helper method to create SQL-implementation of UniversalDao
func NewUniversalDaoSql(sqlc *prom.SqlConnect, tableName string, dbFlavor prom.DbFlavor, mapExtraColNameToField map[string]string) UniversalDao {
	dao := &UniversalDaoSql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	myCols := append([]string{}, cols...)
	myMapFieldToColName := cloneMap(mapFieldToColName)
	myMapColNameToField := cloneMap(mapColNameToField)
	for col, field := range mapExtraColNameToField {
		myCols = append(myCols, col)
		myMapColNameToField[col] = field
		myMapFieldToColName[field] = col
	}
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: myMapFieldToColName},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: myMapColNameToField},
		ColumnsListMap:              map[string][]string{tableName: myCols},
	})
	dao.SetSqlFlavor(dbFlavor)
	return dao
}

const (
	TableApp  = "exter_app"
	TableUser = "exter_user"

	ColId          = "zid"
	ColData        = "zdata"
	ColTimeCreated = "ztcreated"
	ColTimeUpdated = "ztupdated"
	ColAppVersion  = "zaversion"
)

var (
	cols              = []string{ColId, ColData, ColTimeCreated, ColTimeUpdated, ColAppVersion}
	mapFieldToColName = map[string]interface{}{
		FieldId:          ColId,
		FieldData:        ColData,
		FieldTimeCreated: ColTimeCreated,
		FieldTimeUpdated: ColTimeUpdated,
		FieldAppVersion:  ColAppVersion,
	}
	mapColNameToField = map[string]interface{}{
		ColId:          FieldId,
		ColData:        FieldData,
		ColTimeCreated: FieldTimeCreated,
		ColTimeUpdated: FieldTimeUpdated,
		ColAppVersion:  FieldAppVersion,
	}
)

type UniversalDaoSql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UniversalDaoSql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{ColId: bo.GboGetAttrUnsafe(FieldId, reddo.TypeString)}
}

// ToUniversalBo transforms godal.IGenericBo to business object.
func (dao *UniversalDaoSql) ToUniversalBo(gbo godal.IGenericBo) *UniversalBo {
	if gbo == nil {
		return nil
	}
	extraFields := make(map[string]interface{})
	gbo.GboTransferViaJson(&extraFields)
	for _, field := range allFields {
		delete(extraFields, field)
	}
	return &UniversalBo{
		Id:          gbo.GboGetAttrUnsafe(FieldId, reddo.TypeString).(string),
		DataJson:    gbo.GboGetAttrUnsafe(FieldData, reddo.TypeString).(string),
		TimeCreated: gbo.GboGetAttrUnsafe(FieldTimeCreated, reddo.TypeTime).(time.Time),
		TimeUpdated: gbo.GboGetAttrUnsafe(FieldTimeUpdated, reddo.TypeTime).(time.Time),
		AppVersion:  gbo.GboGetAttrUnsafe(FieldAppVersion, reddo.TypeUint).(uint64),
		extraFields: extraFields,
	}
}

// ToGenericBo transforms business object to godal.IGenericBo
func (dao *UniversalDaoSql) ToGenericBo(ubo *UniversalBo) godal.IGenericBo {
	if ubo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(FieldId, ubo.Id)
	gbo.GboSetAttr(FieldData, ubo.DataJson)
	gbo.GboSetAttr(FieldTimeCreated, ubo.TimeCreated)
	gbo.GboSetAttr(FieldTimeUpdated, ubo.TimeUpdated)
	gbo.GboSetAttr(FieldAppVersion, ubo.AppVersion)
	for k, v := range ubo.extraFields {
		gbo.GboSetAttr(k, v)
	}
	return gbo
}

// Delete implements AppDao.Delete
func (dao *UniversalDaoSql) Delete(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}

// Create implements AppDao.Create
func (dao *UniversalDaoSql) Create(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}

// Get implements AppDao.Get
func (dao *UniversalDaoSql) Get(id string) (*UniversalBo, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{ColId: id})
	if err != nil {
		return nil, err
	}
	return dao.ToUniversalBo(gbo), nil
}

// GetN implements AppDao.GetN
func (dao *UniversalDaoSql) GetN(fromOffset, maxNumRows int) ([]*UniversalBo, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(ColId)
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, ordering, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*UniversalBo, 0)
	for _, gbo := range gboList {
		bo := dao.ToUniversalBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements AppDao.GetAll
func (dao *UniversalDaoSql) GetAll() ([]*UniversalBo, error) {
	return dao.GetN(0, 0)
}

// Update implements AppDao.Update
func (dao *UniversalDaoSql) Update(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}

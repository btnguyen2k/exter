package app

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo"
)

// NewAppDaoSql is helper method to create SQL-implementation of AppDao
func NewAppDaoSql(sqlc *prom.SqlConnect, tableName string, dbFlavor prom.DbFlavor) AppDao {
	dao := &AppDaoSql{}
	dao.UniversalDao = bo.NewUniversalDaoSql(sqlc, tableName, dbFlavor)
	return dao
}

// AppDaoSql is SQL-implementation of AppDao
type AppDaoSql struct {
	bo.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *AppDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{bo.ColId: gbo.GboGetAttrUnsafe(bo.FieldId, reddo.TypeString)}
}

// // ‭toBo transforms godal.IGenericBo to business object.
// func (dao *AppDaoSql) toBo(gbo godal.IGenericBo) *App {
// 	return NewAppFromUniversal(dao.ToUniversalBo(gbo))
// }

// // toGbo transforms business object to godal.IGenericBo
// func (dao *AppDaoSql) toGbo(app *App) godal.IGenericBo {
// 	if app == nil {
// 		return nil
// 	}
// 	js, _ := json.Marshal(app)
// 	app.UniversalBo.DataJson = string(js)
// 	return dao.ToGenericBo(app.UniversalBo)
// }

// Delete implements AppDao.Delete
func (dao *AppDaoSql) Delete(app *App) (bool, error) {
	return dao.UniversalDao.Delete(app.UniversalBo)
}

// Create implements AppDao.Create
func (dao *AppDaoSql) Create(app *App) (bool, error) {
	return dao.UniversalDao.Create(app.sync().UniversalBo)
}

// Get implements AppDao.Get
func (dao *AppDaoSql) Get(id string) (*App, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewAppFromUniversal(ubo), nil
}

// GetN implements AppDao.GetN
func (dao *AppDaoSql) GetN(fromOffset, maxNumRows int) ([]*App, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*App, 0)
	for _, ubo := range uboList {
		app := NewAppFromUniversal(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetAll implements AppDao.GetAll
func (dao *AppDaoSql) GetAll() ([]*App, error) {
	return dao.GetN(0, 0)
}

// Update implements AppDao.Update
func (dao *AppDaoSql) Update(app *App) (bool, error) {
	return dao.UniversalDao.Update(app.sync().UniversalBo)
}

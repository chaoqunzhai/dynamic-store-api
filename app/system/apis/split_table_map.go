package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	models2 "go-admin/app/company/models"
	"go-admin/global"
	"time"

	"go-admin/app/system/models"
	"go-admin/app/system/service"
	"go-admin/app/system/service/dto"
	"go-admin/common/actions"
)

type SplitTableMap struct {
	api.Api
}

// GetPage 获取SplitTableMap列表
// @Summary 获取SplitTableMap列表
// @Description 获取SplitTableMap列表
// @Tags SplitTableMap
// @Param layer query string false "排序"
// @Param enable query int64 false "开关"
// @Param cId query int64 false "公司ID"
// @Param type query int64 false "映射表的类型"
// @Param table query string false "对应表的名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SplitTableMap}} "{"code": 200, "data": [...]}"
// @Router /api/v1/split-table-map [get]
// @Security Bearer
func (e SplitTableMap) GetPage(c *gin.Context) {
	req := dto.SplitTableMapGetPageReq{}
	s := service.SplitTableMap{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.SplitTableMap, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SplitTableMap失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SplitTableMap
// @Summary 获取SplitTableMap
// @Description 获取SplitTableMap
// @Tags SplitTableMap
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SplitTableMap} "{"code": 200, "data": [...]}"
// @Router /api/v1/split-table-map/{id} [get]
// @Security Bearer
func (e SplitTableMap) Get(c *gin.Context) {
	req := dto.SplitTableMapGetReq{}
	s := service.SplitTableMap{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.SplitTableMap

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SplitTableMap失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SplitTableMap
// @Summary 创建SplitTableMap
// @Description 创建SplitTableMap
// @Tags SplitTableMap
// @Accept application/json
// @Product application/json
// @Param data body dto.SplitTableMapInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/split-table-map [post]
// @Security Bearer
// 分表的创建,如果是订单类表的类型,因为订单表是有父表和多子表来分表存储订单数据,所以当父表进行了分表,子表也跟随分表
func (e SplitTableMap) Insert(c *gin.Context) {
	req := dto.SplitTableMapInsertReq{}
	s := service.SplitTableMap{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	nowUnix := fmt.Sprintf("%v", time.Now().Unix())[4:]
	tableName := ""
	subTableName :=""
	switch req.Type {
	case global.SplitOrder:
		//订单父表名称
		tableName = fmt.Sprintf("%v_%v_%v", global.SplitOrderDefaultTableName, req.CId, nowUnix)
		subTableName = fmt.Sprintf("%v_%v_%v", global.SplitOrderDefaultSubTableName, req.CId, nowUnix)
	default:
		e.Error(500, nil, "分来类型不存在")
		return
	}
	var count int64
	e.Orm.Model(&models.SplitTableMap{}).Where("c_id = ? and type = ? and enable = ?", req.CId, req.Type, true).Count(&count)
	if count > 0 {

		e.Error(500, errors.New("分表已经存在"), "分表已经存在")
		return
	}

	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	req.Name = tableName
	uid, err := s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建分表失败,%s", err.Error()))
		return
	}
	//数据创建成功后,创建一张分表
	//因为订单表是做了分表的
	switch req.Type {
	case global.SplitOrder:
		//表名称
		SplitRow := models2.Orders{}
		newOrmObject := e.Orm.Model(&SplitRow).Table(tableName)
		createErr := newOrmObject.Migrator().CreateTable(&SplitRow)
		if createErr != nil {
			e.Orm.Unscoped().Delete(&models.SplitTableMap{}, uid)
			e.Error(500, createErr, fmt.Sprintf("创建分表失败,%s", createErr.Error()))
			return
		}
		//订单子表规格表自动创建
		SplitSubRow := models2.OrderSpecs{}
		newSubOrmObject := e.Orm.Model(&SplitSubRow).Table(subTableName)
		createSubErr := newSubOrmObject.Migrator().CreateTable(&SplitSubRow)
		if createSubErr != nil {
			e.Orm.Unscoped().Delete(&models.SplitTableMap{}, uid)
			e.Error(500, createSubErr, fmt.Sprintf("创建子表失败,%s", createSubErr.Error()))
			return
		}
	default:
		e.Error(500, nil, "分来类型不存在")
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SplitTableMap
// @Summary 修改SplitTableMap
// @Description 修改SplitTableMap
// @Tags SplitTableMap
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SplitTableMapUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/split-table-map/{id} [put]
// @Security Bearer
//func (e SplitTableMap) Update(c *gin.Context) {
//	req := dto.SplitTableMapUpdateReq{}
//	s := service.SplitTableMap{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	var count int64
//	e.Orm.Model(&models.SplitTableMap{}).Where("id = ?", req.Id).Count(&count)
//	if count == 0 {
//		e.Error(500, errors.New("数据不存在"), "数据不存在")
//		return
//	}
//	var oldRow models.SplitTableMap
//	e.Orm.Model(&models.SplitTableMap{}).Where("c_id = ? and type = ? and enable = ?", req.CId, req.Type, true).Limit(1).Find(&oldRow)
//
//	if oldRow.Id != 0 {
//		if oldRow.Id != req.Id {
//			e.Error(500, errors.New("分表不可重复"), "分表不可重复")
//			return
//		}
//	}
//	if req.Name == "" {
//		e.Error(500, errors.New("请输入分表名称"), "请输入分表名称")
//		return
//	}
//	var splitRow models.SplitTableMap
//	e.Orm.Model(&models.SplitTableMap{}).Where("name =  ? and enable = ?", req.Name, true).Limit(1).Find(&splitRow)
//	if splitRow.Id > 0 && splitRow.Id != req.Id {
//		e.Error(500, errors.New("分表名称已经存在"), "分表名称已经存在")
//		return
//	}
//	oldTableName, err := s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改分表失败,%s", err.Error()))
//		return
//	}
//
//
//	//查询oldTableName这个表
//	//1.查询原表是否存在,如果存在 重命名表名
//	//2.如果是一个新的表,那就创建吧
//	switch req.Type {
//	case global.SplitOrder:
//		tableName := fmt.Sprintf("%v_%v",global.SplitOrderDefaultTableName,req.Name)
//		subTableName := fmt.Sprintf("%v_%v", global.SplitOrderDefaultSubTableName, req.Name)
//		SplitRow := models2.Orders{}
//		SplitSubRow := models2.OrderSpecs{}
//		valid := e.Orm.Migrator().HasTable(oldTableName)
//		if valid {
//			reNameErr := e.Orm.Migrator().RenameTable(oldTableName, tableName)
//			if reNameErr != nil {
//				e.Error(500, reNameErr, fmt.Sprintf("表名更新失败,请重试,%s", reNameErr.Error()))
//				return
//			}
//		} else {
//			newOrmObject := e.Orm.Model(&SplitRow).Table(tableName)
//			createErr := newOrmObject.Migrator().CreateTable(&SplitRow)
//			if createErr != nil {
//				e.Error(500, createErr, fmt.Sprintf("表名更新失败,请重试,%s", createErr.Error()))
//				return
//			}
//
//			newSubOrmObject := e.Orm.Model(&SplitSubRow).Table(subTableName)
//			createSubErr := newSubOrmObject.Migrator().CreateTable(&SplitRow)
//			if createSubErr != nil {
//				e.Error(500, createSubErr, fmt.Sprintf("子表名更新失败,%s", createSubErr.Error()))
//				return
//			}
//		}
//	}
//
//	e.OK(req.GetId(), "修改成功")
//}

// Delete 删除SplitTableMap
// @Summary 删除SplitTableMap
// @Description 删除SplitTableMap
// @Tags SplitTableMap
// @Param data body dto.SplitTableMapDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/split-table-map [delete]
// @Security Bearer
func (e SplitTableMap) Delete(c *gin.Context) {
	s := service.SplitTableMap{}
	req := dto.SplitTableMapDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SplitTableMap失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

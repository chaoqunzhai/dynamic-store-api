package apis

import (
	"errors"
	"fmt"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type GoodsUnit struct {
	api.Api
}

// GetPage 获取GoodsUnit列表
// @Summary 获取GoodsUnit列表
// @Description 获取GoodsUnit列表
// @Tags GoodsUnit
// @Param cId query string false "大BID"
// @Param name query string false "单位"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsUnit}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-unit [get]
// @Security Bearer
func (e GoodsUnit) GetPage(c *gin.Context) {
    req := dto.GoodsUnitGetPageReq{}
    s := service.GoodsUnit{}
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
	list := make([]models.GoodsUnit, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsUnit失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsUnit
// @Summary 获取GoodsUnit
// @Description 获取GoodsUnit
// @Tags GoodsUnit
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsUnit} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-unit/{id} [get]
// @Security Bearer
func (e GoodsUnit) Get(c *gin.Context) {
	req := dto.GoodsUnitGetReq{}
	s := service.GoodsUnit{}
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
	var object models.GoodsUnit

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsUnit失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建GoodsUnit
// @Summary 创建GoodsUnit
// @Description 创建GoodsUnit
// @Tags GoodsUnit
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsUnitInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-unit [post]
// @Security Bearer
func (e GoodsUnit) Insert(c *gin.Context) {
    req := dto.GoodsUnitInsertReq{}
    s := service.GoodsUnit{}
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var thisCount int64
	e.Orm.Model(&models.GoodsUnit{}).Where("c_id = ?",userDto.CId).Count(&thisCount)
	if thisCount > global.CompanyMaxUnit {
		msg:=fmt.Sprintf("商品单位最大数量上限为:%v",global.CompanyMaxUnit)
		e.Error(500, errors.New(msg), msg)
		return
	}
	var count int64
	var object models.GoodsUnit
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	req.CId = userDto.CId
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建GoodsUnit失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsUnit
// @Summary 修改GoodsUnit
// @Description 修改GoodsUnit
// @Tags GoodsUnit
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsUnitUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-unit/{id} [put]
// @Security Bearer
func (e GoodsUnit) Update(c *gin.Context) {
    req := dto.GoodsUnitUpdateReq{}
    s := service.GoodsUnit{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var count int64
	e.Orm.Model(&models.GoodsUnit{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.GoodsUnit
	e.Orm.Model(&models.GoodsUnit{}).Scopes(actions.PermissionSysUser(oldRow.TableName(),userDto)).Where("name = ? ", req.Name).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}

	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	req.CId = userDto.CId
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改GoodsUnit失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除GoodsUnit
// @Summary 删除GoodsUnit
// @Description 删除GoodsUnit
// @Tags GoodsUnit
// @Param data body dto.GoodsUnitDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-unit [delete]
// @Security Bearer
func (e GoodsUnit) Delete(c *gin.Context) {
    s := service.GoodsUnit{}
    req := dto.GoodsUnitDeleteReq{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	notRemoveIds := make([]int, 0)
	for _, d := range req.Ids {
		var count int64
		e.Orm.Model(&models.GoodsSpecs{}).Where("unit_id = ? and c_id = ?",d,userDto.CId).Count(&count)
		//是否被商品规格绑定
		if count > 0 {
			notRemoveIds = append(notRemoveIds, d)
		}
	}
	if len(notRemoveIds) > 0 {
		e.Error(500, errors.New("存在商品关联不可删除！"), "存在商品关联不可删除！")
		return
	}

	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除GoodsUnit失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

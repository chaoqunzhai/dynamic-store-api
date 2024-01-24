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

type GoodsBrand struct {
	api.Api
}

// GetPage 获取GoodsBrand列表
// @Summary 获取GoodsBrand列表
// @Description 获取GoodsBrand列表
// @Tags GoodsBrand
// @Param layer query string false "排序"
// @Param cId query string false "大BID"
// @Param name query string false "品牌名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsBrand}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-brand [get]
// @Security Bearer
func (e GoodsBrand) GetPage(c *gin.Context) {
    req := dto.GoodsBrandGetPageReq{}
    s := service.GoodsBrand{}
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
	list := make([]models.GoodsBrand, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsBrand失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsBrand
// @Summary 获取GoodsBrand
// @Description 获取GoodsBrand
// @Tags GoodsBrand
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsBrand} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-brand/{id} [get]
// @Security Bearer
func (e GoodsBrand) Get(c *gin.Context) {
	req := dto.GoodsBrandGetReq{}
	s := service.GoodsBrand{}
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
	var object models.GoodsBrand

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsBrand失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建GoodsBrand
// @Summary 创建GoodsBrand
// @Description 创建GoodsBrand
// @Tags GoodsBrand
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsBrandInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-brand [post]
// @Security Bearer
func (e GoodsBrand) Insert(c *gin.Context) {
    req := dto.GoodsBrandInsertReq{}
    s := service.GoodsBrand{}
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
	e.Orm.Model(&models.GoodsBrand{}).Where("c_id = ?",userDto.CId).Count(&thisCount)
	if thisCount > global.CompanyMaxBrand {
		msg:=fmt.Sprintf("商品品牌最大数量上限为:%v",global.CompanyMaxBrand)

		e.Error(500, errors.New(msg), msg)
		return
	}
	var count int64
	var object models.GoodsBrand
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	req.CId = userDto.CId
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建GoodsBrand失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsBrand
// @Summary 修改GoodsBrand
// @Description 修改GoodsBrand
// @Tags GoodsBrand
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsBrandUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-brand/{id} [put]
// @Security Bearer
func (e GoodsBrand) Update(c *gin.Context) {
    req := dto.GoodsBrandUpdateReq{}
    s := service.GoodsBrand{}
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
	e.Orm.Model(&models.GoodsBrand{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.GoodsBrand
	e.Orm.Model(&models.GoodsBrand{}).Scopes(actions.PermissionSysUser(oldRow.TableName(),userDto)).Where("name = ? ", req.Name).Limit(1).Find(&oldRow)

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
		e.Error(500, err, fmt.Sprintf("修改GoodsBrand失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除GoodsBrand
// @Summary 删除GoodsBrand
// @Description 删除GoodsBrand
// @Tags GoodsBrand
// @Param data body dto.GoodsBrandDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-brand [delete]
// @Security Bearer
func (e GoodsBrand) Delete(c *gin.Context) {
    s := service.GoodsBrand{}
    req := dto.GoodsBrandDeleteReq{}
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

	notRemoveIds := make([]int, 0)
	for _, d := range req.Ids {
		var bindCount int64
		whereSql := fmt.Sprintf("SELECT COUNT(*) as count from goods_mark_brand where brand_id = %v", d)
		e.Orm.Raw(whereSql).Scan(&bindCount)
		//是否被商品绑定
		if bindCount > 0 {
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
		e.Error(500, err, fmt.Sprintf("删除GoodsBrand失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

package apis

import (
	"errors"
	"fmt"
	sys "go-admin/app/admin/models"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
)

type Shop struct {
	api.Api
}

// GetPage 获取Shop列表
// @Summary 获取Shop列表
// @Description 获取Shop列表
// @Tags Shop
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param userId query string false "管理员ID"
// @Param name query string false "小B名称"
// @Param phone query string false "联系手机号"
// @Param userName query string false "小B负责人名称"
// @Param lineId query string false "归属配送路线"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Shop}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop [get]
// @Security Bearer
func (e Shop) GetPage(c *gin.Context) {
    req := dto.ShopGetPageReq{}
    s := service.Shop{}
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
	list := make([]models.Shop, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户信息失败,%s", err.Error()))
        return
	}
	result :=make([]interface{},0)
	for _,row:=range list{
		cache :=row
		if row.LineId > 0 {
			var lineRow models2.Line
			e.Orm.Model(&models2.Line{}).Where("id = ? and enable = ?",row.LineId,true).Limit(1).Find(&lineRow)
			if lineRow.Id > 0 {
				cache.Line = lineRow.Name
			}
		}
		result = append(result,cache)

	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Shop
// @Summary 获取Shop
// @Description 获取Shop
// @Tags Shop
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Shop} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop/{id} [get]
// @Security Bearer
func (e Shop) Get(c *gin.Context) {
	req := dto.ShopGetReq{}
	s := service.Shop{}
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
	var object models.Shop

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户信息失败,%s", err.Error()))
        return
	}

	if object.LineId > 0 {
		var lineRow models2.Line
		e.Orm.Model(&models2.Line{}).Where("id = ? and enable = ?",object.LineId,true).Limit(1).Find(&lineRow)
		if lineRow.Id > 0 {
			object.Line = lineRow.Name
		}
	}
	if object.CreateBy  > 0 {
		var userRow sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Where("user_id = ? and enable = ?",object.CreateBy,true).Limit(1).Find(&userRow)
		if userRow.UserId > 0 {
			object.CreateUser = userRow.Username
		}
	}
	//下单记录查询

	//充值记录查询


	e.OK( object, "查询成功")
}

// Insert 创建Shop
// @Summary 创建Shop
// @Description 创建Shop
// @Tags Shop
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop [post]
// @Security Bearer
func (e Shop) Insert(c *gin.Context) {
    req := dto.ShopInsertReq{}
    s := service.Shop{}
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

	var count int64
	e.Orm.Model(&models.Shop{}).Where("c_id = ? and name = ?", userDto.CId, req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	err = s.Insert(userDto.CId,&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("用户创建失败,%s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Shop
// @Summary 修改Shop
// @Description 修改Shop
// @Tags Shop
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop/{id} [put]
// @Security Bearer
func (e Shop) Update(c *gin.Context) {
    req := dto.ShopUpdateReq{}
    s := service.Shop{}
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
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var count int64
	e.Orm.Model(&models.Shop{}).Where("id = ?",req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.Shop
	e.Orm.Model(&models.Shop{}).Where("name = ? and c_id = ?",req.Name,userDto.CId).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改失败,%s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除Shop
// @Summary 删除Shop
// @Description 删除Shop
// @Tags Shop
// @Param data body dto.ShopDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop [delete]
// @Security Bearer
func (e Shop) Delete(c *gin.Context) {
    s := service.Shop{}
    req := dto.ShopDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除Shop失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

func (e Shop)Amount(c *gin.Context)  {
	req := dto.ShopAmountReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
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
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Where("id = ?",req.ShopId).First(&object).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("客户不存在"), "客户不存在")
		return
	}
	Scene:=""
	switch req.Action {
	case global.UserNumberAdd:
		object.Amount += req.Number
		Scene = fmt.Sprintf("手动增加%v元",req.Number)
	case global.UserNumberReduce:
		if req.Number > object.Amount {
			e.Error(500, errors.New("数值非法"), "数值非法")
			return
		}
		object.Amount -=req.Number
		Scene = fmt.Sprintf("手动减少%v元",req.Number)
	case global.UserNumberSet:
		object.Amount = req.Number
		Scene = fmt.Sprintf("手动设置为%v元",req.Number)
	}
	object.UpdateBy = user.GetUserId(c)
	e.Orm.Save(&object)
	row:=models.ShopBalanceLog{
		CId: userDto.CId,
		ShopId: req.ShopId,
		Desc: req.Desc,
		Money: req.Number,
		Scene:Scene,
		Action: req.Action,
	}
	row.CreateBy = user.GetUserId(c)
	e.Orm.Create(&row)
	e.OK("","successful")
	return
}
func (e Shop)Integral(c *gin.Context)  {
	req := dto.ShopIntegralReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
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
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Where("id = ?",req.ShopId).First(&object).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("客户不存在"), "客户不存在")
		return
	}
	Scene:=""
	switch req.Action {
	case global.UserNumberAdd:
		object.Integral += req.Number
		Scene = fmt.Sprintf("手动增加%v个积分",req.Number)
	case global.UserNumberReduce:
		if req.Number > object.Integral {
			e.Error(500, errors.New("数值非法"), "数值非法")
			return
		}
		object.Integral -=req.Number
		Scene = fmt.Sprintf("手动减少%v个积分",req.Number)
	case global.UserNumberSet:
		object.Integral = req.Number
		Scene = fmt.Sprintf("手动积分设置为%v",req.Number)

	}
	object.UpdateBy = user.GetUserId(c)
	e.Orm.Save(&object)
	row:=models.ShopIntegralLog{
		CId: userDto.CId,
		ShopId: req.ShopId,
		Desc: req.Desc,
		Number: req.Number,
		Scene:Scene,
		Action: req.Action,
	}
	row.CreateBy = user.GetUserId(c)
	e.Orm.Create(&row)
	e.OK("","successful")
	return

}
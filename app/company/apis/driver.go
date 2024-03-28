package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	sys "go-admin/app/admin/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Driver struct {
	api.Api
}
func (e Driver) MiniApi(c *gin.Context) {
	err := e.MakeContext(c).
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
	datalist:=make([]models.Driver,0)
	var object models.Driver
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Select("id,name").Order(global.OrderLayerKey).Find(&datalist)

	result:=make([]map[string]interface{},0)
	for _,row:=range datalist{
		result = append(result, map[string]interface{}{
			"id":row.Id,
			"name":row.Name,
		})
	}
	e.OK(result,"操作成功")
	return
}
// GetPage 获取Driver列表
// @Summary 获取Driver列表
// @Description 获取Driver列表
// @Tags Driver
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param userId query string false "关联的用户ID"
// @Param name query string false "司机名称"
// @Param phone query string false "联系手机号"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Driver}} "{"code": 200, "data": [...]}"
// @Router /api/v1/driver [get]
// @Security Bearer
func (e Driver) GetPage(c *gin.Context) {
	req := dto.DriverGetPageReq{}

	s := service.Driver{}
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
	list := make([]models.Driver, 0)

	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Driver失败，\r\n失败信息 %s", err.Error()))
		return
	}
	result := make([]interface{}, 0)
	for _, row := range list {

		var bindLine models.Line
		e.Orm.Model(&models.Line{}).Select("name,id").Where("driver_id = ? and enable = ?", row.Id,true).Limit(1).Find(&bindLine)
		if bindLine.Id > 0 {
			//如果开启了过滤 已经绑定的司机 那就不返回
			if req.Exclude {
				row.Disable = true
			}
			row.LineName = bindLine.Name

		}
		if row.UserId > 0 {
			var userObject sys.SysUser
			e.Orm.Model(&sys.SysUser{}).Select("password").Where("user_id = ? and c_id = ?",row.UserId,row.CId).Limit(1).Find(&userObject)
			row.Password = userObject.Password
		}
		result = append(result, row)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Driver
// @Summary 获取Driver
// @Description 获取Driver
// @Tags Driver
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Driver} "{"code": 200, "data": [...]}"
// @Router /api/v1/driver/{id} [get]
// @Security Bearer
func (e Driver) Get(c *gin.Context) {
	req := dto.DriverGetReq{}
	s := service.Driver{}
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
	var object models.Driver

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Driver失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建Driver
// @Summary 创建Driver
// @Description 创建Driver
// @Tags Driver
// @Accept application/json
// @Product application/json
// @Param data body dto.DriverInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/driver [post]
// @Security Bearer
func (e Driver) Insert(c *gin.Context) {
	req := dto.DriverInsertReq{}
	s := service.Driver{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	var count int64
	var object models.Driver
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Driver
// @Summary 修改Driver
// @Description 修改Driver
// @Tags Driver
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.DriverUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/driver/{id} [put]
// @Security Bearer
func (e Driver) Update(c *gin.Context) {
	req := dto.DriverUpdateReq{}
	s := service.Driver{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
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
	e.Orm.Model(&models.Driver{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.Driver
	e.Orm.Model(&models.Driver{}).Scopes(actions.PermissionSysUser(oldRow.TableName(),userDto)).Where("name = ? ", req.Name).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}


	err = s.Update(userDto.CId, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "修改成功")
}

// Delete 删除Driver
// @Summary 删除Driver
// @Description 删除Driver
// @Tags Driver
// @Param data body dto.DriverDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/driver [delete]
// @Security Bearer
func (e Driver) Delete(c *gin.Context) {
	req := dto.DriverDeleteReq{}

	s := service.Driver{}
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
	p := actions.GetPermissionFromContext(c)

	notRemoveIds := make([]int, 0)
	for _, d := range req.Ids {
		var count int64
		e.Orm.Model(&models.Line{}).Where("driver_id = ? and c_id = ?", d,userDto.CId).Count(&count)
		//如果有路线关联了司机，那就是不可删除
		if count > 0 {
			notRemoveIds = append(notRemoveIds, d)
		}
	}
	if len(notRemoveIds) > 0 {
		e.Error(500, errors.New("存在关联路线不可删除！"), "存在关联路线不可删除！")
		return
	}
	err = s.Remove(&req, p,userDto.CId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("司机信息删除失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

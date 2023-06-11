package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
)

type CompanyRole struct {
	api.Api
}

// GetPage 获取CompanyRole列表
// @Summary 获取CompanyRole列表
// @Description 获取CompanyRole列表
// @Tags CompanyRole
// @Param name query string false ""
// @Param enable query string false ""
// @Param sort query string false ""
// @Param remark query string false ""
// @Param admin query string false ""
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyRole}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-role [get]
// @Security Bearer
func (e CompanyRole) GetPage(c *gin.Context) {
	req := dto.CompanyRoleGetPageReq{}
	s := service.CompanyRole{}
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
	list := make([]models.CompanyRole, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyRole失败，\r\n失败信息 %s", err.Error()))
		return
	}
	result := make([]interface{}, 0)
	for _, row := range list {
		menuIds := make([]int, 0)
		for _, bindMenu := range row.SysMenu {
			menuIds = append(menuIds, bindMenu.Id)
		}
		r := map[string]interface{}{
			"name":       row.Name,
			"id":         row.Id,
			"layer":      row.Layer,
			"desc":       row.Desc,
			"created_at": row.CreatedAt,
			"user_count": len(row.SysUser),
			"menuIds":    menuIds,
			"enable":     row.Enable,
		}
		result = append(result, r)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CompanyRole
// @Summary 获取CompanyRole
// @Description 获取CompanyRole
// @Tags CompanyRole
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyRole} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-role/{id} [get]
// @Security Bearer
func (e CompanyRole) Get(c *gin.Context) {
	req := dto.CompanyRoleGetReq{}
	s := service.CompanyRole{}
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
	var object models.CompanyRole

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyRole失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建CompanyRole
// @Summary 创建CompanyRole
// @Description 创建CompanyRole
// @Tags CompanyRole
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyRoleInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-role [post]
// @Security Bearer
func (e CompanyRole) Insert(c *gin.Context) {
	req := dto.CompanyRoleInsertReq{}
	s := service.CompanyRole{}
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var countAll int64
	e.Orm.Model(&models.CompanyRole{}).Where("c_id = ?", userDto.CId).Count(&countAll)

	CompanyCnf := business.GetCompanyCnf(userDto.CId, "role", e.Orm)
	MaxRole := CompanyCnf["role"]

	if countAll > int64(MaxRole) {
		e.Error(500, errors.New(fmt.Sprintf("角色最多只能创建%v个", MaxRole)), fmt.Sprintf("角色最多只能创建%v个", MaxRole))
		return
	}
	var count int64
	e.Orm.Model(&models.CompanyRole{}).Where("c_id = ? and name = ?", userDto.CId, req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建角色失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyRole
// @Summary 修改CompanyRole
// @Description 修改CompanyRole
// @Tags CompanyRole
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyRoleUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-role/{id} [put]
// @Security Bearer
func (e CompanyRole) Update(c *gin.Context) {
	req := dto.CompanyRoleUpdateReq{}
	s := service.CompanyRole{}
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
	e.Orm.Model(&models.CompanyRole{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.CompanyRole
	e.Orm.Model(&models.CompanyRole{}).Where("name = ? and c_id = ?", req.Name, userDto.CId).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyRole失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除CompanyRole
// @Summary 删除CompanyRole
// @Description 删除CompanyRole
// @Tags CompanyRole
// @Param data body dto.CompanyRoleDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-role [delete]
// @Security Bearer
func (e CompanyRole) Delete(c *gin.Context) {
	s := service.CompanyRole{}
	req := dto.CompanyRoleDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除角色失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

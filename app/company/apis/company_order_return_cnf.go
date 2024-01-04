package apis

import (
    "fmt"
	"github.com/gin-gonic/gin/binding"
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

type CompanyOrderReturnCnf struct {
	api.Api
}

func (e CompanyOrderReturnCnf) Enable(c *gin.Context) {
	req := dto.UpdateEnableReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).

		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	e.Orm.Model(&models.CompanyOrderReturnCnf{}).Where("id = ?",req.Id).Updates(map[string]interface{}{
		"enable":req.Enable,
	})
	e.OK("", "更新成功")
}
// GetPage 获取CompanyOrderReturnCnf列表
// @Summary 获取CompanyOrderReturnCnf列表
// @Description 获取CompanyOrderReturnCnf列表
// @Tags CompanyOrderReturnCnf
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param value query string false "配送文案"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyOrderReturnCnf}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-order-return-cnf [get]
// @Security Bearer
func (e CompanyOrderReturnCnf) GetPage(c *gin.Context) {
    req := dto.CompanyOrderReturnCnfGetPageReq{}
    s := service.CompanyOrderReturnCnf{}
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
	list := make([]models.CompanyOrderReturnCnf, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyOrderReturnCnf失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CompanyOrderReturnCnf
// @Summary 获取CompanyOrderReturnCnf
// @Description 获取CompanyOrderReturnCnf
// @Tags CompanyOrderReturnCnf
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyOrderReturnCnf} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-order-return-cnf/{id} [get]
// @Security Bearer
func (e CompanyOrderReturnCnf) Get(c *gin.Context) {
	req := dto.CompanyOrderReturnCnfGetReq{}
	s := service.CompanyOrderReturnCnf{}
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
	var object models.CompanyOrderReturnCnf

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyOrderReturnCnf失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CompanyOrderReturnCnf
// @Summary 创建CompanyOrderReturnCnf
// @Description 创建CompanyOrderReturnCnf
// @Tags CompanyOrderReturnCnf
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyOrderReturnCnfInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-order-return-cnf [post]
// @Security Bearer
func (e CompanyOrderReturnCnf) Insert(c *gin.Context) {
    req := dto.CompanyOrderReturnCnfInsertReq{}
    s := service.CompanyOrderReturnCnf{}
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
	e.Orm.Model(&models.CompanyOrderReturnCnf{}).Where("c_id = ?",userDto.CId).Count(&count)
	if count >= global.MaxCompanyOrderReturnCnf {
		e.Error(500, nil,fmt.Sprintf("最多只能创建 %v个配置",global.MaxCompanyOrderReturnCnf))
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))


	req.CId = userDto.CId
	req.Enable = true
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建CompanyOrderReturnCnf失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyOrderReturnCnf
// @Summary 修改CompanyOrderReturnCnf
// @Description 修改CompanyOrderReturnCnf
// @Tags CompanyOrderReturnCnf
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyOrderReturnCnfUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-order-return-cnf/{id} [put]
// @Security Bearer
func (e CompanyOrderReturnCnf) Update(c *gin.Context) {
    req := dto.CompanyOrderReturnCnfUpdateReq{}
    s := service.CompanyOrderReturnCnf{}
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

	req.CId = userDto.CId
	req.Enable = true
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyOrderReturnCnf失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CompanyOrderReturnCnf
// @Summary 删除CompanyOrderReturnCnf
// @Description 删除CompanyOrderReturnCnf
// @Tags CompanyOrderReturnCnf
// @Param data body dto.CompanyOrderReturnCnfDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-order-return-cnf [delete]
// @Security Bearer
func (e CompanyOrderReturnCnf) Delete(c *gin.Context) {
    s := service.CompanyOrderReturnCnf{}
    req := dto.CompanyOrderReturnCnfDeleteReq{}
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

	err = s.Remove(userDto.CId,&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CompanyOrderReturnCnf失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

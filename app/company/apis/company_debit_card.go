package apis

import (
    "fmt"
	customUser "go-admin/common/jwt/user"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type CompanyDebitCard struct {
	api.Api
}

// GetPage 获取CompanyDebitCard列表
// @Summary 获取CompanyDebitCard列表
// @Description 获取CompanyDebitCard列表
// @Tags CompanyDebitCard
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "持卡人名称"
// @Param backName query string false "开户行"
// @Param cardNumber query string false "银行卡号"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyDebitCard}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-debit-card [get]
// @Security Bearer
func (e CompanyDebitCard) GetPage(c *gin.Context) {
    req := dto.CompanyDebitCardGetPageReq{}
    s := service.CompanyDebitCard{}
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
	list := make([]models.CompanyDebitCard, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyDebitCard失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CompanyDebitCard
// @Summary 获取CompanyDebitCard
// @Description 获取CompanyDebitCard
// @Tags CompanyDebitCard
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyDebitCard} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-debit-card/{id} [get]
// @Security Bearer
func (e CompanyDebitCard) Get(c *gin.Context) {
	req := dto.CompanyDebitCardGetReq{}
	s := service.CompanyDebitCard{}
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
	var object models.CompanyDebitCard

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyDebitCard失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CompanyDebitCard
// @Summary 创建CompanyDebitCard
// @Description 创建CompanyDebitCard
// @Tags CompanyDebitCard
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyDebitCardInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-debit-card [post]
// @Security Bearer
func (e CompanyDebitCard) Insert(c *gin.Context) {
    req := dto.CompanyDebitCardInsertReq{}
    s := service.CompanyDebitCard{}
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

    req.CId = userDto.CId
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建CompanyDebitCard失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyDebitCard
// @Summary 修改CompanyDebitCard
// @Description 修改CompanyDebitCard
// @Tags CompanyDebitCard
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyDebitCardUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-debit-card/{id} [put]
// @Security Bearer
func (e CompanyDebitCard) Update(c *gin.Context) {
    req := dto.CompanyDebitCardUpdateReq{}
    s := service.CompanyDebitCard{}
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
	req.CId = userDto.CId
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyDebitCard失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CompanyDebitCard
// @Summary 删除CompanyDebitCard
// @Description 删除CompanyDebitCard
// @Tags CompanyDebitCard
// @Param data body dto.CompanyDebitCardDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-debit-card [delete]
// @Security Bearer
func (e CompanyDebitCard) Delete(c *gin.Context) {
    s := service.CompanyDebitCard{}
    req := dto.CompanyDebitCardDeleteReq{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.CId = userDto.CId
	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CompanyDebitCard失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

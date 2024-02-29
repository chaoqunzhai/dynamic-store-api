package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type GradeVip struct {
	api.Api
}

// GetPage 获取GradeVip列表
// @Summary 获取GradeVip列表
// @Description 获取GradeVip列表
// @Tags GradeVip
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "等级名称"
// @Param weight query string false "权重,从小到大"
// @Param upgrade query string false "升级条件,满多少金额,自动升级Weight+1"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GradeVip}} "{"code": 200, "data": [...]}"
// @Router /api/v1/grade-vip [get]
// @Security Bearer
func (e GradeVip) GetPage(c *gin.Context) {
	req := dto.GradeVipGetPageReq{}
	s := service.GradeVip{}
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
	list := make([]models.GradeVip, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GradeVip失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GradeVip
// @Summary 获取GradeVip
// @Description 获取GradeVip
// @Tags GradeVip
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GradeVip} "{"code": 200, "data": [...]}"
// @Router /api/v1/grade-vip/{id} [get]
// @Security Bearer
func (e GradeVip) Get(c *gin.Context) {
	req := dto.GradeVipGetReq{}
	s := service.GradeVip{}
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
	var object models.GradeVip

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GradeVip失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建GradeVip
// @Summary 创建GradeVip
// @Description 创建GradeVip
// @Tags GradeVip
// @Accept application/json
// @Product application/json
// @Param data body dto.GradeVipInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/grade-vip [post]
// @Security Bearer
func (e GradeVip) Insert(c *gin.Context) {
	req := dto.GradeVipInsertReq{}
	s := service.GradeVip{}
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
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "vip", e.Orm)
	MaxNumber := CompanyCnf["vip"]

	if countAll >= int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("VIP最多只可创建%v个", MaxNumber)),
			fmt.Sprintf("VIP最多只可创建%v个", MaxNumber))
		return
	}
	var count int64
	var object models.GradeVip
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where(" name = ?", req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建VIP失败信息,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GradeVip
// @Summary 修改GradeVip
// @Description 修改GradeVip
// @Tags GradeVip
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GradeVipUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/grade-vip/{id} [put]
// @Security Bearer
func (e GradeVip) Update(c *gin.Context) {
	req := dto.GradeVipUpdateReq{}
	s := service.GradeVip{}
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
	e.Orm.Model(&models.GradeVip{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.GradeVip
	e.Orm.Model(&models.GradeVip{}).Scopes(actions.PermissionSysUser(oldRow.TableName(),userDto)).Where("name = ? ", req.Name).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改信息失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除GradeVip
// @Summary 删除GradeVip
// @Description 删除GradeVip
// @Tags GradeVip
// @Param data body dto.GradeVipDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/grade-vip [delete]
// @Security Bearer
func (e GradeVip) Delete(c *gin.Context) {
	s := service.GradeVip{}
	req := dto.GradeVipDeleteReq{}
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

	err = s.Remove(&req,userDto.CId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除GradeVip失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

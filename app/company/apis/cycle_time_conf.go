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

type CycleTimeConf struct {
	api.Api
}

// GetPage 获取CycleTimeConf列表
// @Summary 获取CycleTimeConf列表
// @Description 获取CycleTimeConf列表
// @Tags CycleTimeConf
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param type query string false "类型,每天,每周"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CycleTimeConf}} "{"code": 200, "data": [...]}"
// @Router /api/v1/cycle-time-conf [get]
// @Security Bearer
func (e CycleTimeConf) GetPage(c *gin.Context) {
	req := dto.CycleTimeConfGetPageReq{}
	s := service.CycleTimeConf{}
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
	list := make([]models.CycleTimeConf, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CycleTimeConf
// @Summary 获取CycleTimeConf
// @Description 获取CycleTimeConf
// @Tags CycleTimeConf
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CycleTimeConf} "{"code": 200, "data": [...]}"
// @Router /api/v1/cycle-time-conf/{id} [get]
// @Security Bearer
func (e CycleTimeConf) Get(c *gin.Context) {
	req := dto.CycleTimeConfGetReq{}
	s := service.CycleTimeConf{}
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
	var object models.CycleTimeConf

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建CycleTimeConf
// @Summary 创建CycleTimeConf
// @Description 创建CycleTimeConf
// @Tags CycleTimeConf
// @Accept application/json
// @Product application/json
// @Param data body dto.CycleTimeConfInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/cycle-time-conf [post]
// @Security Bearer
func (e CycleTimeConf) Insert(c *gin.Context) {
	req := dto.CycleTimeConfInsertReq{}
	s := service.CycleTimeConf{}
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

	//时间不能重复,也就是start_time 和end_time 是不能存在相同的
	whereSql :=""
	switch req.Type {

	case global.CyCleTimeDay:
		whereSql =fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v'",
			userDto.CId,true,global.CyCleTimeDay,req.StartTime,req.EndTime)

	case global.CyCleTimeWeek:
		whereSql =fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v' and start_week = %v and end_week = %v",
			userDto.CId,true,global.CyCleTimeWeek,req.StartTime,req.EndTime,req.StartWeek,req.EndWeek)


	default:
		e.Error(500, nil, "非法类型")
		return
	}
	e.Orm.Model(&models.CycleTimeConf{}).Where(whereSql).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("时间不可重复"), "时间不可重复")
		return
	}
	err = s.Insert(userDto.CId,&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("时间区间创建失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CycleTimeConf
// @Summary 修改CycleTimeConf
// @Description 修改CycleTimeConf
// @Tags CycleTimeConf
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CycleTimeConfUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/cycle-time-conf/{id} [put]
// @Security Bearer
func (e CycleTimeConf) Update(c *gin.Context) {
	req := dto.CycleTimeConfUpdateReq{}
	s := service.CycleTimeConf{}
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

	whereSql :=""
	var timeConf models.CycleTimeConf
	switch req.Type {

	case global.CyCleTimeDay:

		whereSql =fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v'",
			userDto.CId,true,global.CyCleTimeDay,req.StartTime,req.EndTime)

	case global.CyCleTimeWeek:
		whereSql =fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v' and start_week = %v and end_week = %v",
			userDto.CId,true,global.CyCleTimeWeek,req.StartTime,req.EndTime,req.StartWeek,req.EndWeek)
	default:
		e.Error(500, nil, "非法类型")
		return
	}
	e.Orm.Model(&models.CycleTimeConf{}).Where(whereSql).Limit(1).Find(&timeConf)
	if timeConf.Id >0 && timeConf.Id != req.Id {
		e.Error(500, nil, "时间不可重复")
		return
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除CycleTimeConf
// @Summary 删除CycleTimeConf
// @Description 删除CycleTimeConf
// @Tags CycleTimeConf
// @Param data body dto.CycleTimeConfDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/cycle-time-conf [delete]
// @Security Bearer
func (e CycleTimeConf) Delete(c *gin.Context) {
	s := service.CycleTimeConf{}
	req := dto.CycleTimeConfDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
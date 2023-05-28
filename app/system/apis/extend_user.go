package apis

import (
	"errors"
	"fmt"
	sys "go-admin/app/admin/models"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/system/models"
	"go-admin/app/system/service"
	"go-admin/app/system/service/dto"
	"go-admin/common/actions"
)

type ExtendUser struct {
	api.Api
}

// GetPage 获取ExtendUser列表
// @Summary 获取ExtendUser列表
// @Description 获取ExtendUser列表
// @Tags ExtendUser
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param platform query string false "注册来源"
// @Param gradeId query string false "会员等级"
// @Param suggestId query string false "推荐人ID"
// @Param invitationCode query string false "本人邀请码"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ExtendUser}} "{"code": 200, "data": [...]}"
// @Router /api/v1/extend-user [get]
// @Security Bearer
func (e ExtendUser) GetPage(c *gin.Context) {
    req := dto.ExtendUserGetPageReq{}
    s := service.ExtendUser{}
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
	list := make([]models.ExtendUser, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ExtendUser失败，\r\n失败信息 %s", err.Error()))
        return
	}
	result :=make([]map[string]interface{},0)
	for _,row:=range list{

		u:=map[string]interface{}{
			"id":row.UserId,
			"platform":row.Platform,
			"create_at":row.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if row.UserId == 0 {continue}

		var userRow sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Where("user_id = ? and enable = ?",row.UserId,true).Limit(1).Find(&userRow)
		u["user_name"] = userRow.Username

		if row.GradeId > 0 {
			var gradeRow models2.GradeVip
			e.Orm.Model(&models2.GradeVip{}).Where("id = ? and enable = ?",row.GradeId,true).Limit(1).Find(&gradeRow)
			if gradeRow.Id > 0 {
				u["grade"] = gradeRow.Name
			}
		}

		var shopRow models2.Shop
		e.Orm.Model(&models2.Shop{}).Where("user_id = ?",row.UserId).Limit(1).Find(&shopRow)
		if shopRow.Id > 0 {
			u["amount"] = shopRow.Amount
			u["integral"] = shopRow.Integral
			u["shop"] = shopRow.Name
		}
		result = append(result,u)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ExtendUser
// @Summary 获取ExtendUser
// @Description 获取ExtendUser
// @Tags ExtendUser
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ExtendUser} "{"code": 200, "data": [...]}"
// @Router /api/v1/extend-user/{id} [get]
// @Security Bearer
func (e ExtendUser) Get(c *gin.Context) {
	req := dto.ExtendUserGetReq{}
	s := service.ExtendUser{}
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
	var object models.ExtendUser

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ExtendUser失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ExtendUser
// @Summary 创建ExtendUser
// @Description 创建ExtendUser
// @Tags ExtendUser
// @Accept application/json
// @Product application/json
// @Param data body dto.ExtendUserInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/extend-user [post]
// @Security Bearer
func (e ExtendUser) Insert(c *gin.Context) {
    req := dto.ExtendUserInsertReq{}
    s := service.ExtendUser{}
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
	e.Orm.Model(&models.ExtendUser{}).Where("user_id = ? ",req.UserId).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("已经存在"), "已经存在")
		return
	}
	var userCount int64
	e.Orm.Model(&sys.SysUser{}).Where("user_id = ? ",req.UserId).Count(&userCount)
	if userCount == 0 {
		e.Error(500, errors.New("用户不存在"), "用户不存在")
		return
	}
	//谁创建的那谁就是推荐人
	req.SuggestId = req.CreateBy
	err = s.Insert(userDto.CId,&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建ExtendUser失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ExtendUser
// @Summary 修改ExtendUser
// @Description 修改ExtendUser
// @Tags ExtendUser
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ExtendUserUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/extend-user/{id} [put]
// @Security Bearer
func (e ExtendUser) Update(c *gin.Context) {
    req := dto.ExtendUserUpdateReq{}
    s := service.ExtendUser{}
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

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改ExtendUser失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

func (e ExtendUser) Grade(c *gin.Context)()  {
	s := service.ExtendUser{}
	req := dto.ExtendUserGradeReq{}
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

	e.Orm.Model(&models.ExtendUser{}).Where("id in ?",req.Ids).Updates(map[string]interface{}{
		"grade_id":req.Grade,
		"update_by":user.GetUserId(c),
	})

	e.OK("","successful")
	return

}
// Delete 删除ExtendUser
// @Summary 删除ExtendUser
// @Description 删除ExtendUser
// @Tags ExtendUser
// @Param data body dto.ExtendUserDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/extend-user [delete]
// @Security Bearer
func (e ExtendUser) Delete(c *gin.Context) {
    s := service.ExtendUser{}
    req := dto.ExtendUserDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ExtendUser失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

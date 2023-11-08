package apis

import (

	"go-admin/app/admin/models"
	"go-admin/global"
	"net/http"


	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type SysUser struct {
	api.Api
}


func (e SysUser) GetInfo(c *gin.Context) {
	req := dto.SysUserById{}
	s := service.SysUser{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	var roles = make([]string, 1)
	roles[0] = user.GetRoleName(c)
	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"
	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	var mp = make(map[string]interface{})
	mp["roles"] = []string{
		"admin",
	}
	mp["buttons"] = buttons
	mp["permissions"] = permissions

	sysUser := models.SysUser{}
	req.Id = user.GetUserId(c)
	err = s.Get(&req, p, &sysUser)
	if err != nil {
		e.Error(http.StatusUnauthorized, err, "登录失败")
		return
	}
	mp["introduction"] = " am a super administrator"
	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if sysUser.Avatar != "" {
		mp["avatar"] = sysUser.Avatar
	}
	mp["userName"] = sysUser.Username
	mp["userId"] = sysUser.UserId
	mp["deptId"] = sysUser.DeptId
	mp["name"] = sysUser.Username
	mp["code"] = 200
	e.OK(mp, "")
}
func (e SysUser) GetUserInfo(c *gin.Context) {
	req := dto.SysUserById{}
	s := service.SysUser{}
	r := service.SysRole{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&r.Service).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)

	var mp = make(map[string]interface{})
	sysUser := models.SysUser{}
	userID := user.GetUserId(c)
	req.Id = userID
	err = s.Get(&req, p, &sysUser)
	if err != nil {
		e.Error(http.StatusUnauthorized, err, "登录失败")
		return
	}
	userInfo := map[string]interface{}{
		"store_user_id": sysUser.UserId,
		"user_name":     sysUser.Username,
		"real_name":     sysUser.Username,
		"is_delete":     0,
		"sort":          0,
		"store_id":      0,
		"create_time":   sysUser.CreatedAt.Format("2006-01-02 15:04:05"),
		"update_time":   sysUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	userInfo["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if sysUser.Avatar != "" {
		userInfo["avatar"] = sysUser.Avatar
	}
	rolesMap := map[string]interface{}{
		"permissionList": make([]string, 0),
	}
	super := false
	//超管是获取所有的菜单的


	switch sysUser.RoleId {
	case global.RoleSuper:
		//超管
		super = true
		rolesMap["permissionList"] = make([]string, 0)
	case global.RoleCompany:
		//大B,这里的菜单最好还是跟超管的区分的
		//大B,默认一定是没有权限的
		super = true //其实呢 因为是大B和超管是完全独立的2个服务,所以这里大B就是超管
		rolesMap["permissionList"] = make([]string, 0)
		//rolesMap["permissionList"] = r.GetCustomAdmin(user.GetUserId(c))
	case global.RoleCompanyUser:
		rolesMap["permissionList"] = r.GetCustomById(user.GetUserId(c))
	default:

		e.Error(http.StatusUnauthorized, err, "您没有权限")
		return

	}
	userInfo["isSuper"] = super
	rolesMap["isSuper"] = super
	mp["userInfo"] = userInfo
	mp["roles"] = rolesMap
	e.OK(mp, "")
}

package service

import (
	"fmt"
	"go-admin/common/utils"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	models2 "go-admin/cmd/migrate/migration/models"
)

type SysRole struct {
	service.Service
}






type RoleRow struct {
	RoleId string
}
type RoleMenu struct {
	MenuId string
}

func (e *SysRole) GetCustomAdmin(userId int) []string{

	var list []models2.DyNamicMenu
	e.Orm.Model(&models2.DyNamicMenu{}).Select("path,role").Where("enable = ?",true).Find(&list)
	perms :=make([]string,0)
	for _,row:=range list{
		role:=strings.Split(row.Role,",")
		if utils.IsArray("company",role){
			perms = append(perms,row.Path)
		}
	}
	return perms
}
//获取这个用户关联的自定义菜单权限列表
func (e *SysRole) GetCustomById(userId int) []string{
	permissions := make([]string, 0)
	roleResult:=make([]RoleRow,0)
	//获取用户关联了哪些角色
	whereSQL := fmt.Sprintf("select role_id from company_role_user where user_id = %v",userId)
	e.Orm.Raw(whereSQL).Scan(&roleResult)

	if len(roleResult) == 0 {
		return permissions
	}

	//获取角色关联了哪些菜单
	roleStr:=make([]string,0)
	for _,r:=range roleResult{
		roleStr =append(roleStr,r.RoleId)
	}

	menuResult:=make([]RoleMenu,0)
	whereMenuSql:=fmt.Sprintf("select menu_id from company_role_menu where role_id in (%v)",strings.Join(roleStr,","))
	e.Orm.Raw(whereMenuSql).Scan(&menuResult)

	if len(menuResult) == 0 {
		return permissions
	}
	//根据绑定的菜单ID获取菜单信息

	menuList :=make([]string,0)
	for _,m:=range menuResult{
		menuList =append(menuList,m.MenuId)
	}
	dyNamingMenu :=make([]models2.DyNamicMenu,0)
	e.Orm.Model(&models2.DyNamicMenu{}).Select("path").Where("id in ?",menuList).Find(&dyNamingMenu)
	for _,t:=range dyNamingMenu {
		permissions = append(permissions,t.Path)
	}
	return permissions
}

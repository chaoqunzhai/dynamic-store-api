package service

import (
	"errors"
	"fmt"
	"go-admin/common/utils"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm/clause"

	"github.com/casbin/casbin/v2"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	models3 "go-admin/app/company/models"
	models2 "go-admin/cmd/migrate/migration/models"
	cDto "go-admin/common/dto"
)

type SysRole struct {
	service.Service
}

// GetPage 获取SysRole列表
func (e *SysRole) GetPage(c *dto.SysRoleGetPageReq, list *[]models.SysRole, count *int64) error {
	var err error
	var data models.SysRole

	err = e.Orm.Model(&data).Preload("SysMenu").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取SysRole对象
func (e *SysRole) Get(d *dto.SysRoleGetReq, model *models.SysRole) error {
	var err error
	db := e.Orm.First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	if err != nil {
		e.Log.Errorf("get menuIds error, %s", err.Error())
		return err
	}
	return nil
}

// Insert 创建SysRole对象
func (e *SysRole) Insert(c *dto.SysRoleInsertReq, cb *casbin.SyncedEnforcer) error {
	var err error
	var data models.SysRole
	var dataMenu []models.SysMenu
	err = e.Orm.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&dataMenu).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	c.SysMenu = dataMenu
	c.Generate(&data)
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx := e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var count int64
	err = tx.Model(&data).Where("role_key = ?", c.RoleKey).Count(&count).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if count > 0 {
		err = errors.New("roleKey已存在，需更换在提交！")
		e.Log.Errorf("db error:%s", err)
		return err
	}

	err = tx.Create(&data).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0)
	for _, menu := range dataMenu {
		for _, api := range menu.SysApi {
			if mp[data.RoleKey+"-"+api.Path+"-"+api.Action] != "" {
				mp[data.RoleKey+"-"+api.Path+"-"+api.Action] = ""
				polices = append(polices, []string{data.RoleKey, api.Path, api.Action})
			}
		}
	}

	if len(polices) <= 0 {
		return nil
	}

	// 写入 sys_casbin_rule 权限表里 当前角色数据的记录
	_, err = cb.AddNamedPolicies("p", polices)
	if err != nil {
		return err
	}

	return nil
}

// Update 修改SysRole对象
func (e *SysRole) Update(c *dto.SysRoleUpdateReq, cb *casbin.SyncedEnforcer) error {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx := e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var model = models.SysRole{}
	var mlist = make([]models.SysMenu, 0)
	tx.Preload("SysMenu").First(&model, c.GetId())
	tx.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&mlist)
	err = tx.Model(&model).Association("SysMenu").Delete(model.SysMenu)
	if err != nil {
		e.Log.Errorf("delete policy error:%s", err)
		return err
	}
	c.Generate(&model)
	model.SysMenu = &mlist
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)

	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, err = cb.RemoveFilteredPolicy(0, model.RoleKey)
	if err != nil {
		e.Log.Errorf("delete policy error:%s", err)
		return err
	}
	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0)
	for _, menu := range mlist {
		for _, api := range menu.SysApi {
			if mp[model.RoleKey+"-"+api.Path+"-"+api.Action] != "" {
				mp[model.RoleKey+"-"+api.Path+"-"+api.Action] = ""
				//_, err = cb.AddNamedPolicy("p", model.RoleKey, api.Path, api.Action)
				polices = append(polices, []string{model.RoleKey, api.Path, api.Action})
			}
		}
	}
	if len(polices) <= 0 {
		return nil
	}

	// 写入 sys_casbin_rule 权限表里 当前角色数据的记录
	_, err = cb.AddNamedPolicies("p", polices)
	if err != nil {
		return err
	}
	return nil
}

// Remove 删除SysRole
func (e *SysRole) Remove(c *dto.SysRoleDeleteReq, cb *casbin.SyncedEnforcer) error {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx := e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var model = models.SysRole{}
	tx.Preload("SysMenu").Preload("SysDept").First(&model, c.GetId())
	//删除 SysRole 时，同时删除角色所有 关联其它表 记录 (SysMenu 和 SysMenu)
	db := tx.Select(clause.Associations).Delete(&model)

	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, _ = cb.RemoveFilteredPolicy(0, model.RoleKey)

	return nil
}

// GetRoleMenuId 获取角色对应的菜单ids
func (e *SysRole) GetRoleMenuId(roleId int) ([]int, error) {
	menuIds := make([]int, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		menuIds = append(menuIds, l[i].MenuId)
	}
	return menuIds, nil
}

func (e *SysRole) UpdateDataScope(c *dto.RoleDataScopeReq) *SysRole {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx := e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var dlist = make([]models.SysDept, 0)
	var model = models.SysRole{}
	tx.Preload("SysDept").First(&model, c.RoleId)
	tx.Where("dept_id in ?", c.DeptIds).Find(&dlist)
	// 删除SysRole 和 SysDept 的关联关系
	err = tx.Model(&model).Association("SysDept").Delete(model.SysDept)
	if err != nil {
		e.Log.Errorf("delete SysDept error:%s", err)
		_ = e.AddError(err)
		return e
	}
	c.Generate(&model)
	model.SysDept = dlist
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Model(&model).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if db.RowsAffected == 0 {
		_ = e.AddError(errors.New("无权更新该数据"))
		return e
	}
	return e
}

// UpdateStatus 修改SysRole对象status
func (e *SysRole) UpdateStatus(c *dto.UpdateStatusReq) error {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx := e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var model = models.SysRole{}
	tx.First(&model, c.GetId())
	c.Generate(&model)
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// GetWithName 获取SysRole对象
func (e *SysRole) GetWithName(d *dto.SysRoleByName, model *models.SysRole) *SysRole {
	var err error
	db := e.Orm.Where("role_name = ?", d.RoleName).First(model)
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	if err != nil {
		e.Log.Errorf("get menuIds error, %s", err.Error())
		_ = e.AddError(err)
		return e
	}
	return e
}

// GetById 获取SysRole对象
func (e *SysRole) GetById(roleId int) ([]string, error) {
	permissions := make([]string, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		if l[i].Permission != "" {
			permissions = append(permissions, l[i].Permission)
		}
	}
	return permissions, nil
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
//获取大B开启的菜单
func (e *SysRole) RoleCompany(userId int,CompanyModel models3.Company) []string{

	permissions := make([]string, 0)
	dyNamingMenu :=make([]models2.DyNamicMenu,0)
	e.Orm.Model(&models2.DyNamicMenu{}).Select("path").Where("enable = ?",true).Find(&dyNamingMenu)
	for _,t:=range dyNamingMenu {

		if !CompanyModel.InventoryModule {//没有开启仓库功能,那就跳出

			if strings.HasPrefix(t.Path,"/inventory/"){
				continue
			}
		}
		if !CompanyModel.SaleUserModule { //没有开启业务员功能,那就跳出
			if strings.HasPrefix(t.Path,"/user/salesman"){
				continue
			}
		}
		permissions = append(permissions,t.Path)
	}
	return permissions

}
//获取这个用户关联的自定义菜单权限列表
func (e *SysRole) GetCustomById(userId int,CompanyModel models3.Company) []string{
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
	e.Orm.Model(&models2.DyNamicMenu{}).Select("path").Where("id in ? and enable = ?",menuList,true).Find(&dyNamingMenu)
	for _,t:=range dyNamingMenu {
		if !CompanyModel.InventoryModule {//没有开启仓库功能,那就跳出

			if strings.HasPrefix(t.Path,"/inventory/"){
				continue
			}
		}
		if !CompanyModel.SaleUserModule { //没有开启业务员功能,那就跳出
			if strings.HasPrefix(t.Path,"/user/salesman"){
				continue
			}
		}
		permissions = append(permissions,t.Path)
	}
	return permissions
}

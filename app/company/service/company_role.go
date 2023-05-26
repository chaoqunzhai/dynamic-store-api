package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	sys "go-admin/app/admin/models"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CompanyRole struct {
	service.Service
}

// GetPage 获取CompanyRole列表
func (e *CompanyRole) GetPage(c *dto.CompanyRoleGetPageReq, p *actions.DataPermission, list *[]models.CompanyRole, count *int64) error {
	var err error
	var data models.CompanyRole

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyRoleService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyRole对象
func (e *CompanyRole) Get(d *dto.CompanyRoleGetReq, p *actions.DataPermission, model *models.CompanyRole) error {
	var data models.CompanyRole

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyRole error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
func (e *CompanyRole) getMenuModels(ids []int) (menuList []models.DyNamicMenu) {
	for _, id := range ids {
		var menu models.DyNamicMenu
		e.Orm.Model(&models.DyNamicMenu{}).Where("id = ?", id).First(&menu)
		if menu.Id == 0 {
			continue
		}
		menuList = append(menuList, menu)
	}
	return menuList
}
func (e *CompanyRole) getUserModels(ids []int) (list []sys.SysUser) {
	for _, id := range ids {
		var menu sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Where("user_id = ?", id).First(&menu)
		if menu.UserId == 0 {
			continue
		}
		list = append(list, menu)
	}
	return list
}


// Insert 创建CompanyRole对象
func (e *CompanyRole) Insert(cId int, c *dto.CompanyRoleInsertReq) error {
	var err error
	var data models.CompanyRole
	c.Generate(&data)
	data.CId = cId
	if len(c.Menu) > 0 {
		data.SysMenu = e.getMenuModels(c.Menu)
	}
	if len(c.User) > 0 {
		data.SysUser = e.getUserModels(c.User)
	}
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("角色创建失败", err)
		return err
	}
	return nil
}

// Update 修改CompanyRole对象
func (e *CompanyRole) Update(c *dto.CompanyRoleUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CompanyRole{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CompanyRoleService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CompanyRole
func (e *CompanyRole) Remove(d *dto.CompanyRoleDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyRole

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCompanyRole error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
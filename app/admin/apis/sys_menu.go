package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	JwtUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"sort"
	"strings"
)

type SysMenu struct {
	api.Api
}

// GetPage Menu列表数据
// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/menu [get]
// @Security Bearer
func (e SysMenu) GetPage(c *gin.Context) {
	s := service.SysMenu{}
	req := dto.SysMenuGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var list = make([]models.SysMenu, 0)
	err = s.GetPage(&req, &list).Error
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

// Get 获取菜单详情
// @Summary Menu详情数据
// @Description 获取JSON
// @Tags 菜单
// @Param id path string false "id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/menu/{id} [get]
// @Security Bearer
func (e SysMenu) Get(c *gin.Context) {
	req := dto.SysMenuGetReq{}
	s := new(service.SysMenu)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object = models.SysMenu{}

	err = s.Get(&req, &object).Error
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(object, "查询成功")
}

// Insert 创建菜单
// @Summary 创建菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysMenuInsertReq true "data"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/menu [post]
// @Security Bearer
func (e SysMenu) Insert(c *gin.Context) {
	req := dto.SysMenuInsertReq{}
	s := new(service.SysMenu)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	err = s.Insert(&req).Error
	if err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.GetId(), "创建成功")
}

// Update 修改菜单
// @Summary 修改菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysMenuUpdateReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/menu/{id} [put]
// @Security Bearer
func (e SysMenu) Update(c *gin.Context) {
	req := dto.SysMenuUpdateReq{}
	s := new(service.SysMenu)
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
	err = s.Update(&req).Error
	if err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// Delete 删除菜单
// @Summary 删除菜单
// @Description 删除数据
// @Tags 菜单
// @Param data body dto.SysMenuDeleteReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/menu [delete]
// @Security Bearer
func (e SysMenu) Delete(c *gin.Context) {
	control := new(dto.SysMenuDeleteReq)
	s := new(service.SysMenu)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	err = s.Remove(control).Error
	if err != nil {
		e.Logger.Errorf("RemoveSysMenu error, %s", err)
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(control.GetId(), "删除成功")
}

type MenuRow struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	ParentId  int       `json:"parent_id"`
	Id        int       `json:"id"`
	MetaTitle string    `json:"title"`
	MetaIcon  string    `json:"icon"`
	Hidden    bool      `json:"hidden"`
	KeepAlive bool      `json:"keep_alive"`
	Layer int
	Children  []MenuRow `json:"children"`
}



func getParentAll(parent int, rr MenuRow, data map[int]MenuRow) map[int]MenuRow {

	newMap := make(map[int]MenuRow)
	for _, row := range data {
		if len(row.Children) == 0 {
			newMap[row.Id] = row
			continue
		}

		netList := row.Children

		for k, c := range row.Children {
			if c.Id == parent {
				c.Children = append(c.Children, rr)

			}
			//重新赋值
			row.Children[k] = c
		}
		row.Children = netList
		newMap[row.Id] = row

	}
	return newMap
}
func (e SysMenu) GetAdminMenuRole(c *gin.Context) {
	s := new(service.SysMenu)
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	result, err := s.SetMenuRole(user.GetRoleName(c))

	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.OK(result, "")
}

//parentId 父ID
//当前
func GetChildrenAll(row MenuRow,allMenuList *[]MenuRow) MenuRow  {

	list := *allMenuList

	thisRowList:= make([]MenuRow, 0)
	for _,j:=range list{

		if j.ParentId != row.ParentId{
			continue
		}
		//ms := GetChildrenAll(row, allMenuList)
		thisRowList = append(thisRowList, j)

	}
	row.Children = thisRowList

	return row
}


func (e SysMenu) MbmList(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	userDto, err := JwtUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var CompanyModel models2.Company
	e.Orm.Model(&models2.Company{}).Where("enable = 1 and id = ? ",userDto.CId).First(&CompanyModel)

	var menu = make([]models2.DyMbmMenu, 0)
	e.Orm.Model(&models2.DyMbmMenu{}).Where("enable = ?",true).Order("layer desc").Find(&menu)
	cacheMap := make(map[int]MenuRow, 0)

	result := make([]MenuRow, 0)
	allMenuList:=make([]MenuRow,0)
	for _, row := range menu {

		r := MenuRow{
			Id:        row.Id,
			Name:      row.Name,
			ParentId:  row.ParentId,
			MetaTitle: row.MetaTitle,
			Children:  make([]MenuRow, 0),
			Layer: row.Layer,
		}
		allMenuList = append(allMenuList,r)
		if row.ParentId == 0 {
			//一级菜单 那就直接放进去
			cacheMap[row.Id] = r
		}
	}
	for _, row := range allMenuList {
		if row.ParentId == 0 {
			continue
		}
		//只查子元素
		data,ok:=cacheMap[row.ParentId]

		if ok {
			data.Children = append(data.Children,row)

			cacheMap[row.ParentId] = data
		}else {
			//下面的层级
			cacheMap = getParentAll(row.ParentId, row, cacheMap)
		}
	}

	layerList :=make([]int,0)
	for _,row :=range cacheMap{
		layerList = append(layerList,row.Layer)
	}
	layerList = utils.RemoveRepeatInt(layerList)
	sort.Slice(layerList, func(i, j int) bool {
		return layerList[i] < layerList[j]
	})
	for _,k:=range layerList{
		for _, row := range cacheMap {
			if row.Layer == k{
				result = append(result, row)
			}
		}
	}

	e.OK(result, "")
}



func (e SysMenu) GetMenuRole(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	userDto, err := JwtUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var CompanyModel models2.Company
	e.Orm.Model(&models2.Company{}).Where("enable = 1 and id = ? ",userDto.CId).First(&CompanyModel)

	var menu = make([]models2.DyNamicMenu, 0)
	e.Orm.Model(&models2.DyNamicMenu{}).Where("enable = ?",true).Order("layer asc").Find(&menu)
	cacheMap := make(map[int]MenuRow, 0)

	result := make([]MenuRow, 0)
	allMenuList:=make([]MenuRow,0)
	for _, row := range menu {
		if !CompanyModel.InventoryModule {//没有开启仓库功能,那就跳出

			if strings.HasPrefix(row.Path,"/inventory/"){
				continue
			}
		}
		if !CompanyModel.SaleUserModule { //没有开启业务员功能,那就跳出
			if strings.HasPrefix(row.Path,"/user/salesman"){
				continue
			}
		}
		r := MenuRow{
			Id:        row.Id,
			Name:      row.Name,
			ParentId:  row.ParentId,
			Path:      row.Path,
			Hidden:    row.Hidden,
			KeepAlive: row.KeepAlive,
			MetaIcon:  row.MetaIcon,
			MetaTitle: row.MetaTitle,
			Component: row.Component,
			Children:  make([]MenuRow, 0),
			Layer: row.Layer,
		}
		allMenuList = append(allMenuList,r)
		if row.ParentId == 0 {
			//一级菜单 那就直接放进去
			cacheMap[row.Id] = r
		}
	}
	for _, row := range allMenuList {
		if row.ParentId == 0 {
			continue
		}
		//只查子元素
		data,ok:=cacheMap[row.ParentId]

		if ok {
			data.Children = append(data.Children,row)

			cacheMap[row.ParentId] = data
		}else {
			//下面的层级
			cacheMap = getParentAll(row.ParentId, row, cacheMap)
		}
	}

	layerList :=make([]int,0)
	for _,row :=range cacheMap{
		layerList = append(layerList,row.Layer)
	}
	layerList = utils.RemoveRepeatInt(layerList)
	sort.Slice(layerList, func(i, j int) bool {
		return layerList[i] < layerList[j]
	})
	for _,k:=range layerList{
		for _, row := range cacheMap {
			if row.Layer == k{
				result = append(result, row)
			}
		}
	}

	e.OK(result, "")
}

//// GetMenuIDS 获取角色对应的菜单id数组
//// @Summary 获取角色对应的菜单id数组，设置角色权限使用
//// @Description 获取JSON
//// @Tags 菜单
//// @Param id path int true "id"
//// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
//// @Router /api/v1/menuids/{id} [get]
//// @Security Bearer
//func (e SysMenu) GetMenuIDS(c *gin.Context) {
//	s := new(service.SysMenu)
//	r := service.SysRole{}
//	m := dto.SysRoleByName{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&m, binding.JSON).
//		MakeService(&s.Service).
//		MakeService(&r.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	var data models.SysRole
//	err = r.GetWithName(&m, &data).Error
//
//	//data.RoleName = c.GetString("role")
//	//data.UpdateBy = user.GetUserId(c)
//	//result, err := data.GetIDS(s.Orm)
//
//	if err != nil {
//		e.Logger.Errorf("GetIDS error, %s", err.Error())
//		e.Error(500, err, "获取失败")
//		return
//	}
//	e.OK(result, "")
//}

// GetMenuTreeSelect 根据角色ID查询菜单下拉树结构
// @Summary 角色修改使用的菜单列表
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/json
// @Product application/json
// @Param roleId path int true "roleId"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/menuTreeselect/{roleId} [get]
// @Security Bearer
func (e SysMenu) GetMenuTreeSelect(c *gin.Context) {
	m := service.SysMenu{}
	r := service.SysRole{}
	req := dto.SelectRole{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&m.Service).
		MakeService(&r.Service).
		Bind(&req, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	result, err := m.SetLabel()
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	menuIds := make([]int, 0)
	if req.RoleId != 0 {
		menuIds, err = r.GetRoleMenuId(req.RoleId)
		if err != nil {
			e.Error(500, err, "")
			return
		}
	}
	e.OK(gin.H{
		"menus":       result,
		"checkedKeys": menuIds,
	}, "获取成功")
}

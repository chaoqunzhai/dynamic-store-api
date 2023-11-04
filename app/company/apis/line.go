package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Line struct {
	api.Api
}


func (e Line) UnusedOneLine(c *gin.Context) {

	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var object models.Line
	e.Orm.Model(&models.Line{}).Where("enable = ? and driver_id = 0",true).Scopes(
		actions.Permission(object.TableName(), p),
	).Limit(1).Find(&object)

	if object.Id == 0 {
		e.Error(500, errors.New("路线不存在"), "路线不存在")
		return
	}
	result :=map[string]interface{}{
		"id":object.Id,
		"name":object.Name,
	}

	e.OK(result, "successful")
	return
}
func (e Line) BindShop(c *gin.Context) {
	req := dto.BindLineUserReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var object models.Line
	e.Orm.Model(&models.Line{}).Where("id = ? and enable = ?", req.LineId, true).Scopes(
		actions.Permission(object.TableName(), p),
	).Limit(1).Find(&object)

	if object.Id == 0 {
		e.Error(500, errors.New("路线不存在"), "路线不存在")
		return
	}
	var shopObject models2.Shop
	e.Orm.Model(&shopObject).Where("id in ?", req.ShopId).Scopes(
		actions.Permission(shopObject.TableName(), p),
	).Updates(map[string]interface{}{
		"line_id":    req.LineId,
		"updated_at": time.Now(),
		"update_by":  user.GetUserId(c),
	})

	e.OK("", "successful")
	return
}

func (e Line) UpdateLineBindShopList(c *gin.Context) {
	req := dto.UpdateLineBindShopReq{}
	s := service.Line{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var shopCount int64
	var object models2.Shop
	e.Orm.Model(&models2.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where("id = ?", req.Id).Count(&shopCount)

	if shopCount == 0 {
		e.Error(500, nil, "客户不存在")
		return
	}
	e.Orm.Model(&models2.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where("id = ? ", req.Id).Updates(map[string]interface{}{
		"layer":     req.Layer,
		"enable":    req.Enable,
		"desc":      req.Desc,
		"address":   req.Address,
		"longitude": req.Longitude,
		"latitude":  req.Latitude,
	})
	e.OK("successful", "successful")
	return
}
func (e Line) LineBindShopList(c *gin.Context) {
	req := dto.LineBindShopGetPageReq{}
	s := service.Line{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//路线下的商户,那就是查询商户数据增加 大B + 路线ID
	lineId := c.Param("id")
	lineIdNumber, _ := strconv.Atoi(lineId)
	var lineObject models.Line
	e.Orm.Model(&lineObject).Scopes(actions.PermissionSysUser(lineObject.TableName(),userDto)).Select("id,name").Where("id = ? and enable = ?", lineIdNumber, true).Limit(1).Find(&lineObject)
	if lineObject.Id == 0 {
		e.Error(500, nil, "路线不存在")
		return
	}
	result := make([]interface{}, 0)
	var list []models2.Shop
	var count int64
	p := actions.GetPermissionFromContext(c)
	var shopObject models2.Shop
	e.Orm.Model(&shopObject).Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(shopObject.Name, p),
		).Where("enable = ?  and line_id = ?", true, lineIdNumber).Order(global.OrderLayerKey).Find(&list).Limit(-1).Offset(-1).
		Count(&count)
	for _, row := range list {
		cc := map[string]interface{}{
			"name":       row.Name,
			"phone":      row.Phone,
			"address":    row.Address,
			"local":      fmt.Sprintf("%v,%v", row.Longitude, row.Latitude),
			"line_name":  lineObject.Name,
			"id":         row.Id,
			"layer":      row.Layer,
			"created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"desc":       row.Desc,
		}
		result = append(result, cc)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return

}
func (e Line) MiniApi(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
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
	datalist:=make([]models.Line,0)
	var object models.Line
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Select("id,name").Where("enable = ?",true).Order(global.OrderLayerKey).Find(&datalist)

	result:=make([]map[string]interface{},0)
	for _,row:=range datalist{
		result = append(result, map[string]interface{}{
			"id":row.Id,
			"name":row.Name,
		})
	}
	e.OK(result,"successful")
	return
}

func (e Line) GetPage(c *gin.Context) {
	req := dto.LineGetPageReq{}
	s := service.Line{}
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
	list := make([]models.Line, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Line失败，\r\n失败信息 %s", err.Error()))
		return
	}
	result := make([]interface{}, 0)
	for _, row := range list {
		if row.DriverId > 0 {
			var driverObject models.Driver
			e.Orm.Model(&models.Driver{}).Select("name,phone,id").Where("id = ? and enable = ?", row.DriverId, true).Limit(1).Find(&driverObject)

			if driverObject.Id > 0 {
				row.DriverName = fmt.Sprintf("%v-%v", driverObject.Name, driverObject.Phone)
			}
		}
		var shopCount int64
		e.Orm.Model(&models2.Shop{}).Where("line_id = ? and enable = ?", row.Id, true).Count(&shopCount)
		row.ShopCount = shopCount
		if !row.ExpirationTime.Time.IsZero() {


			row.ExpirationDay = int(row.ExpirationTime.Sub(time.Now()).Hours() / 24)
			row.ExpirationTimeStr = row.ExpirationTime.Format("2006-01-02 15:04:05")
		}else {
			row.ExpirationTimeStr = "无期限"
		}
		result = append(result, row)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Line
// @Summary 获取Line
// @Description 获取Line
// @Tags Line
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Line} "{"code": 200, "data": [...]}"
// @Router /api/v1/line/{id} [get]
// @Security Bearer
func (e Line) Get(c *gin.Context) {
	req := dto.LineGetReq{}
	s := service.Line{}
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
	var object models.Line

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Line失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建Line
// @Summary 创建Line
// @Description 创建Line
// @Tags Line
// @Accept application/json
// @Product application/json
// @Param data body dto.LineInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/line [post]
// @Security Bearer
func (e Line) Insert(c *gin.Context) {
	req := dto.LineInsertReq{}
	s := service.Line{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models.Line
	var countAll int64
	e.Orm.Model(&models.Line{}).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&countAll)

	CompanyCnf := business.GetCompanyCnf(userDto.CId, "line", e.Orm)
	MaxNumber := CompanyCnf["line"]

	if countAll >= int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("线路最多只可创建%v个", MaxNumber)), fmt.Sprintf("线路最多只可创建%v个", MaxNumber))
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	var count int64
	e.Orm.Model(&models.Line{}).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Where("name = ?", userDto.CId, req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}

	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("路线创建失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Line
// @Summary 修改Line
// @Description 修改Line
// @Tags Line
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.LineUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/line/{id} [put]
// @Security Bearer
func (e Line) Update(c *gin.Context) {
	req := dto.LineUpdateReq{}
	s := service.Line{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)


	var count int64
	e.Orm.Model(&models.Line{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.Line
	e.Orm.Model(&models.Line{}).Scopes(actions.PermissionSysUser(oldRow.TableName(),userDto)).Where("name = ? ", req.Name).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	//如果选择了司机,判断司机是否已经被其他路线关联

	if req.DriverId > 0 {

		var validLine models.Line
		e.Orm.Model(&models.Line{}).Scopes(actions.PermissionSysUser(validLine.TableName(),userDto)).Where("driver_id = ? ", req.DriverId).Limit(1).Find(&validLine)

		if validLine.Id != 0 {
			if validLine.Id != req.Id {
				msg := fmt.Sprintf("司机已被,[%v]路线关联", validLine.Name)
				e.Error(500, errors.New(msg), msg)
				return
			}
		}
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改Line失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}



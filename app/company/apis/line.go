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
	shopModels "go-admin/app/shop/models"
	"go-admin/common/actions"
)

type Line struct {
	api.Api
}

func (e Line) QuotaCnf(c *gin.Context)   {
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
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "line", e.Orm)

	MaxNumber := CompanyCnf["line"]
	var count int64
	e.Orm.Model(&models.Line{}).Where("c_id = ?",userDto.CId).Count(&count)

	res:=make(map[string]interface{},0)
	if int(count) < MaxNumber {
		//还有可以创建的路线

		res["show"] = true
		res["count"] =  MaxNumber - int(count)
	}else {
		res["show"] = false
	}

	e.OK(res,"successful")
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
		actions.Permission("line", p),
	).Limit(1).Find(&object)

	if object.Id == 0 {
		e.Error(500, errors.New("路线不存在"), "路线不存在")
		return
	}

	e.Orm.Model(&models2.Shop{}).Where("id in ?", req.ShopId).Scopes(
		actions.Permission("line", p),
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
	e.Orm.Model(&models2.Shop{}).Where("id = ? and c_id = ?", req.Id, userDto.CId).Count(&shopCount)

	if shopCount == 0 {
		e.Error(500, nil, "客户不存在")
		return
	}
	e.Orm.Model(&models2.Shop{}).Where("id = ? and c_id = ?", req.Id, userDto.CId).Updates(map[string]interface{}{
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
	e.Orm.Model(&lineObject).Select("id,name").Where("id = ? and enable = ? and c_id = ?", lineIdNumber, true, userDto.CId).Limit(1).Find(&lineObject)
	if lineObject.Id == 0 {
		e.Error(500, nil, "路线不存在")
		return
	}
	result := make([]interface{}, 0)
	var list []models2.Shop
	var count int64
	p := actions.GetPermissionFromContext(c)
	e.Orm.Model(&models2.Shop{}).Where("enable = ? and c_id = ? and line_id = ?", true, userDto.CId, lineIdNumber).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission("shop", p),
		).Order(global.OrderLayerKey).Find(&list).Limit(-1).Offset(-1).
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
	e.Orm.Model(&models.Line{}).Select("id,name").Where("c_id = ? and enable = ?",userDto.CId,true).Order(global.OrderLayerKey).Find(&datalist)

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
// GetPage 获取Line列表
// @Summary 获取Line列表
// @Description 获取Line列表
// @Tags Line
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "路线名称"
// @Param driverId query string false "关联司机"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Line}} "{"code": 200, "data": [...]}"
// @Router /api/v1/line [get]
// @Security Bearer
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
		var driverObject models.Driver
		e.Orm.Model(&models.Driver{}).Select("name,phone,id").Where("id = ? and enable = ?", row.DriverId, true).Limit(1).Find(&driverObject)

		if driverObject.Id > 0 {
			row.DriverName = fmt.Sprintf("%v-%v", driverObject.Name, driverObject.Phone)
		}
		var shopCount int64
		e.Orm.Model(&models2.Shop{}).Where("line_id = ? and enable = ?", row.Id, true).Count(&shopCount)
		row.ShopCount = shopCount
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
	var countAll int64
	e.Orm.Model(&models.Line{}).Where("c_id = ?", userDto.CId).Count(&countAll)

	CompanyCnf := business.GetCompanyCnf(userDto.CId, "good_class", e.Orm)
	MaxNumber := CompanyCnf["line"]

	if countAll >= int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("线路最多只可创建%v个", MaxNumber)), fmt.Sprintf("线路最多只可创建%v个", MaxNumber))
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	var count int64
	e.Orm.Model(&models.Line{}).Where("c_id = ? and name = ?", userDto.CId, req.Name).Count(&count)
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
	e.Orm.Model(&models.Line{}).Where("name = ? and c_id = ?", req.Name, userDto.CId).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	//如果选择了司机,判断司机是否已经被其他路线关联

	if req.DriverId > 0 {

		var validLine models.Line
		e.Orm.Model(&models.Line{}).Where("driver_id = ? and c_id = ?", req.DriverId, userDto.CId).Limit(1).Find(&validLine)

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

// Delete 删除Line
// @Summary 删除Line
// @Description 删除Line
// @Tags Line
// @Param data body dto.LineDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/line [delete]
// @Security Bearer
func (e Line) Delete(c *gin.Context) {
	s := service.Line{}
	req := dto.LineDeleteReq{}
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
	p := actions.GetPermissionFromContext(c)
	newIds := make([]int, 0)
	for _, line := range req.Ids {
		var count int64
		e.Orm.Model(&shopModels.Shop{}).Where("line_id = ? and c_id = ?", line, userDto.CId).Count(&count)
		if count == 0 {
			newIds = append(newIds, line)
		}
	}
	if len(newIds) == 0 {
		e.Error(500, errors.New("存在关联不可删除！"), "存在关联不可删除！")
		return
	}
	req.Ids = newIds
	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("路线删除失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

package apis

import (
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
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

type Orders struct {
	api.Api
}

func (e Orders) getTableName(cid int) string {
	//先在split分表中查询

	var splitRow models2.SplitTableMap
	e.Orm.Model(&models2.SplitTableMap{}).Where("c_id = ? and enable = ? and type = ?", cid, true, global.SplitOrder).Limit(1).Find(&splitRow)

	tableName := ""
	if splitRow.Id > 0 {
		tableName = splitRow.Name
	} else {
		tableName = global.SplitOrderDefaultTableName
	}
	return tableName
}

// GetPage 获取Orders列表
// @Summary 获取Orders列表
// @Description 获取Orders列表
// @Tags Orders
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param shopId query string false "关联客户"
// @Param status query string false "配送状态"
// @Param number query string false "下单数量"
// @Param delivery query string false "配送周期"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Orders}} "{"code": 200, "data": [...]}"
// @Router /api/v1/orders [get]
// @Security Bearer
func (e Orders) GetPage(c *gin.Context) {
	req := dto.OrdersGetPageReq{}
	s := service.Orders{}
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
	list := make([]models.Orders, 0)
	var count int64
	req.CId = userDto.CId
	err = s.GetPage(e.getTableName(userDto.CId), &req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单失败,%s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Orders
// @Summary 获取Orders
// @Description 获取Orders
// @Tags Orders
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Orders} "{"code": 200, "data": [...]}"
// @Router /api/v1/orders/{id} [get]
// @Security Bearer
func (e Orders) Get(c *gin.Context) {
	req := dto.OrdersGetReq{}
	s := service.Orders{}
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
	var object models.Orders

	p := actions.GetPermissionFromContext(c)
	err = s.Get(e.getTableName(userDto.CId), &req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建Orders
// @Summary 创建Orders
// @Description 创建Orders
// @Tags Orders
// @Accept application/json
// @Product application/json
// @Param data body dto.OrdersInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/orders [post]
// @Security Bearer
func (e Orders) Insert(c *gin.Context) {
	req := dto.OrdersInsertReq{}
	s := service.Orders{}
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
	err = s.Insert(e.getTableName(userDto.CId), &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Orders
// @Summary 修改Orders
// @Description 修改Orders
// @Tags Orders
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrdersUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/orders/{id} [put]
// @Security Bearer
func (e Orders) Update(c *gin.Context) {
	req := dto.OrdersUpdateReq{}
	s := service.Orders{}
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
	err = s.Update(e.getTableName(userDto.CId), &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除Orders
// @Summary 删除Orders
// @Description 删除Orders
// @Tags Orders
// @Param data body dto.OrdersDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/orders [delete]
// @Security Bearer
func (e Orders) Delete(c *gin.Context) {
	s := service.Orders{}
	req := dto.OrdersDeleteReq{}
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

	err = s.Remove(e.getTableName(userDto.CId), &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

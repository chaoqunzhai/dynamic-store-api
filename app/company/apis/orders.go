package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
)

type Orders struct {
	api.Api
}

// 获取订单列表,
func (e Orders) GetPage(c *gin.Context) {
	req := dto.OrdersGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
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
	whereSql := fmt.Sprintf("c_id = ?", userDto.CId)
	//查询是否分表了

	p := actions.GetPermissionFromContext(c)
	var count int64

	var data models2.Orders

	list := make([]models2.Orders, 0)

	e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).Where(whereSql).
		Find(&list).Limit(-1).Offset(-1).
		Count(&count)

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e Orders) Valet(c *gin.Context) {

}
func (e Orders) Get(c *gin.Context) {

}
/**
@Author: chaoqun
* @Date: 2023/6/1 00:41
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	customUser "go-admin/common/jwt/user"
)

type IndexReq struct {
	Day string  `json:"day"`
}
type DetailReq struct {
	Id     int     `uri:"id" comment:"主键编码"`
}
type OrderShopResult struct {
	ShopId int
}
//获取指定日期的报表
//按配送员区分,每个配送员
//下订单是和商家关联的，而且商家都有一个关联的路线,所以反查即可
//是根据配送周期
func (e Orders) Index(c *gin.Context) {
	req := IndexReq{}
	err := e.MakeContext(c).
		Bind(&req).
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
	//根据选择的日期 + 大B配置的自定义配送时间
	orderTableName := e.getTableName(userDto.CId)
	//获取指定天数的订单的商家列表
	whereSql :=fmt.Sprintf("c_id = %v and enable = %v ",userDto.CId,true)
	orderResult:=make([]OrderShopResult,0)
	e.Orm.Table(orderTableName).Raw(whereSql).Scan(&orderResult)

	//通过商家列表获取到各路线





}
func (e Orders) Detail(c *gin.Context) {
	req := DetailReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}



}
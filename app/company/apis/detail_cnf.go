package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
)

type DetailCnf struct {
	api.Api
}
type DetailInsertReq struct {
	DetailAddCart string `json:"detail_add_cart" ` //详情页面中,加入购物车的文案
	DetailAddCartColor string `json:"detail_add_cart_color" ` //详情页面中,加入购物车的颜色
	DetailAddCartShow bool `json:"detail_add_cart_show"` //是否展示加入购物车
	DetailByNow string `json:"detail_by_now" ` //详情页面中,立即购买的文案
	DetailByNowColor string `json:"detail_by_now_color" ` //详情页面中,立即购买的文案
	DetailByNowShow bool `json:"detail_by_now_show"` //是否展示立即购买
	VisitorShowVip bool `json:"visitor_show_vip" ` //是否展示访问VIP价格

}
func (e DetailCnf) Create(c *gin.Context) {
	req := DetailInsertReq{}
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

	var WeAppExtendCnf models.WeAppExtendCnf
	e.Orm.Model(&models.WeAppExtendCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&WeAppExtendCnf)


	if WeAppExtendCnf.Id > 0 {

		WeAppExtendCnf.DetailAddCart = req.DetailAddCart
		WeAppExtendCnf.DetailAddCartColor = req.DetailAddCartColor
		WeAppExtendCnf.DetailByNowShow = req.DetailByNowShow
		WeAppExtendCnf.DetailByNowColor = req.DetailByNowColor
		WeAppExtendCnf.DetailAddCartShow = req.DetailAddCartShow
		WeAppExtendCnf.DetailByNow = req.DetailByNow
		WeAppExtendCnf.VisitorShowVip = req.VisitorShowVip
		e.Orm.Save(&WeAppExtendCnf)
	}else {
		object:=models.WeAppExtendCnf{
			DetailAddCart: req.DetailAddCart,
			DetailAddCartShow: req.DetailAddCartShow,
			DetailAddCartColor: req.DetailAddCartColor,
			DetailByNow: req.DetailByNow,
			DetailByNowColor: req.DetailByNowColor,
			DetailByNowShow: req.DetailByNowShow,
			VisitorShowVip: req.VisitorShowVip,
		}
		object.CId = userDto.CId
		object.Enable = true
		e.Orm.Create(&object)
	}
	e.OK("","successful")
	return
}
func (e DetailCnf) Detail(c *gin.Context) {
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

	var WeAppExtendCnf models.WeAppExtendCnf
	e.Orm.Model(&models.WeAppExtendCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&WeAppExtendCnf)



	e.OK(WeAppExtendCnf,"successful")
	return
}

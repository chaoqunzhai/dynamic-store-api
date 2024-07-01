package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
)

type DetailCnf struct {
	api.Api
}
type DetailInsertReq struct {
	SpecImageShow bool `json:"spec_image_show" gorm:"default:1"` //规格图片是否展示
	TitleLine int `json:"title_line" gorm:"size:1;default:1;comment:首页标题显示行数"`
	DetailAddName string `json:"detail_add_name" gorm:"size:5;comment:购物车名称"`
	DetailAddCart string `json:"detail_add_cart" ` //详情页面中,加入购物车的文案
	//DetailAddCartColor string `json:"detail_add_cart_color" ` //详情页面中,加入购物车的颜色
	DetailAddCartShow bool `json:"detail_add_cart_show"` //是否展示加入购物车
	DetailByNow string `json:"detail_by_now" ` //详情页面中,立即购买的文案
	//DetailByNowColor string `json:"detail_by_now_color" ` //详情页面中,立即购买的文案
	DetailByNowShow bool `json:"detail_by_now_show"` //是否展示立即购买
	VisitorShowVip bool `json:"visitor_show_vip" ` //是否展示访问VIP价格
	SaleShow bool `json:"sale_show" gorm:"default:0"`//销售量开关
	StockShow bool `json:"stock_show" gorm:"default:0"`//库存开关
	MinBuyShow bool `json:"min_buy_show" gorm:"default:0"`//起售量展示
	MaxBuyShow bool `json:"max_buy_show" gorm:"default:0"`//限购量展示
	ShowBarrageShow bool  `json:"show_barrage_show" gorm:"default:0"`//弹幕展示
	MarketPriceShow bool `json:"market_price_show" gorm:"default:1"`//市场价展示
	BuyingAuthList []int `json:"buying_auth_list"` //购买的权限列表
	PreviewVipList []int `json:"preview_vip_list"` //浏览的权限列表
	BuyingAuth bool `json:"buying_auth"`
	PriceShow bool `json:"price_show"` //是否展示售价
	Preview bool `json:"preview"`
	RecommendShow bool `json:"recommend_show"` //是否显示推荐产品
	ServerShow bool `json:"server_show"` //是否展示产品服务
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

	WeAppExtendCnf.TitleLine = req.TitleLine
	WeAppExtendCnf.DetailAddName = req.DetailAddName
	WeAppExtendCnf.DetailAddCart = req.DetailAddCart

	WeAppExtendCnf.DetailByNowShow = req.DetailByNowShow
	WeAppExtendCnf.DetailAddCartShow = req.DetailAddCartShow
	WeAppExtendCnf.DetailByNow = req.DetailByNow
	WeAppExtendCnf.VisitorShowVip = req.VisitorShowVip
	WeAppExtendCnf.SaleShow = req.SaleShow
	WeAppExtendCnf.StockShow = req.StockShow
	WeAppExtendCnf.ShowBarrageShow = req.ShowBarrageShow
	WeAppExtendCnf.MaxBuyShow = req.MaxBuyShow
	WeAppExtendCnf.SpecImageShow = req.SpecImageShow
	WeAppExtendCnf.MinBuyShow = req.MinBuyShow
	WeAppExtendCnf.MarketPriceShow = req.MarketPriceShow
	WeAppExtendCnf.Preview = req.Preview
	WeAppExtendCnf.BuyingAuth =req.BuyingAuth
	WeAppExtendCnf.PriceShow = req.PriceShow
	WeAppExtendCnf.RecommendShow = req.RecommendShow
	WeAppExtendCnf.ServerShow = req.ServerShow
	if WeAppExtendCnf.Id > 0 {

		e.Orm.Save(&WeAppExtendCnf)
		//更新后 应该存到redis中,因为是不常变化的值
		//redis_db.SetConfigManyInit(userDto.CId,  global.SmallBConfigExtendKey,WeAppExtendCnf)
	}else {
		WeAppExtendCnf.CId = userDto.CId
		WeAppExtendCnf.Enable = true
		e.Orm.Create(&WeAppExtendCnf)
		////更新后 应该存到redis中,因为是不常变化的值
		//redis_db.SetConfigManyInit(userDto.CId,  global.SmallBConfigExtendKey,object)
	}


	//等级购买权限开了
	if !req.BuyingAuth {

		for _,k:=range req.BuyingAuthList{
			var showObject models.VipShowEnable

			e.Orm.Model(&models.VipShowEnable{}).Where("c_id = ? and vip_id = ? and `type` = ? ",userDto.CId,k, global.GoodsAuthVip).Limit(1).Find(&showObject)
			if showObject.Id == 0 {
				showObject =models.VipShowEnable{
					CId: userDto.CId,
					VipId: k,
					Type: global.GoodsAuthVip,
					Enable: true,
				}
				e.Orm.Create(&showObject)
			}else {
				showObject.Enable = true
				e.Orm.Save(&showObject)
			}
		}
	}else{
		//先设置为false
		e.Orm.Model(&models.VipShowEnable{}).Where("c_id = ?  and `type` = ? ",userDto.CId,global.GoodsAuthVip).Updates(map[string]interface{}{
			"enable":false,
		})
	}
	//等级预览权限开了
	if !req.Preview {
		for _,k:=range req.PreviewVipList{
			var showObject models.VipShowEnable

			e.Orm.Model(&models.VipShowEnable{}).Where("c_id = ? and vip_id = ? and `type` = ? ",userDto.CId,k, global.GoodsPreview).Limit(1).Find(&showObject)

			if showObject.Id == 0 {
				showObject =models.VipShowEnable{
					CId: userDto.CId,
					VipId: k,
					Type: global.GoodsPreview,
					Enable: true,
				}
				e.Orm.Create(&showObject)
			}else {
				showObject.Enable = true
				e.Orm.Save(&showObject)
			}
		}
	}else {
		e.Orm.Model(&models.VipShowEnable{}).Where("c_id = ?  and `type` = ? ",userDto.CId,global.GoodsPreview).Updates(map[string]interface{}{
			"enable":false,
		})
	}

	e.OK("","操作成功")
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

	result:=make(map[string]interface{},0)

	result["detail"] = WeAppExtendCnf


	vipCnfList:=make([]models.VipShowEnable,0)
	e.Orm.Model(&models.VipShowEnable{}).Where("c_id = ? and enable = ?",userDto.CId,true).Find(&vipCnfList)

	buyingAuthList :=make([]int,0)
	previewVipList :=make([]int,0)
	for _,row:=range vipCnfList{
		var gradeVip models.GradeVip
		e.Orm.Model(&gradeVip).Select("id,name").Where("c_id = ? and id = ? and enable = ?",userDto.CId,row.VipId,true).Limit(1).Find(&gradeVip)
		if gradeVip.Id == 0 {
			continue
		}
		if row.Type == global.GoodsPreview{
			previewVipList = append(previewVipList,row.VipId)
		}
		if row.Type == global.GoodsAuthVip{
			buyingAuthList = append(buyingAuthList,row.VipId)
		}
	}
	result["preview_vip_list"] = previewVipList
	result["buying_auth_list"] = buyingAuthList
	e.OK(result,"操作成功")
	return
}

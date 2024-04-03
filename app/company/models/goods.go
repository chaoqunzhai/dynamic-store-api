package models

import (
	"go-admin/common/models"
)

type Goods struct {
	models.Model
	Layer     int          `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable    bool         `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId       int          `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	Name      string       `json:"name" gorm:"type:varchar(35);comment:商品名称"`
	SpecName  string       `gorm:"size:8;comment:规格命名,例如是:颜色,重量,系列"`
	Subtitle  string       `json:"subtitle" gorm:"type:varchar(8);comment:宣发文案"`
	Image     string       `json:"image" gorm:"type:varchar(155);comment:商品图片路径"`
	Quota     int          `json:"quota" gorm:"type:tinyint(1);comment:是否限购"`
	EnjoyVipSale bool `json:"enjoy_vip_sale" gorm:"comment:是否享受会员功能"`
	VipSale   bool         `json:"vipSale" gorm:"type:tinyint(1);comment:单独会员价"`
	Inventory int          `json:"inventory" gorm:"comment:库存"`
	Sale      int          `json:"sale" gorm:"comment:销量"`
	Money     string       `gorm:"size:30;comment:价格区间"`
	Recommend bool         `json:"recommend" gorm:"comment:是否推荐"`
	RubikCube bool         `json:"rubik_cube" gorm:"comment:支持魔方"`
	Tag       []GoodsTag   `gorm:"many2many:goods_mark_tag;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:tag_id;"`
	Class     []GoodsClass `gorm:"many2many:goods_mark_class;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:class_id;"`
	Brand []GoodsBrand `gorm:"many2many:goods_mark_brand;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:brand_id;"`

	models.ModelTime
	models.ControlBy
}

func (Goods) TableName() string {
	return "goods"
}

func (e *Goods) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Goods) GetId() interface{} {
	return e.Id
}

type GoodsDesc struct {
	models.Model
	CId int `gorm:"index;comment:大BID"`
	GoodsId int
	Desc      string //描述内容
}
func (GoodsDesc) TableName() string {
	return "goods_desc"
}
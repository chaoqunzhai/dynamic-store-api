package models

import (
	"go-admin/common/models"
)

type Goods struct {
	models.Model

	Layer     int          `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable    bool         `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc      string       `json:"desc" gorm:"type:varchar(200);comment:商品详情"`
	CId       int          `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	Name      string       `json:"name" gorm:"type:varchar(35);comment:商品名称"`
	Subtitle  string       `json:"subtitle" gorm:"type:varchar(8);comment:宣发文案"`
	Image     string       `json:"image" gorm:"type:varchar(155);comment:商品图片路径"`
	Quota     int          `json:"quota" gorm:"type:tinyint(1);comment:是否限购"`
	VipSale   bool         `json:"vipSale" gorm:"type:tinyint(1);comment:会员价"`
	Inventory int          `json:"inventory" gorm:"comment:库存"`
	Sale      int          `json:"sale" gorm:"comment:销量"`
	Money     string       `gorm:"size:12;comment:价格区间"`
	Tag       []GoodsTag   `gorm:"many2many:goods_mark_tag;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:tag_id;"`
	Class     []GoodsClass `gorm:"many2many:goods_mark_class;foreignKey:id;joinForeignKey:goods_id;references:id;joinReferences:class_id;"`
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

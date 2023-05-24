package models

import (
     
     
     
     
     "time"

	"go-admin/common/models"

)

type Goods struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable string `json:"enable" gorm:"type:tinyint(1);comment:开关"` 
    Desc string `json:"desc" gorm:"type:varchar(200);comment:商品详情"` 
    CId string `json:"cId" gorm:"type:bigint(20);comment:大BID"` 
    Name string `json:"name" gorm:"type:varchar(35);comment:商品名称"` 
    Subtitle string `json:"subtitle" gorm:"type:varchar(100);comment:副标题"` 
    Image string `json:"image" gorm:"type:varchar(155);comment:商品图片路径"` 
    Quota string `json:"quota" gorm:"type:tinyint(1);comment:是否限购"` 
    VipSale string `json:"vipSale" gorm:"type:tinyint(1);comment:会员价"` 
    Code string `json:"code" gorm:"type:varchar(50);comment:条形码"` 
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
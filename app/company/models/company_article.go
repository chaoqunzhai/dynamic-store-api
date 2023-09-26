package models

import (
     
     
     
     "go-admin/common/models"

)

type CompanyArticle struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable string `json:"enable" gorm:"type:tinyint(1);comment:开关"` 
    Desc string `json:"desc" gorm:"type:varchar(35);comment:描述信息"` 
    CId string `json:"cId" gorm:"type:bigint(20);comment:大BID"` 
    CoverImage string `json:"coverImage" gorm:"type:varchar(20);comment:封面图片"` 
    Title string `json:"title" gorm:"type:varchar(50);comment:文章标题"` 
    Class string `json:"class" gorm:"type:varchar(10);comment:文章分类"` 
    Context string `json:"context" gorm:"type:longtext;comment:文章内容"` 
    models.ModelTime
    models.ControlBy
}

func (CompanyArticle) TableName() string {
    return "company_article"
}

func (e *CompanyArticle) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyArticle) GetId() interface{} {
	return e.Id
}
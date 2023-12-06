package models

//公告

type Message struct {
	BigBRichGlobal
	Context string `json:"context" gorm:"size:60;comment:公告内容"`
	Link string `json:"link" gorm:"size:30;comment:链接地址"`
}
func (Message) TableName() string {
	return "company_message"
}


// 广告位 轮播图

type Ads struct {
	BigBRichGlobal
	Type int `json:"type"  gorm:"comment:类型: 0:顶部 1:中部"`
	LinkName string `json:"link_name" gorm:"size:30;comment:名称"`
	LinkUrl string `json:"link_url" gorm:"size:50;comment:链接地址"`
	ImageUrl string `json:"image_url" gorm:"size:30;comment:图片地址"`
	ShowImage string `json:"image" gorm:"-"` //展示字段
}
func (Ads) TableName() string {
	return "company_ads"
}
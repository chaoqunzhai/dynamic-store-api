package models

//文章信息
type Article struct {
	BigBRichGlobal
	CoverImage string `json:"cover_image" gorm:"size:20;comment:封面图片"`
	Title string `json:"title" gorm:"size:50;comment:文章标题"`
	Class string `json:"class" gorm:"size:10;comment:文章分类"`
	Context string `json:"context" gorm:"comment:文章内容"`
}
func (Article) TableName() string {
	return "company_article"
}
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
	LinkName string `json:"link_name"`
	LinkUrl string `json:"link_url"`
	ImageUrl string `json:"image_url"`
}
func (Ads) Ads() string {
	return "company_ads"
}
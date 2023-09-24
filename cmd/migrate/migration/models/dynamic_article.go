package models

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
/**
@Author: chaoqun
* @Date: 2023/12/19 17:44
*/
package models
//全局生效配置

type GlobalArticle struct {
	RichGlobal
	Name      string       `gorm:"size:50;comment:标题名称"`
	Subtitle string   `gorm:"size:120;comment:副标题名称"`
	Type int  `json:"type" gorm:"size:1;default:1;index;comment:1:公告 2:帮助文档"`
	Link string `json:"link" gorm:"size:120;comment:跳转地址"`

}
func (GlobalArticle) TableName() string {
		return "global_article"
}
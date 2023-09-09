package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
)

type Theme struct {
	api.Api
}
type ThemeInsertReq struct {
	Theme string `json:"theme"`


}
func (e Theme) Create(c *gin.Context) {
	req := ThemeInsertReq{}
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

	var PayCnf models.WeAppExtendCnf
	e.Orm.Model(&models.WeAppExtendCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&PayCnf)


	if PayCnf.Id > 0 {

		PayCnf.StyleTheme = req.Theme
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.WeAppExtendCnf{
			StyleTheme: req.Theme,
		}
		trade.CId = userDto.CId
		trade.Enable = true
		e.Orm.Create(&trade)
	}
	e.OK("","successful")
	return
}
func (e Theme) Detail(c *gin.Context) {
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
	var PayCnf models.WeAppExtendCnf
	e.Orm.Model(&models.WeAppExtendCnf{}).Select("style_theme").Where("c_id = ?",userDto.CId).Limit(1).Find(&PayCnf)


	e.OK(PayCnf.StyleTheme,"successful")
	return
}

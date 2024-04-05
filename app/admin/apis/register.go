package apis

import (
	"github.com/gin-gonic/gin"
	models2 "go-admin/cmd/migrate/migration/models"
)

type RegisterReq struct {
	Desc    string    `json:"desc" `
	UserName    string    `json:"username" `
	Phone    string    `json:"phone" `
}
func (e GoAdminSystem)Register(c *gin.Context) {
	req:=RegisterReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	var object models2.CompanyRegister
	e.Orm.Model(&object).Where("phone = ?",req.Phone).Limit(1).Find(&object)
	if object.Id > 0 {

		e.Error(500, err, "您已经申请过,请耐心等待反馈")
		return
	}

	e.Orm.Create(&models2.CompanyRegister{
		UserName: req.UserName,
		Phone: req.Phone,
		Desc: req.Desc,
	})
	e.OK("","")
	return
}
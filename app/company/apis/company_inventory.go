package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
)

type CompanyInventory struct {
	api.Api
}

func (e CompanyInventory) GetPage(c *gin.Context) {
	req := dto.CompanyMessageGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

}

func (e CompanyInventory) Info(c *gin.Context) {

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


	var object models2.InventoryCnf
	e.Orm.Model(&models2.InventoryCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)

	if object.Id == 0 {
		e.OK(false,"")
		return
	}

	e.OK(object.Enable,"")
	return
}
func (e CompanyInventory) UpdateCnf(c *gin.Context) {
	req := dto.CompanyInventoryCnfReq{}
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

	var object models2.InventoryCnf
	e.Orm.Model(&models2.InventoryCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)

	if object.Id == 0 {
		row:=models2.InventoryCnf{}
		row.CId = userDto.CId
		row.Enable = req.Enable
		e.Orm.Create(&row)
		return
	}
	e.Orm.Model(&models2.InventoryCnf{}).Where("c_id = ?",userDto.CId).Updates(map[string]interface{}{
		"enable":req.Enable,
	})
	e.OK(object.Enable,"操作成功")
	return
}
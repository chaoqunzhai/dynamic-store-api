/**
@Author: chaoqun
* @Date: 2023/9/25 23:49
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
)


type OfflinePay struct {
	api.Api
}
type OfflinePayReq struct {
	Name string `json:"name" `
}
func (e *OfflinePay) Create(c *gin.Context) {
	req := OfflinePayReq{}
	err := e.MakeContext(c).
		Bind(&req, binding.JSON, nil).
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
	var offPay models.OfflinePay
	var count int64
	e.Orm.Model(&offPay).Where("c_id = ? and name = ?",userDto.CId,req.Name).Count(&count)
	if count > 0 {
		e.Error(500, nil, "已经存在")
		return
	}
	offPay = models.OfflinePay{
		Name: req.Name,
	}
	offPay.CreateBy = userDto.UserId
	offPay.CId = userDto.CId
	e.Orm.Create(&offPay)
	return
}

func (e *OfflinePay) List(c *gin.Context) {
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
	list:=make([]models.OfflinePay,0)
	e.Orm.Model(&models.OfflinePay{}).Where("c_id = ?",userDto.CId).Find(&list)

	result:=make([]map[string]interface{},0)
	for _,row:=range list{
		result = append(result, map[string]interface{}{
			"name":row.Name,
			"create_at":row.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	e.OK(result,"successful")
	return
}


func (e *OfflinePay) Update(c *gin.Context) {
	req := OfflinePayReq{}
	err := e.MakeContext(c).
		Bind(&req, binding.JSON, nil).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	uid:=c.Query("id")

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	e.Orm.Model(&models.OfflinePay{}).Where("c_id = ? and id = ?",userDto.CId,uid).Updates(map[string]interface{}{
		"name":req.Name,
	})

	e.OK("","successful")
	return
}
func (e *OfflinePay) Remove(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	uid:=c.Query("id")
	fmt.Println("uid",uid)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//直接删除
	e.Orm.Model(&models.OfflinePay{}).Unscoped().Where("c_id = ? and id = ？",userDto.CId,uid).Delete(&models.OfflinePay{})

	e.OK("","successful")
	return
}

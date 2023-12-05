package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
)

type CompanyMessAge struct {
	api.Api
}

func (e CompanyMessAge) GetPage(c *gin.Context) {
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

	p := actions.GetPermissionFromContext(c)
	list := make([]models2.Message, 0)
	var count int64

	var data models2.Message

	err = e.Orm.Model(&data).Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")


}

func (e CompanyMessAge) Enable(c *gin.Context) {

	req := dto.UpdateEnableReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).

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
	e.Orm.Model(&models2.Message{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Updates(map[string]interface{}{
		"enable":req.Enable,
	})
	e.OK("", "更新成功")

}

func (e CompanyMessAge) Insert(c *gin.Context) {
	req := dto.CompanyMessageInsertReq{}
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
	req.SetCreateBy(userDto.UserId)
	req.CId = userDto.CId
	req.Enable = true
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "index_message", e.Orm)
	MaxNumber := CompanyCnf["index_message"]
	var count int64
	e.Orm.Model(&models2.Message{}).Where("c_id = ?",userDto.CId).Count(&count)
	if count >= int64(MaxNumber) {
		msg:=fmt.Sprintf("最多只能创建%v条消息",MaxNumber)
		e.Error(500, errors.New(msg), msg)
		return
	}
	dat :=models2.Message{
		Link: req.Link,
		Context: req.Context,
	}
	dat.Layer = 0
	dat.Desc = req.Desc
	dat.CId = userDto.CId
	dat.CreateBy = userDto.UserId
	dat.Enable = true
	e.Orm.Create(&dat)

	e.OK("", "创建成功")
}


func (e CompanyMessAge) Update(c *gin.Context) {
	req := dto.CompanyMessageUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).

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

	var count int64
	e.Orm.Model(&models2.Message{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}

	e.Orm.Model(&models2.Message{}).Where("id = ? and c_id = ?",req.Id,userDto.CId).Updates(map[string]interface{}{
		"context":req.Context,
		"link":req.Link,
		"desc":req.Desc,
		"update_by":userDto.UserId,
	})

	e.OK("", "修改成功")

}


func (e CompanyMessAge) Delete(c *gin.Context) {
	req := dto.CompanyMessageDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)


	var data models2.Message

	e.Orm.Model(&data).Scopes(
			actions.Permission(data.TableName(), p),
		).Unscoped().Delete(&data, req.GetId())
	e.OK("", "删除成功")
}
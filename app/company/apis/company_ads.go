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

type CompanyAds struct {
	api.Api
}

func (e CompanyAds) GetPage(c *gin.Context) {
	req := dto.CompanyAdsGetPageReq{}
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
	list := make([]models2.Ads, 0)
	var count int64

	var data models2.Ads

	err = e.Orm.Model(&data).Scopes(
		cDto.MakeCondition(req.GetNeedSearch()),
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		actions.Permission(data.TableName(), p),
	).
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
}

func (e CompanyAds) Enable(c *gin.Context) {

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
	e.Orm.Model(&models2.Ads{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Updates(map[string]interface{}{
		"enable":req.Enable,
	})
	e.OK("", "更新成功")

}

func (e CompanyAds) Insert(c *gin.Context) {
	req := dto.CompanyAdsInsertReq{}
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
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "index_ads", e.Orm)
	MaxNumber := CompanyCnf["index_ads"]
	var count int64
	e.Orm.Model(&models2.Ads{}).Where("c_id = ?",userDto.CId).Count(&count)
	if count >= int64(MaxNumber) {
		msg:=fmt.Sprintf("最多只能创建%v条广告",MaxNumber)
		e.Error(500, errors.New(msg), msg)
		return
	}
	dat :=models2.Ads{
		LinkName: req.LinkName,
		LinkUrl: req.LinkUrl,
		ImageUrl: req.ImageUrl,
	}
	dat.Layer = 0
	dat.Desc = req.Desc
	dat.CId = userDto.CId
	dat.CreateBy = userDto.UserId
	dat.Enable = true
	e.Orm.Create(&dat)

	e.OK("", "创建成功")
}


func (e CompanyAds) Update(c *gin.Context) {
	req := dto.CompanyAdsUpdateReq{}
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
	e.Orm.Model(&models2.Ads{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}

	e.Orm.Model(&models2.Ads{}).Where("id = ? and c_id = ?",req.Id,userDto.CId).Updates(map[string]interface{}{
		"link_url":req.LinkUrl,
		"link_name":req.LinkName,
		"desc":req.Desc,
		"image_url":req.ImageUrl,
	})

	e.OK("", "修改成功")

}


func (e CompanyAds) Delete(c *gin.Context) {
	req := dto.CompanyAdsDeleteReq{}
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
	var data models2.Ads
	e.Orm.Model(&data).Scopes(
		actions.Permission(data.TableName(), p),
	).Limit(1).Find(&data)
	//删除图片

	//删除数据
	//e.Orm.Model(&data).Scopes(
	//	actions.Permission(data.TableName(), p),
	//).Unscoped().Delete(&data, req.GetId())
	e.OK("", "删除成功")
}
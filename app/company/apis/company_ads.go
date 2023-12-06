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
	"go-admin/common/qiniu"
	utils2 "go-admin/common/utils"
	"go-admin/global"
	"go.uber.org/zap"
	"os"
	"strings"
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

	result:=make([]interface{},0)
	for _,r:=range list{
		r.ShowImage = business.GetGoodsPathFirst(r.CId,r.ImageUrl,global.AdsPath)
		result = append(result,
			r)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

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
		Type: req.Type,
	}
	dat.Layer = 0
	dat.Desc = req.Desc
	dat.CId = userDto.CId
	dat.CreateBy = userDto.UserId
	dat.Enable = true
	saveErr :=e.Orm.Create(&dat).Error
	if saveErr == nil {
		//存储图片

		fileObj,fileErr :=c.FormFile("files")
		if fileErr!=nil{
			//没有文件直接返回
			e.OK("", "创建成功")
			return
		}
		buckClient :=qiniu.QinUi{CId: userDto.CId}
		buckClient.InitClient()
		_,goodsImagePath  :=GetCosImagePath(global.AdsPath,fileObj.Filename,userDto.CId)
		if saveFileErr := c.SaveUploadedFile(fileObj, goodsImagePath); saveFileErr == nil {

			//1.上传到cos中
			fileName,cosErr :=buckClient.PostFile(goodsImagePath)
			if cosErr !=nil{
				zap.S().Errorf("用户:%v,CID:%v 广告图片上传失败:%v",userDto.UserId,userDto.CId,cosErr)
				e.OK("", "创建成功")
				return
			}
			//本地删除
			_=os.RemoveAll(goodsImagePath)
			e.Orm.Model(&models2.Ads{}).Where("id = ?",  dat.Id).Updates(map[string]interface{}{
				"image_url": fileName,
			})
		}


	}

	e.OK("", "创建成功")
}


func (e CompanyAds) Update(c *gin.Context) {
	req := dto.CompanyAdsUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)

		e.Error(500, err, err.Error())
		return
	}
	if bindErr := c.ShouldBind(&req); bindErr != nil {

		e.Error(500, bindErr, bindErr.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	uid := c.Param("id")

	var adsObject models2.Ads
	e.Orm.Model(&models2.Ads{}).Where("id = ?", uid).Limit(1).Find(&adsObject)
	if adsObject.Id == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}

	baseFileList := make([]string, 0)
	if adsObject.ImageUrl != "" {
		baseFileList = strings.Split(adsObject.ImageUrl, ",")
	}

	buckClient:=qiniu.QinUi{
		CId: userDto.CId,
	}
	buckClient.InitClient()
	//原值
	imageUrl :=adsObject.ImageUrl
	if req.FileClear == 1 {
		if adsObject.ImageUrl != "" {
			for _, image := range strings.Split(adsObject.ImageUrl, ",") {
				buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId, global.AdsPath, image))
			}
		}

		imageUrl = ""
	}else {

		fileList := make([]string, 0)
		fileObj,_ :=c.FormFile("files")
		//处理下路径
		if req.BaseFiles != "" {
			for _, baseFile := range strings.Split(req.BaseFiles, ",") {
				ll := strings.Split(baseFile, "/")
				fileList = append(fileList, ll[len(ll)-1])
			}
		}
		//前段更新了,进行文件内容的比对 baseFileList 和 fileList 比对，如果不一样是需要进行删除的
		diffList := utils2.Difference(baseFileList, fileList)

		for _, image := range diffList {

			buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId,global.AdsPath,image))
		}

		//正常获取文件对象
		if fileObj != nil{
			_,goodsImagePath  :=GetCosImagePath(global.AdsPath,fileObj.Filename,userDto.CId)
			if saveFileErr := c.SaveUploadedFile(fileObj, goodsImagePath); saveFileErr == nil {

				//1.上传到cos中
				fileName,cosErr :=buckClient.PostFile(goodsImagePath)
				if cosErr !=nil{
					zap.S().Errorf("用户:%v,CID:%v 广告图片上传失败:%v",userDto.UserId,userDto.CId,cosErr)
					e.OK("", "创建成功")
					return
				}
				//本地删除
				_=os.RemoveAll(goodsImagePath)
				//新图片名称
				imageUrl = fileName
			}
		}

	}
	e.Orm.Model(&models2.Ads{}).Where("id = ? and c_id = ?",uid,userDto.CId).Updates(map[string]interface{}{
		"link_url":req.LinkUrl,
		"link_name":req.LinkName,
		"desc":req.Desc,
		"image_url":imageUrl,
		"type":req.Type,
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var data models2.Ads
	e.Orm.Model(&data).Scopes(
		actions.Permission(data.TableName(), p),
	).Limit(1).Find(&data)
	//删除图片
	if data.ImageUrl != ""{
		//删除cos存储

		buckClient:=qiniu.QinUi{
			CId: userDto.CId,
		}
		buckClient.InitClient()
		buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId,global.AdsPath,data.ImageUrl))
	}
	//直接删除数据
	e.Orm.Model(&data).Scopes(
		actions.Permission(data.TableName(), p),
	).Unscoped().Delete(&data, req.Id)
	e.OK("", "删除成功")
}
package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"github.com/google/uuid"
	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"time"
)

type Goods struct {
	api.Api
}

type ClassData struct {
	ClassId   int        `json:"class_id" `
	ClassName string     `json:"class_name" `
	GoodsList []specsRow `json:"goods_list" `
}

type specsRow struct {
	GoodsId   int     `json:"goods_id" `
	Money     float64 `json:"money" `
	Unit      string  `json:"unit" `
	Name      string  `json:"name" `
	Inventory int     `json:"inventory" ` //库存
}

func (e Goods) ClassSpecs(c *gin.Context) {

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
	var goods []models.Goods
	e.Orm.Model(&models.Goods{}).Where("c_id = ? and enable = ?", userDto.CId, true).
		Order(global.OrderLayerKey).Preload("Class", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name")
	}).Find(&goods)

	result := make(map[int]ClassData, 0)
	for _, row := range goods {

		var specsObject models.GoodsSpecs
		e.Orm.Model(&models.GoodsSpecs{}).Where("c_id = ? and enable = ? and goods_id = ?", userDto.CId, true, row.Id).Limit(1).Find(&specsObject)
		if specsObject.Id == 0 {
			continue
		}
		specData := specsRow{
			GoodsId:   row.Id,
			Money:     specsObject.Price,
			Unit:      specsObject.Unit,
			Name:      specsObject.Name,
			Inventory: specsObject.Inventory,
		}
		for _, class := range row.Class {
			data, ok := result[class.Id]
			if ok {
				data.GoodsList = append(data.GoodsList, specData)
				result[class.Id] = data
			} else {
				result[class.Id] = ClassData{
					ClassId:   class.Id,
					ClassName: class.Name,
					GoodsList: []specsRow{specData},
				}
			}
		}
	}

	e.OK(result, "successful")
	return

}

// GetPage 获取Goods列表
// @Summary 获取Goods列表
// @Description 获取Goods列表
// @Tags Goods
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "商品名称"
// @Param vipSale query string false "会员价"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Goods}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods [get]
// @Security Bearer
func (e Goods) GetPage(c *gin.Context) {
	req := dto.GoodsGetPageReq{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Goods, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Goods失败，\r\n失败信息 %s", err.Error()))
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, row := range list {
		r := map[string]interface{}{
			"id":     row.Id,
			"name":   row.Name,
			"enable": row.Enable,
			"layer":  row.Layer,
			"class": func() []string {
				cache := make([]string, 0)
				for _, cl := range row.Class {
					cache = append(cache, cl.Name)
				}
				return cache
			}(),
			"image": func() string {
				if row.Image == "" {
					return ""
				}
				return business.GetGoodPathName(row.CId) + strings.Split(row.Image, ",")[0]
			}(),
			"created_at": row.CreatedAt,
		}
		result = append(result, r)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Goods
// @Summary 获取Goods
// @Description 获取Goods
// @Tags Goods
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Goods} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods/{id} [get]
// @Security Bearer
func (e Goods) Get(c *gin.Context) {
	req := dto.GoodsGetReq{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
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
	var object models.Goods

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Goods失败，\r\n失败信息 %s", err.Error()))
		return
	}
	goodsMap := map[string]interface{}{
		"name":     object.Name,
		"subtitle": object.Subtitle,
		"desc":     object.Desc,
		"tag": func() []int {
			t := make([]int, 0)
			for _, r := range object.Tag {
				t = append(t, r.Id)
			}
			return t
		}(),
		"class": func() []int {
			t := make([]int, 0)
			for _, r := range object.Class {
				t = append(t, r.Id)
			}
			return t
		}(),
		"vip_sale": object.VipSale,
		"quota":    object.Quota,
		"code":     object.Code,
		"enable":   object.Enable,
		"layer":    object.Layer,
		"image": func() []string {
			i := make([]string, 0)
			if object.Image == "" {
				return i
			}
			for _, im := range strings.Split(object.Image, ",") {
				i = append(i,
					business.GetGoodPathName(object.CId)+im)
			}
			return i
		}(),
	}
	var specsList []models.GoodsSpecs
	e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id = ? and c_id = ?", req.Id, userDto.CId).Find(&specsList)
	specData := make([]interface{}, 0)
	specVipData := make([]interface{}, 0)

	for _, specs := range specsList {
		now := time.Now().Unix()
		specRow := map[string]interface{}{
			"id":        specs.Id,
			"key":       now,
			"name":      specs.Name,
			"price":     specs.Price,
			"original":  specs.Original,
			"inventory": specs.Inventory,
			"limit":     specs.Limit,
			"enable":    specs.Enable,
			"layer":     specs.Layer,
			"unit":      specs.Unit,
		}
		specData = append(specData, specRow)
		vipMap := map[string]interface{}{
			"key":    now,
			"name":   specs.Name,
			"price":  specs.Price,
			"enable": specs.Enable,
		}
		var specVipList []models.GoodsVip
		e.Orm.Model(&models.GoodsVip{}).Where("specs_id = ? and c_id = ?", specs.Id, userDto.CId).Find(&specVipList)

		for _, vip := range specVipList {
			vipKey := fmt.Sprintf("vip_%v", vip.GradeId)
			vipMap[vipKey] = vip.CustomPrice
		}
		specVipData = append(specVipData, vipMap)

	}

	goodsMap["specs"] = specData
	goodsMap["specsVip"] = specVipData
	e.OK(goodsMap, "查询成功")
}

func (e Goods) UpdateState(c *gin.Context) {
	req := dto.GoodsStateReq{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
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
	for _, row := range req.Goods {
		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?", row, userDto.CId).Updates(map[string]interface{}{
			"enable": req.Enable,
		})
	}
	e.OK("更新成功", "更新成功")
	return
}

// Insert 创建Goods
// @Summary 创建Goods
// @Description 创建Goods
// @Tags Goods
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods [post]
// @Security Bearer
func (e Goods) Insert(c *gin.Context) {
	req := dto.GoodsInsertReq{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var count int64
	e.Orm.Model(&models.Goods{}).Where("c_id = ? and name = ?", userDto.CId, req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}

	goodId, goodErr := s.Insert(userDto.CId, &req)
	if goodErr != nil {
		e.Error(500, err, fmt.Sprintf("创建商品失败,%s", goodErr.Error()))
		return
	}
	//商品信息创建成功,才会保存客户的商品照片
	goodsImagePath := business.GetGoodPathName(fmt.Sprintf("%v", userDto.CId))
	// 遍历所有图片
	fileForm, fileErr := c.MultipartForm()
	if fileErr != nil {
		e.Error(500, nil, "请提交表单模式")
		return
	}
	files := fileForm.File["files"]
	for _, file := range files {
		// 逐个存
		guid := strings.Split(uuid.New().String(), "-")
		filePath := guid[0] + utils.GetExt(file.Filename)
		saveFilePath := goodsImagePath + filePath
		fileList := make([]string, 0)
		if saveErr := c.SaveUploadedFile(file, saveFilePath); saveErr == nil {
			//只保留文件名称,防止透露服务器地址
			fileList = append(fileList, filePath)
		}
		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?", goodId, userDto.CId).Updates(map[string]interface{}{
			"image": strings.Join(fileList, ","),
		})
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Goods
// @Summary 修改Goods
// @Description 修改Goods
// @Tags Goods
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods/{id} [put]
// @Security Bearer
func (e Goods) Update(c *gin.Context) {
	req := dto.GoodsUpdateReq{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
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
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	putUid := c.Param("id")
	//手动设置下数据ID
	uid, _ := strconv.Atoi(putUid)
	req.Id = uid
	var count int64
	e.Orm.Model(&models.Goods{}).Where("id = ?", req.Id).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}
	var oldRow models.Goods
	e.Orm.Model(&models.Goods{}).Where("name = ? and c_id = ?", req.Name, userDto.CId).Limit(1).Find(&oldRow)

	if oldRow.Id != 0 {
		if oldRow.Id != req.Id {
			e.Error(500, errors.New("名称不可重复"), "名称不可重复")
			return
		}
	}
	err = s.Update(userDto.CId, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改商品信息失败,%s", err.Error()))
		return
	}

	if req.FileClear == 1 {
		var goodsObject models.Goods
		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?",
			req.Id, userDto.CId).Limit(1).Find(&goodsObject)

		if goodsObject.Image != "" {
			for _, image := range strings.Split(goodsObject.Image, ",") {
				os.Remove(business.GetGoodPathName(userDto.CId) + image)
			}
		}
		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?",
			req.Id, userDto.CId).Updates(map[string]interface{}{
			"image": "",
		})
	} else {

		//商品信息创建成功,才会保存客户的商品照片
		goodsImagePath := business.GetGoodPathName(fmt.Sprintf("%v", userDto.CId))
		// 遍历所有图片
		fileForm, fileErr := c.MultipartForm()
		if fileErr != nil {
			e.Error(500, nil, "请提交表单模式")
			return
		}
		files := fileForm.File["files"]
		for _, file := range files {
			// 逐个存
			guid := strings.Split(uuid.New().String(), "-")
			filePath := guid[0] + utils.GetExt(file.Filename)
			saveFilePath := goodsImagePath + filePath
			fileList := make([]string, 0)
			if saveErr := c.SaveUploadedFile(file, saveFilePath); saveErr == nil {
				//只保留文件名称,防止透露服务器地址
				fileList = append(fileList, filePath)
			}

			e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?", req.Id, userDto.CId).Updates(map[string]interface{}{
				"image": strings.Join(fileList, ","),
			})
		}
	}

	e.OK(req.GetId(), "修改成功")
}

// Delete 删除Goods
// @Summary 删除Goods
// @Description 删除Goods
// @Tags Goods
// @Param data body dto.GoodsDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods [delete]
// @Security Bearer
func (e Goods) Delete(c *gin.Context) {
	s := service.Goods{}
	req := dto.GoodsDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除Goods失败,%s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

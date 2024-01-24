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
	"go-admin/common/qiniu"
	utils2 "go-admin/common/utils"
	"go-admin/global"
	"go.uber.org/zap"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
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
	GoodsId    int     `json:"goods_id" `
	GoodsName  string  `json:"goods_name"`
	GoodsPrice string  `json:"goods_price"`
	GoodsStore int     `json:"goods_store"`
	Image      string  `json:"image" `
	Money      float64 `json:"money" `
	Unit       string  `json:"unit" `
	Name       string  `json:"name" `
	Inventory  int     `json:"inventory" ` //库存
}

func (e Goods) ClassSpecs(c *gin.Context) {
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
	ClassID:=c.Query("class_id")
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if ClassID == "" {
		e.Error(500, err, "请输入类别")
		return
	}
	result := make([]specsRow, 0)
	//获取商品和分类的关联
	var bindGoodsId []int
	e.Orm.Raw(fmt.Sprintf("select goods_id from goods_mark_class where class_id = %v",ClassID)).Scan(&bindGoodsId)

	if len(bindGoodsId) == 0 {
		e.OK(result, "successful")
		return
	}
	var goods []models.Goods
	e.Orm.Model(&models.Goods{}).Select("id,name,c_id,image,inventory,money").Where(
		"c_id = ? and enable = ? and id in ?", userDto.CId, true,bindGoodsId).
		Order(global.OrderLayerKey).Find(&goods).Limit(-1).Offset(-1)


	goodsId:=make([]int,0)
	for _, row := range goods {
		goodsId = append(goodsId,row.Id)

	}
	goodsId = utils2.RemoveRepeatInt(goodsId)
	openInventory,InventoryMap:=s.GetBatchGoodsInventory(userDto.CId,goodsId)

	for _, row := range goods {
		//只返回有规格的数据
		var specsObject models.GoodsSpecs
		e.Orm.Model(&models.GoodsSpecs{}).Scopes(actions.PermissionSysUser(specsObject.TableName(),userDto)).Where(
			"enable = ? and goods_id = ?", true, row.Id).Limit(1).Find(&specsObject)
		if specsObject.Id == 0 {
			continue
		}
		var Inventory int //规格总数
		if openInventory {
			Inventory = InventoryMap[row.Id]
		}else {
			Inventory = row.Inventory
		}
		var unitObject models.GoodsUnit
		unitName:=""
		if specsObject.UnitId > 0 {
			e.Orm.Model(&unitObject).Where("id = ? and c_id = ?",specsObject.UnitId,userDto.CId).Limit(1).Find(&unitObject)
			if unitObject.Id > 0 {
				unitName = unitObject.Name
			}
		}
		specData := specsRow{
			GoodsId: row.Id,
			Image: func() string {
				if row.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(row.CId,row.Image,global.GoodsPath)
			}(),
			GoodsName:  row.Name,
			GoodsPrice: row.Money,
			GoodsStore: Inventory,
			Money:      specsObject.Price,
			Unit:       unitName,
			Name:       specsObject.Name,
			Inventory:  Inventory,
		}
		result = append(result,specData)
	}
	e.OK(result, "successful")
	return

}
func getFileName(fileName string) string {
	guid := strings.Split(uuid.New().String(), "-")

	return  guid[0] + utils.GetExt(fileName)
}
func GetCosImagePath(imageConst,fileName string,CId interface{})  (filePath,goodsImagePath string) {

	//增加一层 cache_image 目录,防止因为大量的客户 产生大量的客户目录文件 堆放在程序目录同层级中

	goodsImagePath = path.Join("cache_image",business.GetSiteCosPath(CId,imageConst,getFileName(fileName)))

	return
}

//todo:存储商品详情中上传的图片
func (e Goods) CosSaveImage(c *gin.Context) {
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
	res :=make(map[string]interface{},0)
	file,_ :=c.FormFile("file")
	guid := strings.Split(uuid.New().String(), "-")

	filePath := guid[0] + utils.GetExt(file.Filename)

	goodsImagePath := business.GetSiteCosPath(userDto.CId,global.GoodsPath,filePath)
	fmt.Println("goodsImagePath",goodsImagePath)
	////1.文件先存本地
	if saveErr :=c.SaveUploadedFile(file,goodsImagePath);saveErr==nil{
		//2.上传到cos中
		cos :=qiniu.QinUi{CId: userDto.CId}
		cos.InitClient()
		fileName,cosErr:=cos.PostImageFile(goodsImagePath)
		fmt.Println("七牛保存的返回",fileName,cosErr)
		if cosErr !=nil{
			zap.S().Errorf("商品图片上传COS失败:%v",cosErr.Error())
			res["code"] = -1
			res["msg"] = "文件上传失败"
			e.OK(res,"")
			return
		}
		////3.上传成功后删除本地文件
		res["code"] = 0
		res["msg"] = "文件上传成功"
		res["url"] = business.GetDomainCosEncodePathName(global.GoodsPath,userDto.CId,fileName,false)
		_=os.Remove(goodsImagePath)

	}else {
		zap.S().Errorf("商品图片上传,本地保存图片失败:%v",saveErr.Error())
		res["code"] = -1
		res["msg"] = "文件上传失败"
	}
	zap.S().Infof("用户%v 上传图片成功,url:%v",userDto.Username,res["url"])
	e.OK(res,"")
	return
}
func (e Goods) CosRemoveImage(c *gin.Context) {
	req:=dto.GoodsRemove{}
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
	//fmt.Println("删除的图片",req.Image)
	//处理下encode的路径
	QueryUnescape,_ :=url.QueryUnescape(req.Image)

	cosImagePath:=business.GetDomainSplitFilePath(QueryUnescape)
	buckClient :=qiniu.QinUi{CId: userDto.CId}
	buckClient.InitClient()

	buckClient.RemoveFile(cosImagePath)
	e.OK("","successful")
	return
}
func (e Goods) MiniApi(c *gin.Context) {
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
	datalist := make([]models.Goods, 0)
	e.Orm.Model(&models.Goods{}).Select("id,name").Where("c_id = ? and enable = ?", userDto.CId, true).Order(global.OrderLayerKey).Find(&datalist)

	result := make([]map[string]interface{}, 0)
	for _, row := range datalist {
		result = append(result, map[string]interface{}{
			"id":   row.Id,
			"name": row.Name,
		})
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	listType := c.Query("listType")

	switch listType {
	case "on_sale":
		//在售中
		req.Enable = "1"
	case "off_sale":
		//下架的
		req.Enable = "0"
	case "sale_out":
		//售罄 ?
	}
	p := actions.GetPermissionFromContext(c)
	list := make([]models.Goods, 0)
	var count int64
	req.CId = userDto.CId
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Goods失败,信息 %s", err.Error()))
		return
	}

	//make数据
	result := make([]map[string]interface{}, 0)

	goodsId:=make([]int,0)
	for _, row := range list {
		goodsId = append(goodsId,row.Id)

	}
	goodsId = utils2.RemoveRepeatInt(goodsId)
	openInventory,InventoryMap:=s.GetBatchGoodsInventory(userDto.CId,goodsId)
	for _, row := range list {
		var Inventory int
		if openInventory{
			Inventory = InventoryMap[row.Id]
		}else {
			Inventory = row.Inventory
		}
		r := map[string]interface{}{
			"id":       row.Id,
			"name":     row.Name,
			"subtitle": row.Subtitle,
			"enable":   row.Enable,
			"layer":    row.Layer,
			"class": func() []string {
				cache := make([]string, 0)
				for _, cl := range row.Class {
					cache = append(cache, cl.Name)
				}
				return cache
			}(),
			"inventory": Inventory,
			"image": func() string {
				if row.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(row.CId,row.Image,global.GoodsPath)
			}(),
			"sale":       row.Sale,
			"created_at": row.CreatedAt,
			"spec_name":row.SpecName,
			//规格的价格从小到大
			"money": row.Money,
		}
		result = append(result, r)
	}

	data :=map[string]interface{}{
		"list":result,
		"count":count,
		"pageIndex": req.GetPageIndex(),
		"pageSize":req.GetPageSize(),
		"imageUrl":"cos",
	}
	e.OK(data,"查询成功")
	return
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
	var GoodsDesc models.GoodsDesc
	e.Orm.Model(&GoodsDesc).Where("goods_id = ?",req.Id).Limit(1).Find(&GoodsDesc)
	goodsMap := map[string]interface{}{
		"name":     object.Name,
		"subtitle": object.Subtitle,
		"desc":     GoodsDesc.Desc,
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
		"brand": func() []int {
			t := make([]int, 0)
			for _, r := range object.Brand {
				t = append(t, r.Id)
			}
			return t
		}(),
		"vip_sale": object.VipSale,
		"quota":    object.Quota,
		"enable":   object.Enable,
		"layer":    object.Layer,
		"spec_name":object.SpecName,
		"recommend":object.Recommend,
		"rubik_cube":object.RubikCube,
		"image": func() []map[string]string {
			i := make([]map[string]string, 0)
			if object.Image == "" {
				return i
			}
			for _, im := range strings.Split(object.Image, ",") {

				imagePath :=business.GetDomainCosEncodePathName(global.GoodsPath,object.CId,im,false)
				i = append(i, map[string]string{
					"url":imagePath,
					"name":im,
				})
			}
			return i
		}(),
	}
	var specsList []models.GoodsSpecs
	e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id = ? and c_id = ?", req.Id, userDto.CId).Find(&specsList)
	specData := make([]interface{}, 0)
	specVipData := make([]interface{}, 0)



	for _, specs := range specsList {
		now := utils.GetUUID()
		var inventory int
		openInventory,getInventory:=s.GetSpecInventory(object.CId,fmt.Sprintf("(goods_id = %v and spec_id = %v)",specs.GoodsId,specs.Id))

		if openInventory{
			inventory = getInventory
		}else {
			inventory = specs.Inventory
		}
		specRow := map[string]interface{}{
			"id":        specs.Id,
			"key":       now,
			"name":      specs.Name,
			"code":      specs.Code,
			"price":     specs.Price,
			"market":specs.Market,
			"original":  specs.Original,
			"inventory": inventory,
			"limit":     specs.Limit,
			"max":specs.Max,
			"enable":    specs.Enable,
			"layer":     specs.Layer,
			"unit_id":      specs.UnitId,
			"image":business.GetDomainCosEncodePathName(global.GoodsPath,object.CId,specs.Image,false),
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
	goodsMap["imageUrl"] = "cos"
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
		//规格也进行更新
		e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id = ? and c_id = ?", row, userDto.CId).Updates(map[string]interface{}{
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
	var countAll int64
	e.Orm.Model(&models.Goods{}).Where("c_id = ?", userDto.CId).Count(&countAll)
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "goods", e.Orm)
	MaxNumber := CompanyCnf["goods"]
	if countAll >= int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("商品最多只可创建%v个", MaxNumber)), fmt.Sprintf("商品最多只可创建%v个", MaxNumber))
		return
	}

	var count int64
	e.Orm.Model(&models.Goods{}).Where("c_id = ? and name = ?", userDto.CId, req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	if req.Specs == ""{
		e.Error(500, errors.New("请配置规格"), "请配置规格")
		return
	}
	goodId,specDbMap, goodErr := s.Insert(userDto.CId, &req)
	if goodErr != nil {
		e.Error(500, err, fmt.Sprintf("创建商品失败,%s", goodErr.Error()))
		return
	}

	// 遍历所有图片
	fileForm, fileErr := c.MultipartForm()
	if fileErr != nil {
		e.Error(500, nil, "请提交表单模式")
		return
	}
	files := fileForm.File["files"]

	fileList := make([]string, 0)
	//txClient:=tx_api.TxCos{}
	buckClient :=qiniu.QinUi{CId: userDto.CId}
	buckClient.InitClient()
	//商品信息创建成功,才会保存客户的商品照片
	for _, file := range files {
		// 逐个存
		_,goodsImagePath  :=GetCosImagePath(global.GoodsPath,file.Filename,userDto.CId)
		if saveErr := c.SaveUploadedFile(file, goodsImagePath); saveErr == nil {

			//1.上传到cos中
			fileName,cosErr :=buckClient.PostImageFile(goodsImagePath)
			if cosErr !=nil{
				zap.S().Errorf("用户:%v,CID:%v 商品规格保存失败:%v",userDto.UserId,userDto.CId,cosErr)
				continue
			}
			//只保留文件名称,防止透露服务器地址
			fileList = append(fileList, fileName)
			//本地删除
			_=os.RemoveAll(goodsImagePath)
		}
		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?", goodId, userDto.CId).Updates(map[string]interface{}{
			"image": strings.Join(fileList, ","),
		})
	}
	//存储规格的图片
	//根据索引来创建
	specFiles := fileForm.File["spec_files"]
	//fmt.Println("规格DB",specDbMap)
	//fmt.Println("规格图片",specFiles)
	for index, file := range specFiles {
		fmt.Println("规格索引",index)
		specId,specOk:=specDbMap[index]
		if !specOk{
			continue
		}
		// 逐个存
		_,goodsImagePath  :=GetCosImagePath(global.GoodsPath,file.Filename,userDto.CId)
		if saveErr := c.SaveUploadedFile(file, goodsImagePath); saveErr == nil {

			//1.上传到cos中
			fileName,cosErr :=buckClient.PostImageFile(goodsImagePath)
			if cosErr !=nil{
				zap.S().Errorf("用户:%v,商品规格保存失败:%v",userDto.UserId,cosErr)
			}
			e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id = ? and c_id = ? and id = ?", goodId, userDto.CId,specId).Updates(map[string]interface{}{
				"image": fileName,
			})
			//本地删除
			_=os.Remove(goodsImagePath)
		}

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
	if req.Specs == ""{
		e.Error(500, errors.New("请配置规格"), "请配置规格")
		return
	}
	//设置桶
	buckClient :=qiniu.QinUi{CId: userDto.CId}
	buckClient.InitClient()


	CacheSpecImageMap,err := s.Update(userDto.CId,buckClient, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改商品信息失败,%s", err.Error()))
		return
	}
	var goodsObject models.Goods
	e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?",
		req.Id, userDto.CId).Limit(1).Find(&goodsObject)


	//原来的图片
	baseFileList := make([]string, 0)
	if goodsObject.Image != "" {
		baseFileList = strings.Split(goodsObject.Image, ",")
	}
	//初始化cos对象存储
	//txClient :=tx_api.TxCos{}


	fileForm, fileErr := c.MultipartForm()
	if fileErr != nil {
		e.Error(500, nil, "请提交表单模式")
		return
	}

	//商品的图片处理
	if req.FileClear == 1 {

		if goodsObject.Image != ""{
			for _, image := range strings.Split(goodsObject.Image,",") {
				buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId,global.GoodsPath,image))
			}
		}

		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?",
			req.Id, userDto.CId).Updates(map[string]interface{}{
			"image": "",
		})
	} else {
		//商品信息创建成功,才会保存客户的商品照片
		// 遍历所有图片
		fileList := make([]string, 0)
		files := fileForm.File["files"]
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

			buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId,global.GoodsPath,image))
		}
		for _, file := range files {
			// 逐个存
			//index

			_,goodsImagePath  :=GetCosImagePath(global.GoodsPath,file.Filename,userDto.CId)

			if saveErr := c.SaveUploadedFile(file, goodsImagePath); saveErr == nil {
				//只保留文件名称,防止透露服务器地址
				fileName,cosErr:=buckClient.PostImageFile(goodsImagePath)
				if cosErr !=nil{
					continue
				}
				fileList = append(fileList, fileName)
			}
			os.Remove(goodsImagePath)
		}
		e.Orm.Model(&models.Goods{}).Where("id = ? and c_id = ?", req.Id, userDto.CId).Updates(map[string]interface{}{
			"image": strings.Join(fileList, ","),
		})

	}
	//规格图片的处理
	if req.SpecFileClear == 1{

		//那就把规格的图片都清空掉
		specsList:=make([]models.GoodsSpecs,0)
		e.Orm.Model(&models.GoodsSpecs{}).Select("image").Where("goods_id = ? and c_id = ? ", req.Id, userDto.CId).Find(&specsList)

		for _,row:=range specsList{
			if row.Image != ""{
				buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId,global.GoodsPath,row.Image))
			}
			e.Orm.Model(&models.GoodsSpecs{}).Select("image").Where("goods_id = ? and c_id = ? ", req.Id, userDto.CId).Updates(map[string]interface{}{
				"image":"",
			})
		}

	}else {

		fmt.Println("规格和图片位置的map",CacheSpecImageMap)
		specFiles := fileForm.File["spec_files"]
		for index, file := range specFiles {
			// 逐个存
			_,goodsImagePath  :=GetCosImagePath(global.GoodsPath,file.Filename,userDto.CId)

			if saveErr := c.SaveUploadedFile(file, goodsImagePath); saveErr == nil {
				//只保留文件名称,防止透露服务器地址
				fileName,cosErr:=buckClient.PostImageFile(goodsImagePath)
				if cosErr !=nil{
					continue
				}

				fmt.Println("文件的索引",index)
				for specIdKey,v:=range CacheSpecImageMap {
					fmt.Println("规格图片的索引",specIdKey,v)
					if index == v {
						//因为图片更新了,那就把旧图删掉
						//保存一个新图
						GoodsImageSpecs:=models.GoodsSpecs{}
						e.Orm.Model(&models.GoodsSpecs{}).Select("id,image").Where("goods_id = ? and c_id = ? and id = ?", req.Id, userDto.CId,specIdKey).Limit(1).Find(&GoodsImageSpecs)
						if GoodsImageSpecs.Id > 0 && GoodsImageSpecs.Image !=""{
							buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId,global.GoodsPath,GoodsImageSpecs.Image))
						}

						e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id = ? and c_id = ? and id = ?", req.Id, userDto.CId,specIdKey).Updates(map[string]interface{}{
							"image": fileName,
						})
					}
				}
			}
			_=os.Remove(goodsImagePath)
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	notDelete,okDelete := s.Remove(&req,userDto.CId, p)

	res:=map[string]interface{}{
		"ok_delete":okDelete,
		"not_delete":notDelete,
	}
	if len(notDelete) > 0 {
		e.OK(business.Response{Code: 1,Data:res },"删除成功")
		return
	}
	e.OK(business.Response{Code: 0,Msg: "success"},"")
	return
}
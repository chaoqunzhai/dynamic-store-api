package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type CompanyInventory struct {
	api.Api
}

type GoodsFindRow struct {
	Name string
	Image string
}
//入库 如果商品不在库中 默认都是0个库存
//出库 如果商品不在库中 直接跳出
func (e CompanyInventory) Goods(c *gin.Context) {
	req := dto.InventoryGoodsReq{}
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

	whereSql :=fmt.Sprintf("c_id = %v and enable = %v",userDto.CId,true)

	goodsCnfMap:=make(map[int]GoodsFindRow,0)
	if req.Name != ""{
		likeVal:=fmt.Sprintf("%%%v%%",req.Name)
		var goodsLists []models2.Goods
		var goodsIds []string
		goodsSearchSql := fmt.Sprintf("%v and `name` like '%v'",whereSql,likeVal)
		e.Orm.Model(&models2.Goods{}).Select("id,name,image").Where(goodsSearchSql).Scopes(cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Find(&goodsLists)

		for _,row:=range goodsLists{
			goodsIds = append(goodsIds,fmt.Sprintf("%v",row.Id))
			goodsCnfMap[row.Id] =GoodsFindRow{
				Name: row.Name,
				Image: row.Image,
			}
		}
		goodsIds = utils.RemoveRepeatStr(goodsIds)

		if len(goodsIds) > 0 {
			whereSql =fmt.Sprintf("%v and goods_id in (%v)",whereSql,strings.Join(goodsIds,","))
		}
	}
	//fmt.Println("查询sql",whereSql)
	var goodsSpecs []models2.GoodsSpecs
	var count int64
	//根据分页获取商品
	e.Orm.Model(&models2.GoodsSpecs{}).Where(whereSql).Scopes(
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Find(&goodsSpecs).Limit(-1).Offset(-1).Count(&count)


	//统一在查一次商品
	//统一查询商品
	var goodsIds []int
	var inventoryKey []string
	for _,row:=range goodsSpecs{

		inventoryKey = append(inventoryKey,fmt.Sprintf("(goods_id = %v and spec_id = %v)",row.GoodsId,row.Id))
		//已经查询过了,那就不用查询了
		if _,ok:=goodsCnfMap[row.GoodsId];ok{
			continue
		}

		goodsIds = append(goodsIds,row.GoodsId)

	}
	inventoryKey = utils.RemoveRepeatStr(inventoryKey)
	//fmt.Println("inventoryKey",inventoryKey)
	//统一查库存
	InventoryDbMap:=make(map[string]models2.Inventory,0)
	if len(inventoryKey) > 0 {
		var InventoryObjectList []models2.Inventory
		e.Orm.Model(&models2.Inventory{}).Where(strings.Join(inventoryKey," or ")).Find(&InventoryObjectList)
		for _,row:=range InventoryObjectList{
			InventoryDbMap[fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId)] = row
		}
	}
	if len(goodsIds) > 0 {
		goodsIds = utils.RemoveRepeatInt(goodsIds)
		var goodsLists []models2.Goods
		e.Orm.Model(&models2.Goods{}).Select("id,name,image").Where("id in ?",goodsIds).Scopes(cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Find(&goodsLists)

		for _,row:=range goodsLists{
			//没有查到商品 在设置一次map
			goodsCnfMap[row.Id] = GoodsFindRow{
				Name: row.Name,
				Image: row.Image,
			}
		}
	}

	result :=make([]dto.GoodsSpecs,0)
	//组装一次数据 + 商品在库存中查询是否有
	for _,row:=range goodsSpecs{
		key :=fmt.Sprintf("%v_%v",row.GoodsId,row.Id)

		InventoryObject,cnfOk:= InventoryDbMap[key]

		goodsData,ok:=goodsCnfMap[row.GoodsId]
		if !ok{continue}
		tableRow :=dto.GoodsSpecs{
			Key: fmt.Sprintf("%v_%v",row.GoodsId,row.Id),
			Name: fmt.Sprintf("%v %v",goodsData.Name,row.Name),
			Unit: row.Unit,
			Image: func() string {
				if row.Image == "" {
					return business.GetGoodsPathFirst(row.CId,goodsData.Image,global.GoodsPath)
				}
				return business.GetGoodsPathFirst(row.CId,row.Image,global.GoodsPath)
			}(),
		}

		if cnfOk{
			//如果有库存 库存的值作为展示
			tableRow.Stock = InventoryObject.Stock
			tableRow.CostPrice = InventoryObject.OriginalPrice
			tableRow.Code = InventoryObject.Code
			tableRow.ArtNo = InventoryObject.ArtNo
		}else {
			//如果没有库存 默认就拿首次录入的库存 直接设置为0
			//tableRow.Stock = row.Inventory
			//tableRow.Price = float64(row.Original)
		}

		result = append(result,tableRow)
	}


	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
}



func (e CompanyInventory) ManageEdit(c *gin.Context) {
	req := dto.EditReq{}
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
	var object models2.Inventory

	e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Limit(1).Find(&object)

	if object.Id == 0 {
		e.Error(500, nil,"数据不存在")
		return
	}
	e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and id = ?",userDto.CId,req.Id).Updates(map[string]interface{}{
		"art_no":req.EditArtNo,
		"code":req.EditCode,
		"original_price":req.EditOriginalPrice,
	})
	e.OK("","保存成功")
	return
}
func (e CompanyInventory) ManageGetPage(c *gin.Context) {
	req := dto.ManageListGetPageReq{}
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
	whereSql :=fmt.Sprintf("c_id = %v ",userDto.CId)
	//进行商品名称的查询
	if req.Name != "" {
		likeVal:=fmt.Sprintf("%%%v%%",req.Name)
		var goodsLists []models2.Goods
		var goodsIds []string
		goodsSearchSql := fmt.Sprintf("c_id = %v and `name` like '%v'",userDto.CId,likeVal)
		e.Orm.Model(&models2.Goods{}).Select("id,name,image").Where(goodsSearchSql).Scopes(cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Find(&goodsLists)

		for _,row:=range goodsLists{
			goodsIds = append(goodsIds,fmt.Sprintf("%v",row.Id))
		}
		goodsIds = utils.RemoveRepeatStr(goodsIds)

		if len(goodsIds) > 0 {
			whereSql =fmt.Sprintf("%v and goods_id in (%v)",whereSql,strings.Join(goodsIds,","))
		}
	}

	if req.Action == "out" { //如果是出库 那只查看大于0的数据
		whereSql += " and stock > 0 "
	}
	var count int64
	result:=make([]interface{},0)
	InventoryList :=make([]models2.Inventory,0)
	e.Orm.Model(&models2.Inventory{}).Where(whereSql).Scopes(
		cDto.MakeCondition(req.GetNeedSearch()),
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Order(global.OrderLayerKey).Find(&InventoryList).Limit(-1).Offset(-1).Count(&count)

	findGoodsSpecKey:=make([]string,0) //规格列表
	findGoodsKey :=make([]int,0) //商品列表
	for _,row:=range InventoryList{
		findGoodsSpecKey = append(findGoodsSpecKey,fmt.Sprintf(" (goods_id = %v and id = %v) ",row.GoodsId,row.SpecId))
		findGoodsKey = append(findGoodsKey,row.GoodsId)
	}
	findGoodsSpecKey = utils.RemoveRepeatStr(findGoodsSpecKey)
	findGoodsKey = utils.RemoveRepeatInt(findGoodsKey)
	GoodsSpecsInfoMap:=make(map[string]dto.GoodsInfo,0) //商品规格: 信息
	GoodsInfoMap:=make(map[int]dto.GoodsInfo,0) //商品: 信息
	if len(findGoodsSpecKey) > 0 {
		var GoodsObjectList []models2.Goods
		e.Orm.Model(&models2.Goods{}).Select("id,name,image").Where("id in ?",findGoodsKey).Find(&GoodsObjectList)
		for _,row:=range GoodsObjectList{
			GoodsInfoMap[row.Id] = dto.GoodsInfo{
				Name: row.Name,
				Image: row.Image,
			}
		}

		var GoodsSpecsObjectList []models2.GoodsSpecs
		e.Orm.Model(&models2.GoodsSpecs{}).Select("id,unit,name,image,goods_id").Where(strings.Join(findGoodsSpecKey," or ")).Find(&GoodsSpecsObjectList)
		for _,row:=range GoodsSpecsObjectList{
			GoodsSpecsInfoMap[fmt.Sprintf("%v_%v",row.GoodsId,row.Id)] = dto.GoodsInfo{
				Name: row.Name,
				Unit: row.Unit,
				Image: row.Image,
			}
		}
	}

	for _,row:=range InventoryList{
		key:=fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId)

		goodsData:=GoodsInfoMap[row.GoodsId]
		goodsSpecs:=GoodsSpecsInfoMap[key]

		imageVal := goodsSpecs.Image
		if goodsSpecs.Image == ""{
			//商品如果有图片,那获取第一张图片即可
			if goodsData.Image != ""{
				imageVal = strings.Split( goodsData.Image,",")[0]
			}else {
				imageVal = ""
			}
		}

		data:=map[string]interface{}{
			"key":fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId),
			"name":fmt.Sprintf("%v %v",goodsData.Name,goodsSpecs.Name),
			"image":func() string {
				if imageVal == "" {
					return ""
				}
				return business.GetGoodsPathFirst(row.CId,imageVal,global.GoodsPath)
			}(),
			"unit":goodsSpecs.Unit,
			"cost_price":row.OriginalPrice,
			"original_price":utils.StringDecimal(row.OriginalPrice),
			"time":row.UpdatedAt.Format("2006-01-02 15:04:05"),
			"stock":row.Stock,
			"code":row.Code,
			"art_no":row.ArtNo,
			"id":row.Id,
		}
		result = append(result,data)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return


}

func (e CompanyInventory) ManageRecords(c *gin.Context) {
	req := dto.RecordsListGetPageReq{}
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
	skuId:=c.Param("skuId")
	if skuId == ""{
		e.Error(500, nil,"请选择库存条目")
		return
	}
	var  Inventory models2.Inventory
	e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and id = ?",userDto.CId,skuId).Limit(1).Find(&Inventory)
	if Inventory.Id == 0 {
		e.Error(500, nil,"库存不存在")
		return
	}
	var RecordsList []models2.InventoryRecord
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	var count int64
	result:=make([]interface{},0)
	e.Orm.Table(splitTableRes.InventoryRecordLog).Where(
		"c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,Inventory.GoodsId,Inventory.SpecId).Scopes(
				cDto.MakeSplitTableCondition(req.GetNeedSearch(),splitTableRes.InventoryRecordLog),
				cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			).Order("id desc").Find(&RecordsList).Limit(-1).Offset(-1).Count(&count)


	for _,row:=range RecordsList{
		data :=map[string]interface{}{
			"id":row.Id,
			"create_at":row.CreatedAt.Format("2006-01-02 15:04:05"),
			"user":row.CreateBy,
			"source_number":row.SourceNumber,
			"current_number":row.CurrentNumber,
			"action":row.Action,
			"original_price":utils.StringDecimal(row.OriginalPrice),
			"source_price":utils.StringDecimal(row.SourcePrice),
		}
		switch row.Action {
		case global.InventoryIn:
			data["action_cn"] = "入库"
			data["action_number"] = fmt.Sprintf("+%v",row.ActionNumber)
		case global.InventoryOut:
			data["action_number"] = fmt.Sprintf("-%v",row.ActionNumber)
			data["action_cn"] = "出库"
		case global.InventoryRefundIn:
			data["action_number"] = fmt.Sprintf("+%v",row.ActionNumber)
			data["action_cn"] = "退货入库"
		}
		result = append(result,data)
	}


	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
}

func (e CompanyInventory) RecordsLog(c *gin.Context) {
	req := dto.RecordsListGetPageReq{}
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


	var RecordsList []models2.InventoryRecord
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	var count int64
	result:=make([]interface{},0)
	e.Orm.Table(splitTableRes.InventoryRecordLog).Where("c_id = ? ",userDto.CId).Scopes(
		cDto.MakeSplitTableCondition(req.GetNeedSearch(),splitTableRes.InventoryRecordLog),
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Order("id desc").Find(&RecordsList).Limit(-1).Offset(-1).Count(&count)


	for _,row:=range RecordsList{
		data :=map[string]interface{}{
			"id":row.Id,
			"create_at":row.CreatedAt.Format("2006-01-02 15:04:05"),
			"user":row.CreateBy,
			"goods_name":fmt.Sprintf("%v %v",row.GoodsName,row.GoodsSpecName),
			"source_number":row.SourceNumber,
			"current_number":row.CurrentNumber,
			"original_price":utils.StringDecimal(row.OriginalPrice),
			"source_price":utils.StringDecimal(row.SourcePrice),
		}
		switch row.Action {
		case global.InventoryIn:
			data["action_cn"] = "入库"
			data["action_number"] = fmt.Sprintf("+%v",row.ActionNumber)
		case global.InventoryOut:
			data["action_number"] = fmt.Sprintf("-%v",row.ActionNumber)
			data["action_cn"] = "出库"
		case global.InventoryRefundIn:
			data["action_number"] = fmt.Sprintf("+%v",row.ActionNumber)
			data["action_cn"] = "退货入库"
		}
		result = append(result,data)
	}


	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
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

func (e CompanyInventory) OrderList(c *gin.Context) {
	req := dto.OrderListGetPageReq{}
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
	var count int64
	result :=make([]interface{},0)

	whereSql :=fmt.Sprintf("c_id = %v ",userDto.CId)
	if req.OrderId != ""{

		likeVal:=fmt.Sprintf("%%%v%%",req.OrderId)

		whereSql = fmt.Sprintf("%v and `order_id` like '%v'",whereSql,likeVal)
	}
	switch req.Action {
	case "in":
		whereSql += fmt.Sprintf(fmt.Sprintf(" and action = %v",1))
	case "out":

		whereSql += fmt.Sprintf(fmt.Sprintf(" and action = %v",2))
	}
	var list []models2.InventoryOrder
	e.Orm.Model(&models2.InventoryOrder{}).Where(whereSql).Scopes(
		cDto.MakeCondition(req.GetNeedSearch()),
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Order("id desc").Find(&list).Limit(-1).Offset(-1).Count(&count)

	for _,row:=range list{
		result = append(result,map[string]interface{}{
			"order_id":row.OrderId,
			"create_time":row.CreatedAt.Format("2006-01-02 15:04:05"),
			"user":row.CreateBy,
			"desc": row.Desc,
			"money":utils.StringDecimal(row.DocumentMoney),
			"number":row.Number,
		})
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
}

func (e CompanyInventory) WarehousingCreate(c *gin.Context) {
	req := dto.InventoryCreateReq{}
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

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	OrderId:=fmt.Sprintf("%v",utils.GenUUID())
	//创建入库单的的入库流水
	DocumentMoney:=0.00
	Number :=0
	//循环数据开始入库
	for rowKey:=range req.Data{

		data,ok :=req.Data[rowKey]
		if !ok{continue}
		//负数不处理
		if data.ActionNumber < 0 || data.CostPrice < 0  {
			continue
		}
		rowKeyList :=strings.Split(rowKey,"_")
		if len(rowKeyList) != 2{
			continue
		}
		goodsId,stroveErr :=strconv.Atoi(rowKeyList[0])
		if stroveErr!=nil{
			continue
		}
		specsId,stroveErr :=strconv.Atoi(rowKeyList[1])
		if stroveErr!=nil{
			continue
		}
		var InventoryObject models2.Inventory
		e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,goodsId,specsId).Limit(1).Find(&InventoryObject)
		var SourceNumber int //原来数量
		var OriginalPrice float64 //原来入库价
		var goodsObject models2.Goods
		var goodsSpecs models2.GoodsSpecs
		e.Orm.Model(&models2.Goods{}).Where("c_id = ? and id = ? ",userDto.CId,goodsId).Limit(1).Find(&goodsObject)
		if goodsObject.Id == 0 {
			continue
		}

		e.Orm.Model(&models2.GoodsSpecs{}).Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,goodsId,specsId).Limit(1).Find(&goodsSpecs)
		if goodsSpecs.Id == 0 {
			continue
		}

		//规格没有录图片的时候 拿商品的图片
		imageVal := goodsSpecs.Image
		if goodsSpecs.Image == ""{
			//商品如果有图片,那获取第一张图片即可
			if goodsObject.Image != ""{
				imageVal = strings.Split( goodsObject.Image,",")[0]
			}else {
				imageVal = ""
			}
		}

		if InventoryObject.Id == 0 {
			//创建 在查询一次
			//需要使用 规格中的商品总数 方便兼容 数据的融合
			InventoryObject = models2.Inventory{
				Stock: data.ActionNumber, // 废弃goodsSpecs.Inventory + data.ActionNumber 只以前段录入数据为主
				OriginalPrice: data.CostPrice,
				GoodsId: goodsId,
				SpecId: specsId,
				Code: data.Code,
				ArtNo: data.ArtNo,
			}
			InventoryObject.CId = userDto.CId
			InventoryObject.CreateBy = userDto.UserId
			if createErr:=e.Orm.Create(&InventoryObject).Error;createErr!=nil{
				zap.S().Errorf("客户 %v 仓库入库创建数据失败,数据:%v 原因:%v",userDto.UserId,data,createErr.Error())
				continue
			}
			//新数据,那原库存就是0
			SourceNumber = 0
			//使用规格的入库价
			OriginalPrice = float64(goodsSpecs.Original)
		}else {
			//原库存
			SourceNumber = InventoryObject.Stock
			OriginalPrice = InventoryObject.OriginalPrice

			InventoryObject.ArtNo = data.ArtNo
			InventoryObject.Code = data.Code
			InventoryObject.Stock += data.ActionNumber
			InventoryObject.OriginalPrice = data.CostPrice
			if saveErr:=e.Orm.Save(&InventoryObject).Error;saveErr!=nil{
				zap.S().Errorf("客户 %v 仓库入库保存数据失败,数据:%v 原因:%v",userDto.UserId,data,saveErr.Error())
				continue
			}
		}

		//流水创建
		RecordLog:=models2.InventoryRecord{
			CId: userDto.CId,
			CreateBy:userDto.Username,
			OrderId: OrderId,
			Action: global.InventoryIn, //入库
			Image: imageVal,
			GoodsId: InventoryObject.GoodsId,
			GoodsName: goodsObject.Name,
			GoodsSpecName: goodsSpecs.Name,
			SpecId: InventoryObject.SpecId,
			SourceNumber:SourceNumber, //原库存
			ActionNumber:data.ActionNumber, //操作的库存
			CurrentNumber:SourceNumber + data.ActionNumber, //那现库存 就是 原库存 + 操作的库存
			OriginalPrice:data.CostPrice,
			SourcePrice:OriginalPrice, //原入库价
			Unit:data.Unit,
		}
		e.Orm.Table(splitTableRes.InventoryRecordLog).Create(&RecordLog)


		//创建成功 金额叠加
		DocumentMoney += utils.RoundDecimalFlot64(float64(data.ActionNumber) * data.CostPrice)
		Number +=data.ActionNumber

	}
	//创建一条入库单记录
	object :=models2.InventoryOrder{
		OrderId: OrderId,
		Action: global.InventoryIn,
		DocumentMoney:DocumentMoney,
		Number: Number,
	}
	object.Desc = req.Desc
	object.CId = userDto.CId
	object.CreateBy = userDto.Username
	e.Orm.Create(&object)

	e.OK("","入库成功")
	return
}



func (e CompanyInventory) OrderDetail(c *gin.Context) {
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
	orderId:=c.Param("orderId")
	var object models2.InventoryOrder
	e.Orm.Model(&models2.InventoryOrder{}).Where("c_id = ? and order_id = ?",userDto.CId,orderId).Limit(1).Find(&object)

	if object.Id == 0 {
		e.Error(500,nil,"数据不存在")
		return
	}
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	result:=make(map[string]interface{},0)
	result["record"] = object
	RecordLog:=make([]models2.InventoryRecord,0)
	e.Orm.Table(splitTableRes.InventoryRecordLog).Where("c_id = ? and order_id = ?",userDto.CId,orderId).Find(&RecordLog)

	table:=make([]interface{},0)
	var number int
	var classCount int
	var money float64
	for _,row:=range RecordLog{
		thisPrice := utils.RoundDecimalFlot64(float64(row.ActionNumber) * row.OriginalPrice)
		classCount +=1
		number +=row.ActionNumber
		money += thisPrice
		table = append(table,map[string]interface{}{
			"image":func() string {
				if row.Image == "" {
				return ""
			}
				return business.GetGoodsPathFirst(row.CId,row.Image,global.GoodsPath)
			}(),
			"id":row.Id,
			"goods_name":row.GoodsName,
			"goods_spec_name":row.GoodsSpecName,
			"unit":row.Unit,
			"action_number":row.ActionNumber,
			"original_price":utils.StringDecimal(row.OriginalPrice),
			"money":utils.StringDecimal(thisPrice),
		})
	}
	result["CountAll"] = map[string]interface{}{
		"number":number,
		"class":classCount,
		"money":utils.StringDecimal(money),
	}
	result["table"] = table
	e.OK(result,"successful")
	return
}

func (e CompanyInventory) OutboundCreate(c *gin.Context) {
	req := dto.InventoryCreateReq{}
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

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	OrderId:=fmt.Sprintf("%v",utils.GenUUID())
	//创建出库单的的入库流水
	DocumentMoney:=0.00
	Number :=0
	//循环查一次数据

	//循环数据开始入库
	for rowKey:=range req.Data{

		data,ok :=req.Data[rowKey]
		if !ok{continue}
		//负数不处理
		if data.ActionNumber < 0 || data.CostPrice < 0  {
			continue
		}
		rowKeyList :=strings.Split(rowKey,"_")
		if len(rowKeyList) != 2{
			continue
		}
		goodsId, stroveErr :=strconv.Atoi(rowKeyList[0])
		if stroveErr !=nil{
			continue
		}
		specsId, stroveErr :=strconv.Atoi(rowKeyList[1])
		if stroveErr !=nil{
			continue
		}
		var InventoryObject models2.Inventory
		e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,goodsId,specsId).Limit(1).Find(&InventoryObject)

		//必须在库存中 才可以进行出库
		if InventoryObject.Id == 0 {

			continue
		}
		var goodsObject models2.Goods
		var goodsSpecs models2.GoodsSpecs
		e.Orm.Model(&models2.Goods{}).Where("c_id = ? and id = ? ",userDto.CId,goodsId).Limit(1).Find(&goodsObject)
		if goodsObject.Id == 0 {
			continue
		}

		e.Orm.Model(&models2.GoodsSpecs{}).Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,goodsId,specsId).Limit(1).Find(&goodsSpecs)
		if goodsSpecs.Id == 0 {
			continue
		}
		//规格没有录图片的时候 拿商品的图片
		imageVal := goodsSpecs.Image
		if goodsSpecs.Image == ""{
			//商品如果有图片,那获取第一张图片即可
			if goodsObject.Image != ""{
				imageVal = strings.Split( goodsObject.Image,",")[0]
			}else {
				imageVal = ""
			}
		}
		//原库存
		SourceNumber := InventoryObject.Stock
		//进行减少覆盖即可,但是不能超用
		if data.ActionNumber > InventoryObject.Stock{
			continue
		}
		InventoryObject.Stock -= data.ActionNumber
		if saveErr:=e.Orm.Save(&InventoryObject).Error;saveErr!=nil{
			zap.S().Errorf("客户 %v 仓库出库保存数据失败,数据:%v 原因:%v",userDto.UserId,data,saveErr.Error())
			continue
		}

		//流水创建
		RecordLog:=models2.InventoryRecord{
			CId: userDto.CId,
			CreateBy:userDto.Username,
			OrderId: OrderId,
			Action: global.InventoryOut, //入库
			Image: imageVal,
			GoodsId: InventoryObject.GoodsId,
			GoodsName: goodsObject.Name,
			GoodsSpecName: goodsSpecs.Name,
			SpecId: InventoryObject.SpecId,
			SourceNumber:SourceNumber, //原库存
			ActionNumber:data.ActionNumber, //操作的库存
			CurrentNumber:SourceNumber - data.ActionNumber, //那现库存 就是 原库存 - 操作的库存
			OriginalPrice:data.CostPrice,
			Unit:data.Unit,
		}
		e.Orm.Table(splitTableRes.InventoryRecordLog).Create(&RecordLog)


		//创建成功 金额叠加
		DocumentMoney += utils.RoundDecimalFlot64(float64(data.ActionNumber) * data.CostPrice)
		Number +=data.ActionNumber

	}
	//创建一条入库单记录
	object :=models2.InventoryOrder{
		OrderId: OrderId,
		Action: global.InventoryOut,
		DocumentMoney:DocumentMoney,
		Number: Number,
	}
	object.Desc = req.Desc
	object.CId = userDto.CId
	object.CreateBy = userDto.Username
	e.Orm.Create(&object)

	e.OK("","出库成功")
	return
}


func (e CompanyInventory) OutboundDetail(c *gin.Context) {

}
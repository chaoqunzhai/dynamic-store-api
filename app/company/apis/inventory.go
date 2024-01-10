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

	goodsCnfMap:=make(map[int]string,0)
	if req.Name != ""{
		likeVal:=fmt.Sprintf("%%%v%%",req.Name)
		var goodsLists []models2.Goods
		var goodsIds []string
		goodsSearchSql := fmt.Sprintf("%v and `name` like '%v'",whereSql,likeVal)
		e.Orm.Model(&models2.Goods{}).Select("id,name").Where(goodsSearchSql).Scopes(cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Find(&goodsLists)

		for _,row:=range goodsLists{
			goodsIds = append(goodsIds,fmt.Sprintf("%v",row.Id))
			goodsCnfMap[row.Id] = row.Name
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
	).Find(&goodsSpecs).Count(&count)


	//统一在查一次商品
	//统一查询商品
	var goodsIds []int
	var inventoryKey []string
	for _,row:=range goodsSpecs{
		//已经查询过了,那就不用查询了
		if _,ok:=goodsCnfMap[row.GoodsId];ok{
			continue
		}

		goodsIds = append(goodsIds,row.GoodsId)

		inventoryKey = append(inventoryKey,fmt.Sprintf("(goods_id = %v and spec_id = %v)",row.GoodsId,row.Id))

	}
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
		e.Orm.Model(&models2.Goods{}).Select("id,name").Where("id in ?",goodsIds).Scopes(cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Find(&goodsLists)

		for _,row:=range goodsLists{
			//没有查到商品 在设置一次map
			goodsCnfMap[row.Id] = row.Name
		}
	}

	result :=make([]dto.GoodsSpecs,0)
	//组装一次数据 + 商品在库存中查询是否有
	for _,row:=range goodsSpecs{
		goodsName,ok:=goodsCnfMap[row.GoodsId]
		if !ok{continue}
		tableRow :=dto.GoodsSpecs{
			Key: fmt.Sprintf("%v_%v",row.GoodsId,row.Id),
			Name: fmt.Sprintf("%v %v",goodsName,row.Name),
			Unit: row.Unit,
			Image: func() string {
				if row.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(row.CId,row.Image,global.GoodsPath)
			}(),
		}

		key :=fmt.Sprintf("%v_%v",row.GoodsId,row.Id)

		InventoryObject,cnfOk:= InventoryDbMap[key]
		if cnfOk{
			//如果有库存 库存的值作为展示
			tableRow.Stock = InventoryObject.Stock
			tableRow.Price = InventoryObject.OriginalPrice
			tableRow.Image = func() string {
				if InventoryObject.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(row.CId,InventoryObject.Image,global.GoodsPath)
			}()
		}else {
			//如果没有库存 默认就拿首次录入的库存
			tableRow.Stock = row.Inventory
			tableRow.Price = float64(row.Original)
		}

		result = append(result,tableRow)
	}


	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
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
func (e CompanyInventory) InRecords(c *gin.Context) {

}
func (e CompanyInventory) Warehousing(c *gin.Context) {
	req := dto.WarehousingGetPageReq{}
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

	whereSql :=fmt.Sprintf("c_id = %v and action = %v",userDto.CId,1)
	if req.OrderId != ""{

		likeVal:=fmt.Sprintf("%%%v%%",req.OrderId)

		whereSql = fmt.Sprintf("%v and `order_id` like '%v'",whereSql,likeVal)
	}
	var list []models2.InventoryOrder
	e.Orm.Model(&models2.InventoryOrder{}).Where(whereSql).Scopes(cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Order("id desc").Find(&list).Count(&count)

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
	req := dto.WarehousingCreateReq{}
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
		goodsId,err :=strconv.Atoi(rowKeyList[0])
		if err!=nil{
			continue
		}
		specsId,err :=strconv.Atoi(rowKeyList[1])
		if err!=nil{
			continue
		}
		var InventoryObject models2.Inventory
		e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,goodsId,specsId).Limit(1).Find(&InventoryObject)
		var SourceNumber int
		if InventoryObject.Id == 0 {
			//创建 在查询一次
			var goods models2.Goods
			var goodsSpecs models2.GoodsSpecs
			e.Orm.Model(&models2.Goods{}).Where("c_id = ? and id = ? ",userDto.CId,goodsId).Limit(1).Find(&goods)

			e.Orm.Model(&models2.GoodsSpecs{}).Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,goodsId,specsId).Limit(1).Find(&goodsSpecs)

			//规格没有录图片的时候 拿商品的图片
			imageVal := goodsSpecs.Image
			if goodsSpecs.Image == ""{
				imageVal = goods.Image
			}
			//需要使用 规格中的商品总数 方便兼容 数据的融合
			InventoryObject = models2.Inventory{
				Stock: goodsSpecs.Inventory + data.ActionNumber,
				OriginalPrice: data.CostPrice,
				GoodsId: goodsId,
				GoodsName: goods.Name,
				GoodsSpecName: goodsSpecs.Name,
				SpecId: specsId,
				Image: imageVal,
			}

			InventoryObject.CId = userDto.CId
			InventoryObject.CreateBy = userDto.UserId
			if createErr:=e.Orm.Create(&InventoryObject).Error;createErr!=nil{
				zap.S().Errorf("客户 %v 仓库录入创建数据失败,数据:%v 原因:%v",userDto.UserId,data,createErr.Error())
				continue
			}
			//新数据,那原库存就是0
			SourceNumber = 0
		}else {
			//原库存
			SourceNumber = InventoryObject.Stock
			//增加覆盖即可
			InventoryObject.Stock += data.ActionNumber
			InventoryObject.OriginalPrice = data.CostPrice
			if saveErr:=e.Orm.Save(&InventoryObject).Error;saveErr!=nil{
				zap.S().Errorf("客户 %v 仓库录入数据失败,数据:%v 原因:%v",userDto.UserId,data,saveErr.Error())
				continue
			}
		}

		//流水创建
		RecordLog:=models2.InventoryRecord{
			CId: userDto.CId,
			CreateBy:userDto.Username,
			OrderId: OrderId,
			Action: 1, //入库
			Image: InventoryObject.Image,
			GoodsId: InventoryObject.GoodsId,
			GoodsName: InventoryObject.GoodsName,
			GoodsSpecName: InventoryObject.GoodsSpecName,
			SpecId: InventoryObject.SpecId,
			SourceNumber:SourceNumber, //原库存
			ActionNumber:data.ActionNumber, //操作的库存
			CurrentNumber:SourceNumber + data.ActionNumber, //那现库存 就是 原库存 + 操作的库存
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
		Action: 1,
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



func (e CompanyInventory) WarehousingDetail(c *gin.Context) {

}



func (e CompanyInventory) Outbound(c *gin.Context) {

}

func (e CompanyInventory) OutboundCreate(c *gin.Context) {

}



func (e CompanyInventory) OutboundDetail(c *gin.Context) {

}
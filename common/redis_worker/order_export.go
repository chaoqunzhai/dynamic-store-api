/**
@Author: chaoqun
* @Date: 2023/12/25 18:56
*/
package redis_worker

import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/redis_db"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type OrderExportObj struct {
	Orm *gorm.DB
	Dat ExportReq
	RedisKey string
	FileName string
}
type SheetRow struct {
	OrderA2 string //索引标记而已
	SheetName string //分页名称 小B名称
	OrderCreateTime string //订单创建时间
	ShopAddress string //小B地址
	ShopPhone string //联系电话
	ShopUserValue string //联系人信息
	ExportTime string `json:"export_time"` //导出操作时间
	ExportUser string //操作人
	DeliveryMoney string //运费
	Table []*XlsxTableRow `json:"table"`

}
type XlsxTableRow struct {

}

//开始执行数据导出
//多个订单save为一个excel文件
//多个订单 同一个小B放到一个sheet中

func (e *OrderExportObj)ReadOrderDetail() error  {

	fmt.Println("开始解析",e.Dat)
	//多个订单 同一个小B的数据放在一起

	siteMap:=make(map[int][]*SheetRow,0)
	//
	nowTimeObj :=time.Now()
	var orderList []models.Orders

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)
	e.Orm.Table(splitTableRes.OrderTable).Where("order_id in ?",e.Dat.Order).Find(&orderList)

	for index,orderRow:=range orderList{
		var shopRow models.Shop
		e.Orm.Model(&models.Shop{}).Where("id = ? ", orderRow.ShopId).Limit(1).Find(&shopRow)
		if shopRow.Id == 0 {
			continue
		}
		sheetRow:=&SheetRow{
			OrderA2: fmt.Sprintf("DCY.%v.%v",nowTimeObj.Format("20060102"),index + 1),
			ExportTime: nowTimeObj.Format("2006-01-02 15:04:05"),
			SheetName: shopRow.Name,
			ShopPhone: shopRow.Phone,
			ShopAddress: shopRow.Address,
			ShopUserValue: shopRow.UserName,
			OrderCreateTime: orderRow.CreatedAt.Format("2006-01-02 15:04"),
		}
		data,ok:=siteMap[orderRow.ShopId]
		if ok{
			data = append(data, sheetRow)
		}else {
			data = make([]*SheetRow,0)
			data = append(data, sheetRow)

		}
		siteMap[orderRow.ShopId] = data

	}

	return nil
}


//导出数据保存在云端

func (e *OrderExportObj)SaveExportZIP() error {
	//fmt.Println("保存到本地zip文件",e.Dat)
	return nil
}

//更新table中状态

func (e *OrderExportObj)SaveExportDb(successTag bool,msg string)  error{
	var status int
	if !successTag{
		status = 2
	}else {
		status = 1
	}
	if msg !=""{
		if len(msg) > 60{
			msg = msg[:60]
		}
	}
	e.Orm.Model(&models.CompanyTasks{}).Where("id = ? and c_id = ?",
		e.Dat.OrmId,e.Dat.CId).Updates(map[string]interface{}{
		"status":status,
		"path":e.FileName,
		"msg":msg,
	})
	return nil

}
//如果key下的list位空 ,那就支持清空这个key

func  (e *OrderExportObj)EmptyKey(keyLen int) {
	err :=redis_db.RedisCli.LTrim(RedisCtx,e.RedisKey,1,int64(keyLen)).Err()
	if err!=nil{
		zap.S().Errorf("清理redis key:%v 数据清理失败:%v",e.RedisKey,err)
	}else {
		zap.S().Infof("redis key:%v 消费完毕,数据清理成功",e.RedisKey)
	}
}
package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	models2 "go-admin/common/models"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"
	"time"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)
var (
	Loc, _ = time.LoadLocation("Local")
)
type CycleTimeConf struct {
	service.Service
}

type ValidCycle struct {
	IsTime bool `json:"is_time"` //是否在这个时间区间
	DeliveryMsg string `json:"delivery_msg"` //送达时间
	DeliveryId int `json:"delivery_id"` //配送ID
	StartTime time.Time `json:"-" gorm:"comment:记录可下单周期开始时间"`
	EndTime   time.Time `json:"-" gorm:"comment:记录可下单周期结束时间"`
	//下单周期的文案也是保持最新的
	CreateStr string `json:"-" gorm:"size:30;comment:下单日期的文案内容"`
	//配送周期的统一查询
	DeliveryTime  models2.XTime   `json:"-" gorm:"type:date;comment:计算的配送时间"`
	//展示,也是保持最新的
	DeliveryStr string `json:"-" gorm:"size:30;comment:配送文案"`

}

// GetPage 获取CycleTimeConf列表
func (e *CycleTimeConf) GetPage(c *dto.CycleTimeConfGetPageReq, p *actions.DataPermission, list *[]models.CycleTimeConf, count *int64) error {
	var err error
	var data models.CycleTimeConf

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("created_at desc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CycleTimeConfService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CycleTimeConf对象
func (e *CycleTimeConf) Get(d *dto.CycleTimeConfGetReq, p *actions.DataPermission, model *models.CycleTimeConf) error {
	var data models.CycleTimeConf

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCycleTimeConf error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CycleTimeConf对象
func (e *CycleTimeConf) Insert(cid int, c *dto.CycleTimeConfInsertReq) error {
	var err error
	var data models.CycleTimeConf
	c.Generate(&data)
	data.CId = cid
	//data.Enable = true
	data.Uid = utils.CreateCode()
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CycleTimeConfService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CycleTimeConf对象
func (e *CycleTimeConf) Update(c *dto.CycleTimeConfUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CycleTimeConf{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	//判断时间是否有变动,如果有变动更新uid标记
	uidTag := false
	if c.StartTime != data.StartTime || c.EndTime != data.EndTime {
		uidTag = true
	}
	switch c.Type {
	case global.CyCleTimeWeek:
		if c.StartWeek != data.StartWeek || c.EndWeek != data.EndWeek {
			uidTag = true
		}
	}
	if uidTag {
		data.Uid = utils.CreateCode()
		fmt.Println("数据发生变更更新code")
	}
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CycleTimeConfService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CycleTimeConf
func (e *CycleTimeConf) Remove(d *dto.CycleTimeConfDeleteReq, p *actions.DataPermission) error {
	var data models.CycleTimeConf

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCycleTimeConf error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
func GetOrderCyClyCnf(CyCleCnf models.CycleTimeConf) (day models2.XTime, dayValue string) {

	//时间检测返回以下内容
	//配送的具体日期Time对象
	//配送的具体时间文案做展示
	nowDay:=models2.XTime{
		Time:time.Now(),
	}

	if CyCleCnf.GiveDay == 0 {
		if CyCleCnf.GiveTime != ""{
			return nowDay,fmt.Sprintf("%v", CyCleCnf.GiveTime)
		}
		return nowDay,"当天配送"
	}

	cycleTimeValue :=CalculateTime(CyCleCnf.GiveDay)

	cycleVal := fmt.Sprintf("%v %v",cycleTimeValue.Format("2006-01-02"), CyCleCnf.GiveTime)

	return cycleTimeValue,cycleVal

}
func GetOrderCreateStr(row models.CycleTimeConf) string {
	orderCreateStr:=""
	switch row.Type {
	case global.CyCleTimeDay:
		orderCreateStr = fmt.Sprintf("每天 %v-%v", row.StartTime, row.EndTime)
	case global.CyCleTimeWeek:
		orderCreateStr = fmt.Sprintf("每周%v %v-每周%v %v", global.WeekIntToMsg(row.StartWeek), row.StartTime,
			global.WeekIntToMsg(row.EndWeek), row.EndTime,
		)
	}

	return orderCreateStr

}

//只有支付成功的时候才会调用这个方法
//1.记录当前时间的哪个下单的时间段,获取到这个下单开始和结束时间 和文案 记录下来
//2.计算这个配送时间也需要录入DB中

//只需要返回这个上层的UID
func CheckOrderCyCleCnfIsDb(cid int,table string,DeliveryObject models.CycleTimeConf,orm *gorm.DB) (uid,deliveryMsg string)  {

	//计算出配送的周期
	deliveryTime,deliveryMsg:=GetOrderCyClyCnf(DeliveryObject)

	//检测下,然后直接返回吧
	var cycleCnf models.OrderCycleCnf

	//查询大B + 配送时间为同一天
	//如果没有,那就创建
	//如果有,那就统一订单的这个UID
	//查询天！！
	orm.Table(table).Model(&models.OrderCycleCnf{}).Where("c_id = ? and delivery_time = ?",cid,deliveryTime.Format("2006-01-02")).Find(&cycleCnf)
	if cycleCnf.Id == 0 {

		//生成一个新的uid,让订单来记录
		uid = RandomCode()

		//把生成的配送信息录入到DB中,做一个订单统筹

		nowTime:=time.Now()
		noYearMonth:=nowTime.Format("2006-01-02")
		if DeliveryObject.StartTime != "" && DeliveryObject.EndTime != ""{

			startStr:=fmt.Sprintf("%v %v",noYearMonth,DeliveryObject.StartTime)
			endStr:=fmt.Sprintf("%v %v",noYearMonth,DeliveryObject.EndTime)
			sTime, _:=time.ParseInLocation("2006-01-02 15:04", startStr, Loc)

			eTime,_:=time.ParseInLocation("2006-01-02 15:04", endStr, Loc)

			//获取到配送的文案和天数
			CreateStr:=GetOrderCreateStr(DeliveryObject)

			orm.Table(table).Create(&models.OrderCycleCnf{
				CId:cid,
				Uid: uid,
				StartTime: sTime,
				EndTime: eTime,
				CreateStr: CreateStr,
				DeliveryTime: deliveryTime,
				DeliveryStr: deliveryMsg,
			})
		}


	}else {
		uid = cycleCnf.Uid
	}
	return uid,deliveryMsg

}

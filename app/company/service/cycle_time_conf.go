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
	Uid string `json:"uid"` //唯一标识
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
		).Order("layer desc").
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
	//判断时间是否有变动,如果有变动更新uid标记.因为是要产生新的订单周期UID
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
	data.Layer = c.Layer
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


	var dayNumber int

	switch CyCleCnf.Type {
	case global.CyCleTimeDay:
		if CyCleCnf.GiveDay == 0 {
			dayNumber = 0
			if CyCleCnf.GiveTime == "" {
				CyCleCnf.GiveTime = "当天配送"
			}

		}else {
			dayNumber = CyCleCnf.GiveDay
		}

		cycleTimeValue := CalculateTime(dayNumber)

		cycleVal := fmt.Sprintf("%v %v", cycleTimeValue.Format("2006-01-02"), CyCleCnf.GiveTime)

		return cycleTimeValue, cycleVal
	case global.CyCleTimeWeek: //截止时间送达
		//通过截止的周配置 获取具体的日期

		dayNumber = CyCleCnf.GiveDay

		endTimeValue,_:=utils.GetWeekdayTimestamps(CyCleCnf.EndWeek)

		cycleTimeValue := CalculateSendTime(dayNumber,endTimeValue)

		cycleVal := fmt.Sprintf("%v %v", cycleTimeValue.Format("2006-01-02"), CyCleCnf.GiveTime)

		return cycleTimeValue, cycleVal

	default:
		return models2.XTime{}, "时间非法"
	}


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

//应该是计算出配送的时间 如果配送时间发生了变化,那这个计算的时间就是变化后的时间
//配送时间 +  大BID + 周期配送的随机ID 这个ID也就是订单的周期ID
func GetThisDayCycleUid(cid int,uid string,DeliveryTime models2.XTime) string {
	//应该是计算出配送的时间 如果配送时间发生了变化,那这个计算的时间就是变化后的时间
	//nowTime:=time.Now().Format("060102")
	cycleTime:=DeliveryTime.Format("060102")
	return fmt.Sprintf("%v_%v_%v",cycleTime,cid,uid)
}

//只有支付成功的时候才会调用这个方法
//1.记录当前时间的哪个下单的时间段,获取到这个下单开始和结束时间 和文案 记录下来
//2.计算这个配送时间也需要录入DB中

//只需要返回这个上层的UID
func CheckOrderCyCleCnfIsDb(cid int,table string,DeliveryObject models.CycleTimeConf,orm *gorm.DB) (cycleUid,deliveryMsg string)  {



	//计算出配送的周期
	deliveryTime,deliveryMsg:=GetOrderCyClyCnf(DeliveryObject)

	//检测下,然后直接返回吧
	var cycleCnf models.OrderCycleCnf

	//计算出当天+选择的周期配送UID
	cycleUid =GetThisDayCycleUid(cid,DeliveryObject.Uid,deliveryTime)

	orm.Table(table).Model(&models.OrderCycleCnf{}).Where("c_id = ? and uid = ?",cid,cycleUid).Find(&cycleCnf)

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


		if cycleCnf.Id == 0 {
			orm.Table(table).Create(&models.OrderCycleCnf{
				CId:cid,
				Uid: cycleUid,
				StartTime: sTime,
				EndTime: eTime,
				CreateStr: CreateStr,
				DeliveryTime: deliveryTime,
				DeliveryStr: deliveryMsg,
			})
		}else {
			orm.Table(table).Model(&models.OrderCycleCnf{}).Where(
				"c_id = ? and uid = ?",cid,cycleUid).Updates(map[string]interface{}{
				"delivery_str":deliveryMsg,
				"delivery_time":deliveryTime,
				"create_str":CreateStr,
				"start_time":sTime,
				"end_time":eTime,
			})
		}
	}



	return cycleUid,deliveryMsg

}

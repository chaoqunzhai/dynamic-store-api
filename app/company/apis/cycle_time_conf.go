package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"time"
)

type CycleTimeConf struct {
	api.Api
}


func (e CycleTimeConf) Switch(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	switchType:=c.Param("type")
	fmt.Println("switchType",switchType)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	switch switchType {
	case "1","2":
	default:
		e.Error(500, nil,"非法类型")
		return
	}
	var object models.CycleTimeConf
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Updates(map[string]interface{}{
		"type":switchType,
		"update_by":userDto.UserId,
	})
	e.OK("切换成功","操作成功")
	return

}
func (e CycleTimeConf) TypeCnf(c *gin.Context) {
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
	cnf := make(map[string]interface{}, 0)
	var cycleObject models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Scopes(actions.PermissionSysUser(cycleObject.TableName(),userDto)).Select("id,type").Where("enable = ? ", true).Limit(1).Find(&cycleObject)
	if cycleObject.Id == 0 {
		cnf["type"] = false
	} else {
		cnf["type"] = cycleObject.Type
	}
	e.OK(cnf, "操作成功")
	return
}

// GetPage 获取CycleTimeConf列表
// @Summary 获取CycleTimeConf列表
// @Description 获取CycleTimeConf列表
// @Tags CycleTimeConf
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param type query string false "类型,每天,每周"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CycleTimeConf}} "{"code": 200, "data": [...]}"
// @Router /api/v1/cycle-time-conf [get]
// @Security Bearer
func (e CycleTimeConf) GetPage(c *gin.Context) {
	req := dto.CycleTimeConfGetPageReq{}
	s := service.CycleTimeConf{}
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
	list := make([]models.CycleTimeConf, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}

	result := make([]map[string]interface{}, 0)

	for _, row := range list {
		cnf := map[string]interface{}{
			"id":         row.Id,
			"type":       row.Type,
			"start_time": row.StartTime,
			"end_time":   row.EndTime,
			"give_day":   row.GiveDay,
			"give_time":  row.GiveTime,
			"start_week": row.StartWeek,
			"end_week":   row.EndWeek,
			"enable":       row.Enable,
			"layer":      row.Layer,
			"created_at": row.CreatedAt,
			"order_range":service.GetOrderCreateStr(row),
			"give_range": func()  string {
				_,giveRange :=service.GetOrderCyClyCnf(row)
				return giveRange
			}(),
		}
		switch row.Type {
		case global.CyCleTimeDay:
			cnf["type_cn"] = "每天"

		case global.CyCleTimeWeek:
			cnf["type_cn"] = "每周"
		}
		result = append(result, cnf)

	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CycleTimeConf
// @Summary 获取CycleTimeConf
// @Description 获取CycleTimeConf
// @Tags CycleTimeConf
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CycleTimeConf} "{"code": 200, "data": [...]}"
// @Router /api/v1/cycle-time-conf/{id} [get]
// @Security Bearer
func (e CycleTimeConf) Get(c *gin.Context) {
	req := dto.CycleTimeConfGetReq{}
	s := service.CycleTimeConf{}
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
	var object models.CycleTimeConf

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建CycleTimeConf
// @Summary 创建CycleTimeConf
// @Description 创建CycleTimeConf
// @Tags CycleTimeConf
// @Accept application/json
// @Product application/json
// @Param data body dto.CycleTimeConfInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/cycle-time-conf [post]
// @Security Bearer
// 模式只能存在一种,不可每天,每周混用
func (e CycleTimeConf) Insert(c *gin.Context) {
	req := dto.CycleTimeConfInsertReq{}
	s := service.CycleTimeConf{}
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
	//特殊处理下
	if req.EndTime == "24:00" {
		req.EndTime = "23:59"
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var count int64

	//周期必须相同
	var cycleObject models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id =? and enable = ? ", userDto.CId, true).Limit(1).Find(&cycleObject)
	if cycleObject.Id > 0 && cycleObject.Type != req.Type {
		e.Error(500, errors.New("周期类型必须一致"), "周期类型必须一致")
		return
	}

	//时间不能重复,也就是start_time 和end_time 是不能存在相同的
	//当每天时,最多支持3个
	whereSql := ""
	switch req.Type {

	case global.CyCleTimeDay:
		whereSql = fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v'",
			userDto.CId, true, global.CyCleTimeDay, req.StartTime, req.EndTime)


		var dayCount int64
		e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and `type` = ?",userDto.CId,global.CyCleTimeDay).Count(&dayCount)
		if dayCount >= global.CyCleTimeDayMaxNumber{
			msg :=fmt.Sprintf("当前周期配置为每天, 仅支持 %v个配置",global.CyCleTimeDayMaxNumber)
			e.Error(500, errors.New(msg), msg)
			return
		}
		//检测时间是否重叠
		if e.CheckOverlapDayCycle(userDto.CId,req.StartTime,req.EndTime,0){
			e.Error(500, nil, "时间不可重叠")
			return
		}

	case global.CyCleTimeWeek:
		whereSql = fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v' and start_week = %v and end_week = %v",
			userDto.CId, true, global.CyCleTimeWeek, req.StartTime, req.EndTime, req.StartWeek, req.EndWeek)

		var dayCount int64
		e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and `type` = ?",userDto.CId,global.CyCleTimeWeek).Count(&dayCount)
		if dayCount >= global.CyCleTimeWeekMaxNumber{
			msg :=fmt.Sprintf("当前周期配置为每周, 仅支持 %v个配置",global.CyCleTimeWeekMaxNumber)
			e.Error(500, errors.New(msg), msg)
			return
		}
		//每周的时间是否吻合
		if e.CheckOverWeekCycle(userDto.CId,req.StartWeek,req.EndWeek,req.StartTime,req.EndTime,0){
			e.Error(500, nil, "时间不可重叠")
			return
		}


	default:
		e.Error(500, nil, "非法类型")
		return
	}
	e.Orm.Model(&models.CycleTimeConf{}).Where(whereSql).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("时间不可重复"), "时间不可重复")
		return
	}

	err = s.Insert(userDto.CId, &req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("时间区间创建失败,%s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

func WeekDay(t time.Weekday) int {
	addDay :=0
	switch t {
	case time.Sunday:
		addDay = 6
	case time.Monday:
		addDay = 0
	case time.Tuesday:
		addDay = 1
	case time.Wednesday:
		addDay = 2
	case time.Thursday:
		addDay = 3
	case time.Friday:
		addDay = 4
	case time.Saturday:
		addDay = 5
	}
	return  addDay
}
func (e CycleTimeConf)CheckOverWeekCycle(CID int,reqStartWeek,reqEndWeek int,reqStartTime,reqEndTime string, ExcludeId int) bool {
	// isTimeRangeOverlap 判断两个时间段是否有重叠
	var allTimeConf []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and type = ?",CID,global.CyCleTimeWeek).Find(&allTimeConf)


	for _,row:=range allTimeConf {

		if ExcludeId > 0 && row.Id == ExcludeId {
			continue
		}

		//原数据开始
		parsedStartTime, _ := time.Parse("15:04", row.StartTime)
		parsedEndTime, _ := time.Parse("15:04", row.EndTime)
		startWeek1 := time.Weekday(row.StartWeek)
		endWeek1 := time.Weekday(row.EndWeek)
		start1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), parsedStartTime.Hour(),
			parsedStartTime.Minute(), 0, 0, time.Local).AddDate(0, 0, int(startWeek1-time.Now().Weekday()) + WeekDay(startWeek1))
		end1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), parsedEndTime.Hour(),
			parsedEndTime.Minute(), 0, 0, time.Local).AddDate(0, 0, int(endWeek1-time.Now().Weekday())  + WeekDay(endWeek1))

		// 每周二 18:00 - 20:00
		startWeek2 := time.Weekday(reqStartWeek)
		endWeek2 := time.Weekday(reqEndWeek)

		fmt.Println("endWeek2!!",endWeek2)
		parsedStartTime2, _ := time.Parse("15:04", reqStartTime)
		parsedEndTime2, _ := time.Parse("15:04",reqEndTime)

		start2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), parsedStartTime2.Hour(),
			parsedStartTime2.Minute(), 0, 0, time.Local).AddDate(0, 0, int(startWeek2-time.Now().Weekday())  + WeekDay(startWeek2))
		end2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), parsedEndTime2.Hour(),
			parsedEndTime2.Minute(), 0, 0, time.Local).AddDate(0, 0, int(endWeek2-time.Now().Weekday()) + WeekDay(endWeek2))

		// 检查两个时间段是否有重叠
		overlap := utils.IsTimeOverlap(start1, end1, start2, end2)
		//fmt.Println("是否范围内:",row.Id, overlap,"|| 数据库时间",start1, end1,"|| 请求的时间",start2, end2,"reqStartWeek",reqStartWeek,reqEndWeek)
		if overlap {
			return true
		}


	}

	return false
}
//检测时间是否重叠
func (e CycleTimeConf) CheckOverlapDayCycle(CID int,reqStart,reqEnd string, ExcludeId int) bool {

	//获取所有的配置,检测是否有重复的
	var allTimeConf []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and type = ?",CID,global.CyCleTimeDay).Find(&allTimeConf)

	for _,row:=range allTimeConf{

		if ExcludeId >0 && row.Id == ExcludeId {continue}
		if row.StartTime!="" && row.EndTime!="" {

			timeRange1Start, _ := time.Parse("15:04", reqStart)
			timeRange1End, _ := time.Parse("15:04", reqEnd)

			timeRange2Start, _ := time.Parse("15:04", row.StartTime)
			timeRange2End, _ := time.Parse("15:04", row.EndTime)

			overlap := utils.IsTimeOverlap(timeRange1Start, timeRange1End, timeRange2Start, timeRange2End)

			if overlap {
				return true
			}
		}
	}
	return false


}

func (e CycleTimeConf) Update(c *gin.Context) {
	req := dto.CycleTimeConfUpdateReq{}
	s := service.CycleTimeConf{}
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
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	whereSql := ""
	var timeConf models.CycleTimeConf
	switch req.Type {

	case global.CyCleTimeDay:

		whereSql = fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v'",
			userDto.CId, true, global.CyCleTimeDay, req.StartTime, req.EndTime)


		 if e.CheckOverlapDayCycle(userDto.CId,req.StartTime,req.EndTime,req.Id){
			 e.Error(500, nil, "时间不可重叠")
			 return
		 }
	case global.CyCleTimeWeek:
		whereSql = fmt.Sprintf("c_id = %v and enable = %v and type = %v and start_time = '%v' and end_time = '%v' and start_week = %v and end_week = %v",
			userDto.CId, true, global.CyCleTimeWeek, req.StartTime, req.EndTime, req.StartWeek, req.EndWeek)


		if e.CheckOverWeekCycle(userDto.CId,req.StartWeek,req.EndWeek,req.StartTime,req.EndTime,req.Id){
			e.Error(500, nil, "时间不可重叠")
			return
		}

	default:
		e.Error(500, nil, "非法类型")
		return
	}
	e.Orm.Model(&models.CycleTimeConf{}).Where(whereSql).Limit(1).Find(&timeConf)
	if timeConf.Id > 0 && timeConf.Id != req.Id {
		e.Error(500, nil, "时间不可重复")
		return
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

func (e CycleTimeConf) Delete(c *gin.Context) {
	s := service.CycleTimeConf{}
	req := dto.CycleTimeConfDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

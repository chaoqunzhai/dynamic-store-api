package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/qiniu"
	"go-admin/config"
	"go-admin/global"
	"path"
	"time"

	models2 "go-admin/cmd/migrate/migration/models"

	"go-admin/common/business"
	"go-admin/common/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/jwt/user"
	"go-admin/common/redis_worker"
)

type Worker struct {
	api.Api
}
type WorkerExportReq struct {
	Order   []string   `json:"order"`
	Cycle int `json:"cycle"`
	Type int `json:"type"`
	LineId []int `json:"line_id"`
	LineName []string `json:"line_name"`
}
type GetPageReq struct {
	dto.Pagination `search:"-"`
	Type string `json:"type" form:"type" search:"-" `
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_tasks" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_tasks" comment:"创建时间"`

}
func (m *GetPageReq) GetNeedSearch() interface{} {
	return *m
}

//获取大B的下载中心任务队列
func (e *Worker)Get(c *gin.Context)  {
	req := GetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models2.CompanyTasks
	list:=make([]models2.CompanyTasks,0)
	var count int64
	orm:=e.Orm.Model(&data).Where("c_id = ?",userDto.CId).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		)

	switch req.Type {
	case "report":
		orm.Where("`type` > ?",global.ExportTypeOrder)
	case "export":
		orm.Where("`type` = ?",global.ExportTypeOrder)
	}
	err = orm.Order("id desc").
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error
	//直接读取DB中的数据

	result:=make([]map[string]interface{},0)
	for _,row:=range list{
		result = append(result, map[string]interface{}{
			"create_time":row.CreatedAt.Format("2006-01-02 15:04:05"),
			"path":row.Path,
			"status":row.Status,
			"type":row.Type,
			"id":row.Id,
			"title":row.Title,
		})
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), fmt.Sprintf("云端文件最多保留%v天",config.ExtConfig.ExportDay))
}


func (e *Worker)Download(c *gin.Context)  {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	uid:=c.Param("uid")
	var data models2.CompanyTasks
	e.Orm.Model(&models2.CompanyTasks{}).Where("c_id = ? and id = ?",userDto.CId,uid).Limit(1).Find(&data)
	if data.Id == 0 {
		e.Error(500, nil,"数据不存在")
		return
	}

	downloadUrl :=config.ExtConfig.CloudObsUrl + path.Join(fmt.Sprintf("%v",userDto.CId),global.CloudExportOrderFilePath,data.Path)

	e.OK(downloadUrl,"")
	return

}

func (e *Worker)Create(c *gin.Context)  {
	req:=WorkerExportReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if req.Type == 0 {
		if len(req.Order) == 0 {
			e.Error(500, nil, "请选择订单ID")
			return
		}
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//查询队列限制
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "export_worker", e.Orm)
	MaxNumber := CompanyCnf["export_worker"]
	//获取订单的导出队列
	thisWorker:=redis_worker.GetExportQueueLength(userDto.CId,global.WorkerOrderStartName)
	if thisWorker >= MaxNumber {
		e.Error(500, nil,fmt.Sprintf("最多同时支持%v个任务执行,请稍后重试",MaxNumber))
		return
	}
	//用一个key来做防止大量的重复提交,给DB带来非必要的查询
	//如果是订单 那就查询 当前时间戳  精确到分 ,查询到这个key 已经存在 + 处理中   = 返回请一分钟后重试

	//如果是汇总 查询cid:周期ID
	var Queue string
	var title string
	var LineExport int
	//2分钟限制一次
	mathKey := time.Now().Add(-2 * time.Minute).Format("200601021504")

	switch req.Type {
	case global.ExportTypeOrder:
		Queue = global.WorkerOrderStartName
		title = "选中导出配送订单"
		mathKey = fmt.Sprintf("%v_order",mathKey)
	case global.ExportTypeSummary:
		title = "导出配送汇总表"
		Queue = global.WorkerReportSummaryStartName
		mathKey = fmt.Sprintf("%v_summary",mathKey)
	case global.ExportTypeLine:
		if len(req.LineName) > 1{
			title = fmt.Sprintf("导出【%v】条路线表",len(req.LineName))
		}else {
			title = fmt.Sprintf("导出【%v】路线表",req.LineName[0])
		}
		LineExport = 0
		Queue = global.WorkerReportLineStartName
		mathKey = fmt.Sprintf("%v_line",mathKey)
	case global.ExportTypeLineShopDelivery:
		if len(req.LineName) > 1{
			title = fmt.Sprintf("导出【%v】条路线配送表",len(req.LineName))
		}else {
			title = fmt.Sprintf("导出【%v】路线配送表",req.LineName[0])
		}
		LineExport = 1
		Queue = global.WorkerReportLineDeliveryStartName
		mathKey = fmt.Sprintf("%v_delivery",mathKey)
	}
	var count int64

	e.Orm.Model(&models2.CompanyTasks{}).Where("`key` = ? and c_id = ? and type = ?",mathKey,userDto.CId,req.Type).Count(&count)

	if count > 0 {
		e.Error(500, nil,"请勿在2分钟内重复提交相同任务")
		return
	}

	taskTable:=models2.CompanyTasks{
		Title: title,
		CreateBy: userDto.UserId,
		CId:      userDto.CId,
		Type:     req.Type,
		Key:      mathKey,
		Status:   0,
	}
	//记录在DB中
	e.Orm.Create(&taskTable)
	exportReq :=global.ExportRedisInfo{
		CId: userDto.CId,
		Order: req.Order,
		Cycle: req.Cycle,
		OrmId: taskTable.Id,
		LineId: req.LineId,
		ExportUser: userDto.Username,
		ExportTime: time.Now().Format("2006-01-02 15:04:05"),
		Queue: Queue,
		LineExport: LineExport,
	}

	//先发送到redis中
	_=redis_worker.SendExportQueue(exportReq)
	e.OK("","任务创建成功.请前往任务中心查看 ！")
	return
}


func (e *Worker)Remove(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	uid:=c.Param("uid")
	var data models2.CompanyTasks
	e.Orm.Model(&models2.CompanyTasks{}).Where("c_id = ? and id = ?",userDto.CId,uid).Limit(1).Find(&data)
	if data.Id == 0 {
		e.Error(500, nil,"数据不存在")
		return
	}

	e.Orm.Model(&data).Where("id = ?",data.Id).Delete(&data)
	//
	obsUrl :=path.Join(fmt.Sprintf("%v",userDto.CId),global.CloudExportOrderFilePath,data.Path)

	//数据删除
	buckClient:=qiniu.QinUi{
		CId: userDto.CId,
	}
	buckClient.InitClient()

	buckClient.RemoveFile(obsUrl)
	e.OK("","删除成功")
	return

}
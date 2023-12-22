package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

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
type WorkerOrderId struct {
	Order   []string   `json:"order"`
	Type int `json:"type"`
}
type GetPageReq struct {
	dto.Pagination `search:"-"`
	Type int `json:"type" form:"type" search:"type:exact;column:type;table:company_tasks" `
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
	err = e.Orm.Model(&data).Where("c_id = ?",userDto.CId).
		Scopes(
		cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).Order("id desc").
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
		})
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e *Worker)Create(c *gin.Context)  {
	req:=WorkerOrderId{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if len(req.Order) == 0 {
		e.Error(500, nil, "请选择订单ID")
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//查询队列限制
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "export_worker", e.Orm)
	MaxNumber := CompanyCnf["export_worker"]
	thisWorker:=redis_worker.GetExportQueueLength(userDto.CId)
	if thisWorker >= MaxNumber {
		e.Error(500, nil,fmt.Sprintf("最多同时支持%v个任务执行,请稍后重试",MaxNumber))
		return
	}
	taskTable:=models2.CompanyTasks{
		CreateBy: userDto.UserId,
		CId:userDto.CId,
		Type: req.Type,
		Status: 0,
	}
	//记录在DB中
	e.Orm.Create(&taskTable)

	exportReq :=redis_worker.ExportReq{
		CId: userDto.CId,
		Order: req.Order,
		OrmId: taskTable.Id,
	}
	//先发送到redis中
	redis_worker.SendExportQueue(exportReq)
	e.OK("","任务创建成功.请前往任务中心查看 ！")
	return
}
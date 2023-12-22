package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/jwt/user"
	"go-admin/common/redis_worker"
)

type Worker struct {
	api.Api
}
type WorkerOrderId struct {
	Order   []string   `json:"order" form:"order"`
}
//获取大B的下载中心任务队列
func (e *Worker)Get(c *gin.Context)  {

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

	fmt.Println("获取消息队列",userDto)
	//直接读取DB中的数据

	return
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
		Type: 0,
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
	e.OK("","任务创建成功.耐心等待...")
	return
}
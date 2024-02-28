/**
@Author: chaoqun
* @Date: 2024/1/17 11:16
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/common/crontab"
	"go.uber.org/zap"
)

type Crontab struct {
	api.Api
}


func (e *Crontab)SyncOrder(c *gin.Context)  {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	zap.S().Infof("开始进行 订单状态自动化变更")

	crontabMain :=crontab.ThisDaySyncOrderCycle{
		Orm: e.Orm,
	}
	crontabMain.RunCompanySplitOrderSync()

	e.OK("","操作成功")
	return
}
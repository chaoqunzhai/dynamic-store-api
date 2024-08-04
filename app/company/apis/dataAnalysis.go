/**
@Author: chaoqun
* @Date: 2024/8/4 17:16
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/service"
)
type DataAnalysis struct {
	api.Api
}

// 销售统计

func (e DataAnalysis) Count(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	return
}

func (e DataAnalysis) Gross(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	return
}

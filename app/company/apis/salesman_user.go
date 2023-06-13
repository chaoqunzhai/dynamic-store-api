package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

type SaleUser struct {
	api.Api
}
//因为用户的所有信息都在一张表里面,所以呢，应该是给用户打一个标签,
//例如哪些用户是业务员,又是小B什么的
func (e SaleUser) List(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}


}
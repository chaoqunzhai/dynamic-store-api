package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/common/global"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":200,
		"msg": fmt.Sprintf( "dcy-store %s",global.Version),
	})
}

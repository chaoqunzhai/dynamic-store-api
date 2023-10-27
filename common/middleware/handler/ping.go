package handler

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":200,
		"msg": "dcy-store.api",
	})
}

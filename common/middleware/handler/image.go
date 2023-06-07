package handler

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/business"
	"go-admin/global"
	"io/ioutil"
)

func ImageShow(c *gin.Context) {
	t := c.Param("type")
	name := c.Param("name")
	companyId := c.Param("cid")
	path := ""
	switch t {
	case global.GoodsPath:
		path = business.GetGoodPathName(companyId) + name

	}

	file, readError := ioutil.ReadFile(path)
	if readError != nil {
		c.JSON(200, gin.H{
			"message": "ok",
		})
		return
	}
	c.Writer.WriteString(string(file))

}

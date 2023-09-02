package apis

import (
	"github.com/gin-gonic/gin"
)

func GoAdmin(c *gin.Context) {
	c.String(200, "dcy-store.api")
}

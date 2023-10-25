package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go-admin/app/admin/service/dto"
	"go-admin/common"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		// 开始时间
		startTime := time.Now()
		// 处理请求
		var body string
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete:
			bf := bytes.NewBuffer(nil)
			wt := bufio.NewWriter(bf)
			_, err := io.Copy(wt, c.Request.Body)
			if err != nil {
				log.Warnf("copy body error, %s", err.Error())
				err = nil
			}
			rb, _ := ioutil.ReadAll(bf)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rb))
			body = string(rb)
		}

		c.Next()
		url := c.Request.RequestURI
		if strings.Index(url, "logout") > -1 ||
			strings.Index(url, "login") > -1 {
			return
		}
		// 结束时间
		endTime := time.Now()
		if c.Request.Method == http.MethodOptions {
			return
		}

		rt, bl := c.Get("result")
		var result = ""
		if bl {
			rb, err := json.Marshal(rt)
			if err != nil {
				log.Warnf("json Marshal result error, %s", err.Error())
			} else {
				result = string(rb)
			}
		}

		st, bl := c.Get("status")
		var statusBus = 0
		if bl {
			statusBus = st.(int)
		}

		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := common.GetClientIP(c)
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 日志格式
		logData := map[string]interface{}{
			"statusCode":  statusCode,
			"latencyTime": latencyTime,
			"clientIP":    clientIP,
			"method":      reqMethod,
			"uri":         reqUri,
		}
		log.WithFields(logData).Info()

		if c.Request.Method != "OPTIONS" && config.LoggerConfig.EnabledDB && statusCode != 404 {
			SetDBOpenerLog(c, clientIP, statusCode, reqUri, reqMethod, latencyTime, body, result, statusBus)
		}
	}
}

// SetDBOpenerLog 数据操作的时候,把操作的内容写入到到日志中
func SetDBOpenerLog(c *gin.Context, clientIP string, statusCode int, reqUri string, reqMethod string, latencyTime time.Duration, body string, result string, status int) {

	l := make(map[string]interface{})
	l["_fullPath"] = c.FullPath()
	l["operUrl"] = reqUri
	l["operIp"] = clientIP
	l["operName"] = user.GetUserName(c)
	l["requestMethod"] = reqMethod
	l["operParam"] = body
	l["operTime"] = time.Now()
	l["jsonResult"] = result
	l["latencyTime"] = latencyTime.String()
	l["statusCode"] = statusCode
	l["userAgent"] = c.Request.UserAgent()
	l["createBy"] = user.GetUserId(c)
	l["updateBy"] = user.GetUserId(c)
	if status == http.StatusOK {
		l["status"] = dto.OperaStatusEnable
	} else {
		l["status"] = dto.OperaStatusDisable
	}

	fmt.Println("lll!!",l)
}

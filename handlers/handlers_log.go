package handlers

import (
	"github.com/gin-gonic/gin"
	"login/logging"
	"time"
)

// LogStat 记录日志
func LogStat(logName string, c *gin.Context, t1 time.Time) {
	r := c.Request
	logging.Debugf("%v: request: url:%v client_ip:%v body:%v cost:%v\n",
		logName, r.RequestURI, c.ClientIP(), r.PostForm.Encode(), time.Since(t1))
}

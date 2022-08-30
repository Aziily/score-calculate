package service

import (
	"encoding/base64"
	"net/http"
	"score-calculate/logger"
	"sync"
	"unsafe"

	"github.com/gin-gonic/gin"
)

var (
	lock sync.Mutex
)

func init() {
	service = startService()
}

func ShutDown() {
	service.Stop()
}

func StartService() {
	router := gin.Default()
	router.Use(CrosHandler())
	router.Use(logger.LoggerToFile())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/search", func(ctx *gin.Context) {
		lock.Lock()
		var input QueryInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			logger.Logger().WithField(
				"Place", "load",
			).Error("获取传入json错误")
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			lock.Unlock()
			return
		}
		decodePasswd, _ := base64.StdEncoding.DecodeString(input.Passwd)
		input.Passwd = *(*string)(unsafe.Pointer(&decodePasswd))

		result := GetAllMess(input.Id, input.Passwd)
		if len(result) == 0 {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "wrong input of id or passwd"},
			)
			lock.Unlock()
			return
		}

		ctx.JSON(http.StatusOK, result)
		lock.Unlock()
	})
	router.Run(":9510")
}

func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json")

		//处理请求
		context.Next()
	}
}

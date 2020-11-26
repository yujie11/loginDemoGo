package my_middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"login/handlers"
	"login/logging"
	"login/models"
	"login/utils/config"
	"login/utils/middlewares"
	"net/http"
)

type Result handlers.Result
//初始化中间件
func InitMiddleware(router *gin.Engine, conf *config.Config){
	router.Use(gin.Recovery())
	router.Use(middlewares.SetMiddleware(handlers.MiddlewareConfig, conf))
	dbServer, ok := conf.DBServerConf(handlers.MiddlewareLoginDB)
	if !ok {
		panic(fmt.Sprintf("InitMiddleware: %v配置不存在\n", handlers.MiddlewareLoginDB))
	}
	logindborm, err := dbServer.NewGormDB(10)
	if err != nil {
		panic(fmt.Sprintf("InitMiddleware：%v数据库连接信息错误 err:%v\n", handlers.MiddlewareLoginDBORM, err))
	}
	router.Use(middlewares.SetMiddleware(handlers.MiddlewareLoginDBORM, logindborm))
	logindbsqlx, err := dbServer.NewPostgresDB(10)
	if err != nil {
		panic(fmt.Sprintf("InitMiddleware：%v数据库连接信息错误 err:%v\n", handlers.MiddlewareLoginDBSQLX, err))
	}
	router.Use(middlewares.SetMiddleware(handlers.MiddlewareLoginDBSQLX, logindbsqlx))
	loginredis, err := conf.NewRedisPool(handlers.MiddlewareLoginREDIS, 10)
	if err != nil {
		panic(fmt.Sprintf("InitMiddleware：%vRedis连接信息错误 err:%v\n", handlers.MiddlewareLoginREDIS, err))
	}
	router.Use(middlewares.SetMiddleware(handlers.MiddlewareLoginREDIS, loginredis))
}

func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Content-Type","application/json")
		next.ServeHTTP(w, r)
	})
}
//生成token
func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mobilenumber" : user.MobileNumber,
		"username": user.UserName,
	})
	return token.SignedString([]byte("secret"))
}

//解析token
func ParseToken(tokenStr string) (cliams jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenStr, func( *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		log.Println("[ParseToken] parse token error,err_msg = %s", err.Error())
		return
	}
	var ok bool
	cliams, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		logging.Errorf("[ParseToken]token claims type error,can't convert")
		return cliams, fmt.Errorf("token claims can't convert to JwtCustomerClaims")
	}
	return
}

//添加token中间件
func TokenMiddleware(c *gin.Context) {
	tokenStr := c.Request.Header.Get("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": Result{
				Message: "未登陆",
				Status:  0,
			},
		})
		//终止后面的HandlerFunc
		c.Abort()
		return
	}
	claims, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": Result{
				Message: "解析token失败",
				Status:  0,
			},
		})
		c.Abort()
		return
	}
	c.Set("mobilenumber", claims["mobilenumber"])
	c.Set("username", claims["username"])
	c.Next()
}

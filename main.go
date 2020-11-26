package main

import (
	"flag"
	"login/logging"
	"login/utils/config"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"login/handlers"
	"login/my_middlewares"
	"net/http"
)
type Config config.Config

var (
	configPath string
	conf       = &Config{}
)

func init() {
	//读取配置文件
	flag.StringVar(&configPath, "conf", "config.toml", "数据库配置文件")
}

func main() {
	conf, err := config.New(configPath)
	if err != nil {
		//数据库配置信息错误
		logging.Errorf("解析配置文件出错")
		return
	}
	r := gin.New()
	//初始化中间件
	my_middlewares.InitMiddleware(r, conf)
	v1 := r.Group("/api")
	{
		//用户注册
		v1.POST("/user/add", handlers.UserRegister)
		v1m := v1.Group("", my_middlewares.TokenMiddleware)
		//更新用户信息:修改密码或用户名
		v1m.POST("/user/update", handlers.UserUpdate)
		//获取用户自身信息
		v1m.POST("/user/get", handlers.UserGet)
		//用户列表
		v1m.POST("/user/list", handlers.UserGetList)
	}
	v2 := v1.Group("/user", my_middlewares.TokenMiddleware)
	{
		v2.GET("", handlers.UserGet)
		v2.POST("", handlers.UserRegister)
		v2.PUT("", handlers.UserUpdate)
	}
	r.LoadHTMLFiles("./statics/login.html", "./statics/index.html", "./statics/register.html")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", handlers.UserLogin)
	r.GET("/register", func(c *gin.Context) {
		//注册页面
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.POST("/register", handlers.UserRegister)
	r.POST("/verify", handlers.Verify)
	r.Run()
}

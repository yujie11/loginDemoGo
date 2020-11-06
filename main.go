package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"login/handles"
	"login/middlewares"
	"login/utils/config"
	"net/http"
)

//Json返回结果
type Result middlewares.Result

var DB *gorm.DB
var Rdb *redis.Client

func main() {
	//数据库连接配置信息
	DBInfo := config.DBServer{
		Host:     "127.0.0.1",
		Port:     5432,
		DBName:   "test",
		User:     "",
		Password: ""}
	var err error
	DB, err = DBInfo.NewGromDB()
	if err != nil {
		fmt.Println("数据库连接失败")
		return
	}
	//连接redis
	Rdb = config.NewRedis()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	v1 := r.Group("/api")
	{
		//更新用户信息:修改密码或用户名
		v1.POST("/user/update", middlewares.TokenMiddleware(), handles.UserUpdate(DB, Rdb))
		//获取用户自身信息
		v1.POST("/user/get", middlewares.TokenMiddleware(), handles.UserGet(DB, Rdb))
		//用户列表
		v1.POST("/user/list", middlewares.TokenMiddleware(), handles.UserGetList(DB, Rdb))
		//增加用户
		v1.POST("/user/add", handles.Userregister(DB, Rdb))
	}
	v2 := v1.Group("")
	{
		v2.GET("/user", middlewares.TokenMiddleware(), handles.UserGet(DB, Rdb))
		v2.POST("/user", middlewares.TokenMiddleware(), handles.Userregister(DB, Rdb))
		v2.PUT("/user", middlewares.TokenMiddleware(), handles.UserUpdate(DB, Rdb))
	}
	r.LoadHTMLFiles("./statics/login.html", "./statics/index.html", "./statics/register.html")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", handles.UserLogin(DB, Rdb))
	r.GET("/register", func(c *gin.Context) {
		//注册页面
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.POST("/register", handles.Userregister(DB, Rdb))
	r.POST("/verify", handles.Verify(Rdb))
	r.Run()
	defer DB.Close()
}

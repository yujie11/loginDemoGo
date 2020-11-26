package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"login/models"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//用户注册
func UserRegister(c *gin.Context) {
	DB := c.MustGet(MiddlewareLoginDBORM).(*gorm.DB)
	redisPool := c.MustGet(MiddlewareLoginREDIS).(*redis.Pool)
	Rdb := redisPool.Get()
	//注册操作
	var user models.User
	user.MobileNumber = c.PostForm("mobilenumber")
	user.UserName = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.UserName = c.PostForm("username")
	password := c.PostForm("checkpassword")
	if strings.Compare(password, user.Password) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "两次输入密码不一致",
				Status:  0,
			},
		})
		return
	}
	verify := c.PostForm("verify") //验证码
	//checkverify, _ := Rdb.Get(user.MobileNumber).Result()
	checkverify, err := redis.String(Rdb.Do("GET", user.MobileNumber))
	fmt.Println(checkverify)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "验证码已过期",
				Status:  0,
			},
		})
		return
	}
	if verify == "" || checkverify == "" || strings.Compare(verify, checkverify) != 0 {
		//验证码输入错误
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "验证码输入错误或超时",
				Status:  0,
			},
		})
		return
	}
	if num, _, _ := models.UserIsExist(DB, user.MobileNumber); num == 0 {
		//存储到数据库中
		_, err := models.AddUser(DB, user.MobileNumber, user.UserName, user.Password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": Result{
					Message: err.Error(),
					Status:  0,
				},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "注册成功",
				Status:  1,
			},
		})
	} else {
		//用户已存在
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "用户已存在",
				Status:  0,
			},
			"user": models.User{
				MobileNumber: user.MobileNumber,
				UserName:     user.UserName,
			},
		})
	}
}

//模拟发送验证码
func Verify(c *gin.Context) {
	redisPool := c.MustGet(MiddlewareLoginREDIS).(*redis.Pool)
	Rdb := redisPool.Get()
	mobilenumber := c.PostForm("mobilenumber")
	//初始化随机数的资源库
	rand.Seed(time.Now().UnixNano())
	num := rand.Int31n(10000)
	_, err := Rdb.Do("SET", mobilenumber, fmt.Sprintf("%d", num))
	if err != nil {
		log.Println(err)
	}
	//Rdb.Set(mobilenumber, num, 60*time.Second)
	c.JSON(http.StatusOK, gin.H{
		"result": Result{
			Message: "ok",
			Status:  1,
		},
		"verify": num,
	})
}


func UserLogin(c *gin.Context) {
	DB := c.MustGet(MiddlewareLoginDBORM).(*gorm.DB)
	mobileNumber := c.PostForm("mobilenumber")
	password := c.PostForm("password")
	password = models.PassAddMd5(password)
	if _, user, err := models.UserIsExist(DB, mobileNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": Result{
				Message: "手机号信息错误",
				Status:  0,
			},
		})
		fmt.Println("手机号码错误：", err)
		return
	} else {
		if user.Password == password {
			if token, err := GenerateToken(user); err != nil {
				log.Println(err)
				c.JSON(http.StatusOK, gin.H{
					"result": Result{
						Message: "token 生成失败",
						Status:  0},
				})
			} else {
				//登陆成功
				c.JSON(http.StatusOK, gin.H{
					"result": Result{
						Message: "登陆成功",
						Status:  1},
					"token": token,
				})
			}
		} else {
			//密码错误
			c.JSON(http.StatusOK, gin.H{
				"result": Result{
					Message: "密码错误",
					Status:  0},
			})
		}
	}
}

//生成token
func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mobilenumber" : user.MobileNumber,
		"username": user.UserName,
	})
	return token.SignedString([]byte("secret"))
}

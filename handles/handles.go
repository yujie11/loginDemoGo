package handles

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"log"
	"login/middlewares"
	"login/models"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Result middlewares.Result

//用户注册
func Userregister(DB *gorm.DB , Rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		checkverify, _ := Rdb.Get(user.MobileNumber).Result()
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
}

func Verify(Rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		mobilenumber := c.PostForm("mobilenumber")
		//初始化随机数的资源库
		rand.Seed(time.Now().UnixNano())
		num := rand.Int31n(10000)
		Rdb.Set(mobilenumber, num, 60*time.Second)
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "ok",
				Status:  1,
			},
			"verify": num,
		})
	}
}

func UserLogin(DB *gorm.DB, Rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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
				if token, err := middlewares.GenerateToken(user); err != nil {
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
}

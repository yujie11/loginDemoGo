package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"login/models"
	"net/http"
)
//用户更新
func UserUpdate(c *gin.Context) {
	DB := c.MustGet(MiddlewareLoginDBORM).(*gorm.DB)
	var user models.User
	var err error
	mobilenumber, ok := c.Get("mobilenumber")
	user.MobileNumber = mobilenumber.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "用户信息登陆错误",
				Status:  0,
			},
		})
		return
	}
	user.UserName = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.Password = models.PassAddMd5(user.Password)
	if user.Validate() != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: fmt.Sprintf("更新的信息不符合:%v", err),
				Status:  0,
			},
		})
		return
	}
	_, err = user.UpdateUser(DB)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: fmt.Sprintf("更新数据库失败:%v", err),
				Status:  0,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": Result{
			Message: "OK",
			Status:  1,
		},
		"user": models.User{
			MobileNumber: user.MobileNumber,
			UserName:     user.UserName,
		},
	})
}

//获取用户信息
func UserGet(c *gin.Context) {
	DB := c.MustGet(MiddlewareLoginDBORM).(*gorm.DB)
	mobilenumber, ok := c.GetQuery("mobilenumber")
	if mobilenumber == "" {
		mobilenumber = c.PostForm("mobilenumber")
	}
	if !ok && mobilenumber == "" {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "用户信息登陆错误",
				Status:  0,
			},
		})
		return
	}
	user, err := models.GetUserByMobileNumber(DB, mobilenumber)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: fmt.Sprintf("查询手机号%v信息错误:%v", mobilenumber, err),
				Status:  0,
			},
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "OK",
				Status:  1,
			},
			"user": models.User{
				MobileNumber: user.MobileNumber,
				UserName:     user.UserName,
			},
		})
	}
}

//获取用户列表
func UserGetList(c *gin.Context) {
	DB := c.MustGet(MiddlewareLoginDBORM).(*gorm.DB)
	users, err := models.GetUserList(DB, "", "", "")
	if err != nil {
		log.Println(err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": Result{
				Message: "成功获取用户列表",
				Status:  1,
			},
			"users": users,
		})
	}
}
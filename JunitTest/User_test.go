package test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"login/models"
	"testing"
)

var DB *gorm.DB

func init(){
	var err error
	DB, err = gorm.Open(
		"postgres",
		"host=127.0.0.1 dbname=test sslmode=disable",
	)
	if err != nil {
		panic(err)
	}
}

func TestGetUserByMobileNumber(t *testing.T){
	defer DB.Close()
	mobileNumber := "13767503027"
	user, err := models.GetUserByMobileNumber(DB, mobileNumber)
	if err != nil {
		panic(err)
	}else {
		fmt.Println(user)
	}
}


func TestAddUser(t *testing.T){
	user := new(models.User)
	user.MobileNumber = "18970945184"
	user.UserName = "添加测试"
	user.Password = models.PassAddMd5("18970945184")
	models.AddUser(DB,user.MobileNumber,user.UserName,user.Password)
	DB.Commit()
	defer DB.Close()
}

func TestCheckPassword(t *testing.T){
	password := "13767503027"
	mobileNumber := "13767503027"
	user,err := models.GetUserByMobileNumber(DB,mobileNumber)
	if err != nil {
		fmt.Println("Get user by mobile number is error", err)
	}
	if user.CheckPassword(password) {
		fmt.Println("密码正确")
	}else {
		fmt.Println("密码错误")
	}
	defer DB.Close()
}


func TestUserIsExit(t *testing.T){
	mobileNumber := "13767503027"
	num, user, err := models.UserIsExist(DB, mobileNumber)
	if num ==0 && err !=nil {
		fmt.Println("用户号码为", mobileNumber, "的用户不存在")
	}else {
		fmt.Println(user)
	}
	defer DB.Close()
}

func TestGetUserList(t *testing.T) {
	users, err := models.GetUserList(DB,"", "", "")
	if err != nil {
		panic(err)
	}else {
		fmt.Println(users)
	}
	DB.Close()
}

func TestUpdateUser(t *testing.T) {
	user := models.User{
		MobileNumber: "18970945184",
		UserName: "更新测试",
		Password: models.PassAddMd5("18970945185"),
		}
		var err error
		preUpdateUser ,err := models.GetUserByMobileNumber(DB, user.MobileNumber)
		if err != nil {
			panic(err)
		}
		_, err = user.UpdateUser(DB)
		if err !=nil {
			panic(err)
		}
		updateUser, err := models.GetUserByMobileNumber(DB, user.MobileNumber)
		if err != nil {
			panic(err)
		}
		fmt.Println("更新前数据库的信息:", preUpdateUser)
		fmt.Println("需要更新的用户:", user)
		fmt.Println("更新后数据库的数据:", updateUser)
		defer DB.Close()
}
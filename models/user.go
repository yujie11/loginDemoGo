package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"strings"
)

//用户
type User struct {
	MobileNumber string `gorm:"column:mobile_number"; json:"mobileNumber"`
	UserName     string `gorm:"column:user_name" json:"userName"`
	Password     string `gorm:"column:password" json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

func (user *User)Prepare(){
	user.MobileNumber = strings.TrimSpace(user.MobileNumber)
	user.UserName = strings.TrimSpace(user.UserName)
	user.Password = strings.TrimSpace(user.Password)
}

func (user *User) Validate() error{
	if len(user.MobileNumber) != 11 {
		return errors.New("MobileNumber is error")
	}
	if user.UserName == "" {
		return errors.New("UserName of venue is required")
	}
	if user.Password == "" {
		return errors.New("password of venue is required")
	}
	return nil
}

func PassAddMd5(password string) string {
	w := md5.New()
	io.WriteString(w, password)
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
	return md5str2
}

func (user *User) CheckPassword(loginPassword string) bool {
	str := user.Password
	md5LoginPassword := PassAddMd5(loginPassword)
	fmt.Println(md5LoginPassword)
	if md5LoginPassword == str {
		return true
	}
	return false
}

func UserIsExist(db *gorm.DB, mobileNumber string) (int32, *User, error) {
	data := &User{}
	var total int32
	if err := db.Debug().Table("users").Where("mobile_number = ?",mobileNumber).Count(&total).First(data).Error; err != nil {
		return 0, nil, err
	}
	return total, data, nil
}


func AddUser(db *gorm.DB, mobileNumber, userName,password string) (*User, error) {
	user :=  new(User)
	if len(mobileNumber) != 11 {
		return nil, errors.New("mobile number is error")
	}
	user.MobileNumber = mobileNumber
	user.UserName = userName
	user.Password = PassAddMd5(password)
	err := db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user,nil
}

func GetUserList(db *gorm.DB, mobileNumber, userName, password string) (data []User, err error) {
	sql := db.Debug().Table("users")
	if mobileNumber != "" {
		sql = sql.Where("mobile_number like ?",fmt.Sprintf("%"+mobileNumber+"%"))
	}
	if userName != "" {
		sql = sql.Where("user_name like ?", fmt.Sprintf("%"+userName+"%"))
	}
	if password != "" {
		sql = sql.Where("password like ?", fmt.Sprintf("%"+password+"%"))
	}
	//if err = sql.Order("mobile_number asc").Find(&data).Error; err != nil {
	//	return nil, err
	//}

	sql.Find(&data)
	return data, nil
}

func GetUserByMobileNumber(db *gorm.DB, mobileNumber string) (*User, error) {
	user := new(User)
	if err := db.Debug().Table("users").Where("mobile_number = ?",mobileNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//修改密码或用户名
func (user *User) UpdateUser(db *gorm.DB) (*User, error){
	if err := db.Table("users").Where("mobile_number = ?",user.MobileNumber).Update(User{
		MobileNumber: user.MobileNumber,
		UserName:     user.UserName,
		Password:     user.Password}).Error; err != nil {
		return user, err
	}
	return user, nil
}



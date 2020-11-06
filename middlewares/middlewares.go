package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"login/models"
	"net/http"
	"strings"
)

//Json返回结果
type Result struct {
	Message string `json:"msg"`
	Status int `json:"status"`
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
		log.Println("[ParseToken]token claims type error,can't convert")
		return cliams, fmt.Errorf("token claims can't convert to JwtCustomerClaims")
	}
	return
}

//添加token中间件
func TokenMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("token")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized,gin.H{
				"result": Result{
					Message: "未登陆",
					Status: 0,
				},
					})
			//终止后面的HandlerFunc
			c.Abort()
			return
		}
		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized,gin.H{
				"result": Result{
					Message: "解析token失败",
					Status: 0,
				},
			})
			c.Abort()
			return
		}
		mobilenumber, ok := c.GetQuery("mobilenumber")
		if !ok {
			c.JSON(http.StatusUnauthorized,gin.H{
				"result": Result{
					Message: "未登陆",
					Status: 0,
				},
			})
			c.Abort()
			return
		}
		number := claims["mobilenumber"].(string)
		if strings.Compare(number,mobilenumber) != 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"result": Result{
					Message: "token有误",
					Status: 0,
				},
			})
			c.Abort()
			return
		}
		c.Set("mobilenumber",claims["mobilenumber"])
		c.Set("username",claims["username"])
		c.Next()
	}
}

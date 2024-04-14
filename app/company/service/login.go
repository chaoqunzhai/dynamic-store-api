/**
@Author: chaoqun
* @Date: 2024/1/22 11:33
*/
package service

import (
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	Identity int `json:"identity"`
	UserId int    `json:"user_id"`
	Phone  string `json:"phone"`
	UserName string `json:"user_name"`
	Enable bool `json:"enable"`
	jwt.StandardClaims
}

func BuildToken(userId int,username, phone string) (tokenString string, expire time.Time, err error) {
	//只需生成token即可,无需给token设置一些值,因为是需要实时查询的
	// 定义过期时间,7天后过期
	expireToken := time.Duration(config.JwtConfig.Timeout) * time.Second
	ExpiresAt := time.Now().Add(expireToken)
	//fmt.Println("过期时间为:", ExpiresAt)
	claims := &Claims{
		Identity: userId,
		UserId: userId,
		Phone:  phone,
		UserName:username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpiresAt.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "dcyServer",     // 签名颁发者
			Subject:   "dynamic-store", //签名主题
		},
		Enable: true,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(config.JwtConfig.Secret))
	if err != nil {
		return "", ExpiresAt, err
	}

	return tokenString, ExpiresAt, err
}
// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string)(*Claims,error){

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtConfig.Secret), nil
	})

	if tokenClaims!=nil{
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims,ok:=tokenClaims.Claims.(*Claims);ok&&tokenClaims.Valid{
			return claims,nil
		}
	}
	return nil,err
}
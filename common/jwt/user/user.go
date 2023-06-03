package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	sys "go-admin/app/admin/models"
	"gorm.io/gorm"
)

func GetUserId(c *gin.Context) int {
	data := ExtractClaims(c)

	//fmt.Println("GetUserIdData",data,"identity",data["identity"] )
	if data["identity"] != nil {
		return int((data["identity"]).(float64))
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId 缺少 identity")
	return 0
}

func ExtractClaims(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get(jwt.JwtPayloadKey)
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}
func GetUserDto(orm *gorm.DB, c *gin.Context) (row *sys.SysUser, err error) {
	data := ExtractClaims(c)
	userId := 0
	if data["identity"] != nil {
		userId = int((data["identity"]).(float64))
	}
	if userId == 0 {
		return nil, errors.New("用户不存在")
	}
	var user *sys.SysUser
	orm.Model(&sys.SysUser{}).Select("user_id,c_id,role_id,username,phone").Where("enable = ? and user_id = ? ", true, userId).Limit(1).Find(&user)
	if user.UserId == 0 {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

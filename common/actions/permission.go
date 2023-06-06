package actions

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
	"time"
)

type DataPermission struct {
	DataScope int
	UserId    int
	DeptId    int
	RoleId    int
	Enable    bool
	CId       int
}

func PermissionSuperRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := user.GetUserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				response.Error(c, 500, err, "权限范围鉴定错误")
				c.Abort()
				return
			}
		}
		if !p.Enable {
			response.Error(c, 401, errors.New("您账户已被停用！"), "您账户已被停用！")
			c.Abort()
			return
		}

		if p.RoleId == 0 {
			response.Error(c, 401, errors.New("您没有权限访问"), "您没有权限访问")
			c.Abort()
			return
		}
		//权限校验
		if p.DataScope != global.RoleSuper {
			response.Error(c, 401, errors.New("您没有权限访问"), "您没有权限访问")
			c.Abort()
			return
		}

		c.Set(PermissionKey, p)
		c.Next()
	}
}
func init() {

}

// 大B的权限
func PermissionCompanyRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := user.GetUserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				response.Error(c, 500, err, "权限范围鉴定错误")
				c.Abort()
				return
			}
		}
		if !p.Enable {
			response.Error(c, 401, errors.New("您账户已被停用！"), "您账户已被停用！")
			c.Abort()
			return
		}

		if p.RoleId == 0 {
			response.Error(c, 401, errors.New("您没有权限访问"), "您没有权限访问")
			c.Abort()
			return
		}
		//权限校验
		if p.DataScope == 0 {
			response.Error(c, 401, errors.New("您没有权限访问"), "您没有权限访问")
			c.Abort()
			return
		}
		if p.DataScope == global.RoleShop || p.DataScope == global.RoleUser {
			response.Error(c, 401, errors.New("您没有权限访问"), "您没有权限访问")
			c.Abort()
			return
		}
		//是否过期校验
		var companyObject models.Company
		if p.CId == 0 {
			response.Error(c, 500, errors.New("您暂无系统"), "您暂无系统")
			c.Abort()
			return
		}
		db.Model(&models.Company{}).Select("id,renewal_time").Where("id = ? and enable = ?", p.CId, true).First(&companyObject)
		if companyObject.Id == 0 {
			response.Error(c, 401, errors.New("您的系统已下线"), "您的系统已下线")
			c.Abort()
			return
		}
		if companyObject.RenewalTime.Before(time.Now()) {
			response.Error(c, 401, errors.New("账号已到期,请及时续费"), "账号已到期,请及时续费")
			c.Abort()
			return
		}

		c.Set(PermissionKey, p)
		c.Next()
	}
}
func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := user.GetUserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				response.Error(c, 500, err, "权限范围鉴定错误")
				c.Abort()
				return
			}
		}
		c.Set(PermissionKey, p)
		c.Next()
	}
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*DataPermission, error) {
	var err error
	p := &DataPermission{}

	err = tx.Table("sys_user").
		Select("sys_user.user_id", "sys_role.role_id", "sys_user.c_id", "sys_user.enable", "sys_role.data_scope").
		Joins("left join sys_role on sys_role.data_scope = sys_user.role_id").
		Where("sys_user.user_id = ?", userId).
		Scan(p).Error
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}
	return p, nil
}

// todo:针对DB层的权限校验
func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !config.ApplicationConfig.EnableDP {
			return db
		}
		fmt.Println("p.dataScope", p.DataScope, p.CId)

		switch p.DataScope {
		case global.RoleSuper:
			fmt.Println("RoleSuper")
			return db
		case global.RoleCompany:

			return db.Where(tableName+".c_id = ?", p.CId)
		case global.RoleCompanyUser:

			return db.Where(tableName+".c_id = ?", p.CId)

		default:
			fmt.Println("default")
			return db
		}

	}
}

func getPermissionFromContext(c *gin.Context) *DataPermission {
	p := new(DataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch pm.(type) {
		case *DataPermission:
			p = pm.(*DataPermission)
		}
	}
	return p
}

// GetPermissionFromContext 提供非action写法数据范围约束
func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return getPermissionFromContext(c)
}

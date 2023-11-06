package handler

import (
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
	"time"
)

type Login struct {
	Phone    string `form:"phone" json:"phone" `
	UserName string `form:"username" json:"username"`
	Password string `form:"password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func LoginValidCompany(userId int,tx *gorm.DB) error {

	var companyObject models.Company
	tx.Model(&models.Company{}).Select("id,expiration_time").Where("leader_id = ? and enable = ?", userId, true).Limit(1).Find(&companyObject)
	if companyObject.Id == 0 {


		return errors.New("您的系统已下线")
	}
	if companyObject.ExpirationTime.Before(time.Now()) {

		return errors.New("账号已到期,请及时续费")
	}
	return nil
}

func (u *Login) GetUserPhone(tx *gorm.DB) (user SysUser, role SysRole, err error) {
	err = tx.Table("sys_user").Where("phone = ?  and status = ? and enable = ?",
		u.Phone, global.SysUserSuccess, true).Limit(1).First(&user).Error
	if err != nil {
		err = errors.New("手机号不存在")
		return
	}
	_, loginErr := pkg.CompareHashAndPassword(user.Password, u.Password)
	if loginErr != nil {
		log.Errorf("user login error, %s", loginErr.Error())
		err = errors.New("手机号或者密码错误")
		return
	}

	err =LoginValidCompany(user.UserId,tx)
	return
}
func (u *Login) GetUser(tx *gorm.DB) (user SysUser, role SysRole, err error) {
	err = tx.Table("sys_user").Where("username = ?  and status = '2'", u.UserName).First(&user).Error
	if err != nil {
		log.Errorf("get user error, %s", err.Error())
		return
	}
	_, err = pkg.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		err = errors.New("用户名或密码错误")
		return
	}

	return
}

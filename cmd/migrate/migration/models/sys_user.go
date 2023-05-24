package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId   int    `gorm:"primaryKey;autoIncrement;comment:编码"  json:"userId"`
	Username string `json:"username" gorm:"type:varchar(20);comment:用户名"`
	Password string `json:"-" gorm:"type:varchar(66);comment:密码"`
	Phone    string `json:"phone" gorm:"type:varchar(11);comment:手机号"`
	RoleId   int    `json:"roleId" gorm:"type:bigint;comment:角色ID"`
	Avatar   string `json:"avatar" gorm:"type:varchar(60);comment:头像"`
	Sex      string `json:"sex" gorm:"type:varchar(10);comment:性别"`
	Email    string `json:"email" gorm:"type:varchar(30);comment:邮箱"`
	DeptId   int    `json:"deptId" gorm:"type:bigint;comment:部门"`
	PostId   int    `json:"postId" gorm:"type:bigint;comment:岗位"`
	Remark   string `json:"remark" gorm:"type:varchar(50);comment:备注"`
	Status   string `json:"status" gorm:"type:varchar(4);comment:状态"`
	UnionId  string `json:"union" gorm:"size:30;"`     //微信唯一的ID
	OOpenId  string `json:"o_open_id" gorm:"size:30;"` //微信公众号的openid
	ControlBy
	ModelTime
}

func (SysUser) TableName() string {
	return "sys_user"
}

// Encrypt 加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

func (e *SysUser) BeforeCreate(_ *gorm.DB) error {
	return e.Encrypt()
}

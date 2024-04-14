package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type SysUser struct {

	UserId         int    `gorm:"primaryKey;autoIncrement;comment:编码"  json:"userId"`
	Layer          int    `gorm:"default:1;index;comment:排序"` //排序
	NickName       string `json:"nickname" gorm:"type:varchar(20);comment:昵称"`
	Username       string `json:"username" gorm:"type:varchar(20);comment:用户名"`
	Password       string `json:"-" gorm:"type:varchar(66);comment:密码"`
	Phone          string `json:"phone" gorm:"type:varchar(11);comment:手机号"`
	CId            int    `json:"c_id" gorm:"comment:关联大B"`
	Enable         bool   `gorm:"comment:是否开启"`
	RoleId         int    `json:"roleId" gorm:"type:bigint;comment:系统角色ID"`
	Avatar         string `json:"avatar" gorm:"type:varchar(60);comment:头像"`
	Sex            string `json:"sex" gorm:"type:varchar(10);comment:性别"`
	Email          string `json:"email" gorm:"type:varchar(30);comment:邮箱"`
	DeptId         int    `json:"deptId" gorm:"type:bigint;comment:部门"`
	PostId         int    `json:"postId" gorm:"type:bigint;comment:岗位"`
	Remark         string `json:"remark" gorm:"type:varchar(50);comment:备注"`
	Status         string `json:"status" gorm:"type:varchar(4);default:2;comment:状态"`
	UnionId        string `json:"union" gorm:"size:30;"`     //微信唯一的ID
	OffOpenId        string `json:"off_open_id" gorm:"size:30;"` //微信公众号的openid
	InvitationCode string `json:"invitationCode" gorm:"type:varchar(10);comment:本人邀请码"`
	AuthExamine bool `json:"auth_examine" gorm:"type:varchar(10);comment:是否拥有订单审批/操作权限"`

	models.ControlBy
	models.ModelTime
}

func (SysUser) TableName() string {
	return "sys_user"
}

func (e *SysUser) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysUser) GetId() interface{} {
	return e.UserId
}

// 加密
func (e *SysUser) Encrypt() (err error) {
	return nil
	//if e.Password == "" {
	//	return
	//}
	//
	//var hash []byte
	//
	////同时生成他的邀请码
	////e.InvitationCode = utils.GenValidateCode(6)
	//if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
	//	return
	//} else {
	//	e.Password = string(hash)
	//	return
	//}
}

func (e *SysUser) BeforeCreate(_ *gorm.DB) error {
	return e.Encrypt()
}

func (e *SysUser) AfterFind(_ *gorm.DB) error {

	return nil
}

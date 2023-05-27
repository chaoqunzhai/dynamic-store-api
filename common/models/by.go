package models

import (
	"time"

	"gorm.io/gorm"
)

type ControlBy struct {
	CreateBy int `json:"create_by" gorm:"index;comment:创建者"`
	UpdateBy int `json:"update_by" gorm:"index;comment:更新者"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy

}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
}

type ModelTime struct {
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

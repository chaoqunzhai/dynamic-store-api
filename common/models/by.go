package models

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"time"
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
	CreatedAt XTime `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt XTime `json:"updated_at" gorm:"comment:最后更新时间"`
	//保留软删除使用gorm.DeletedAt字段, 如果不保留软删除,使用其他time字段
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// 1. 创建 time.Time 类型的副本 XTime；
type XTime struct {
	time.Time
}

// 2. 为 Xtime 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
func (t XTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(output), nil

}

// 3. 为 Xtime 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t XTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// 4. 为 Xtime 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *XTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = XTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

package models


type ChinaData struct {
	Model
	Level    int    `json:"level" gorm:"column:level"`
	Latitude    string    `json:"latitude" gorm:"size:30;column:latitude"`
	Name    string    `json:"name" gorm:"size:30;column:name"`
	Pid    int    `json:"pid" gorm:"column:pid"`
	Sort    int    `json:"sort" gorm:"column:sort"`
	Shortname    string    `json:"shortname" gorm:"size:30;column:shortname"`
	Longitude    string    `json:"longitude" gorm:"size:30;column:longitude"`
	Status    int    `json:"status" gorm:"column:status"`
}
func (ChinaData) TableName() string {
	return "china_data"
}
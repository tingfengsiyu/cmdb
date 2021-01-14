package model

import (
	"gorm.io/gorm"
	"time"
)

type Switch struct {
	ID             uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name           string    `gorm:"type:varchar(30);not null" json:"name"`
	Models          string    `gorm:"type:varchar(30);not null" json:"model"`
	Location       string    `gorm:"type:varchar(30);not null" json:"location"`
	Ipaddress      string    `gorm:"type:varchar(30);not null" json:"ipaddress"`
	Create_time    time.Time `gorm:"type:datetime;" json:"create_time"`
	Update_time    time.Time `gorm:"type:datetime;" json:"update_time"`
	Cabinet_number string    `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Idc            string    `gorm:"type:varchar(30);not null" json:"idc"`
	gorm.Model
}

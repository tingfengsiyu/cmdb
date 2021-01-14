package model

import (
	"cmdb/utils/errmsg"
	"fmt"
	"gorm.io/gorm"
)

type Idc struct {
	ID             int      `gorm:"primary_key;auto_increment" json:"id"`
	City           string    `gorm:"type:varchar(30);not null" json:"city"`
	Name           string    `gorm:"type:varchar(30);not null" json:"name"`
	Cabinet           string    `gorm:"type:varchar(30);not null" json:"cabinet"`
    gorm.Model
	Cabinet_number string    `gorm:"type:varchar(30);not null" json:"cabinet_number"`
}
// 查询服务器是否存在
func CheckIdc(name string) (code int) {
	var sv Server
	db.Select("id").Where("name = ?", name).First(&sv)
	if sv.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST //2001
	}
	return errmsg.SUCCSE
}

func CreateIdc(data *Idc) int {
	err := db.Create(&data).Error
	fmt.Println(data)
	if err != nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
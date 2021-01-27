package model

import (
	"cmdb/utils/errmsg"
	"fmt"
	"gorm.io/gorm"
)

type Idc struct {
	ID      int    `gorm:"primary_key;auto_increment" json:"id"`
	City    string `gorm:"type:varchar(30);not null" json:"city"`
	Name    string `gorm:"type:varchar(30);not null" json:"name"`
	Cabinet string `gorm:"type:varchar(30);not null" json:"cabinet"`
	gorm.Model
	Cabinet_number string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
}

// 查询IDC是否存在
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
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func EditIdc(id int, data *Idc) int {
	var idc Idc
	var maps = make(map[string]interface{})
	maps["name"] = data.Name
	maps["city"] = data.City
	maps["cabinet"] = data.Cabinet
	maps["cabinet_number"] = data.Cabinet_number
	fmt.Println(maps)
	err = db.Model(&idc).Where("id=?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

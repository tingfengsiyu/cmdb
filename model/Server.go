package model

import (
	"cmdb/utils/errmsg"
	"fmt"

	"gorm.io/gorm"
)

type Server struct {
	ID        uint   `gorm:"primary_key;auto_increment;int" json:"id"`
	Name      string `gorm:"type:varchar(30);not null" json:"name"`
	Models    string `gorm:"type:varchar(30);not null" json:"models"`
	Location  string `gorm:"type:varchar(30);not null" json:"location"`
	Ipaddress string `gorm:"type:varchar(30);not null" json:"ipaddress"`
	Label     string `gorm:"type:varchar(30);not null" json:"label"`
	Cpu       string `gorm:"type:varchar(30);not null" json:"cpu"`
	Memory    string `gorm:"type:varchar(30);not null" json:"memory"`
	Disk      string `gorm:"type:varchar(30);not null" json:"disk"`
	gorm.Model
	Cabinet_number string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Idc            string `gorm:"type:varchar(30);not null" json:"idc"`
}

func CreateServer(data *Server) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 查询服务器是否存在
func CheckServer(name string) (code int) {
	var svc Server
	db.Select("id").Where("name = ?", name).First(&svc)
	if svc.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST //2001
	}
	return errmsg.SUCCSE
}

func GetServer(pageSize int, pageNum int) ([]Server, int64) {
	var svc []Server
	var total int64
	err = db.Unscoped().Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}

func EditServer(id int, data *Server) int {
	var servers Server
	var maps = make(map[string]interface{})
	maps["name"] = data.Name
	maps["models"] = data.Models
	maps["location"] = data.Location
	maps["ipaddress"] = data.Ipaddress
	maps["label"] = data.Label
	maps["cpu"] = data.Cpu
	maps["memory"] = data.Memory
	maps["disk"] = data.Disk
	maps["cabinet_number"] = data.Cabinet_number
	maps["idc"] = data.Idc
	//fmt.Println(maps)
	err = db.Model(&servers).Where("id=?", id).Updates(maps).Error
	if err != nil {
		//fmt.Println(err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteServer(id int) int {
	var servers Server
	err = db.Debug().Unscoped().Where("id = ? ", id).Delete(&servers).Error
	if err != nil {
		fmt.Println(err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

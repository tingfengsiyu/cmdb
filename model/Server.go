package model

import (
	"cmdb/utils/errmsg"
	"gorm.io/gorm"
)

type  Server struct {
	ID             uint      `gorm:"primary_key;auto_increment;int" json:"id"`
	Name           string    `gorm:"type:varchar(30);not null" json:"name"`
	Models          string    `gorm:"type:varchar(30);not null" json:"model"`
	Location       string    `gorm:"type:varchar(30);not null" json:"location"`
	Ipaddress      string    `gorm:"type:varchar(30);not null" json:"ipaddress"`
	Label      string    `gorm:"type:varchar(30);not null" json:"label"`
	Cpu      string    `gorm:"type:varchar(30);not null" json:"cpu"`
	Memory      string    `gorm:"type:varchar(30);not null" json:"memory"`
	Disk      string    `gorm:"type:varchar(30);not null" json:"disk"`
	gorm.Model
	Cabinet_number string    `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Idc            string    `gorm:"type:varchar(30);not null" json:"idc"`
}

func CreateServer(data *Server) int {
	err := db.Create(&data).Error
	if err != nil{
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
	err = db.Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}

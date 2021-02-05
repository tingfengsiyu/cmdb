package model

import (
	"cmdb/utils/errmsg"
	"fmt"
	"gorm.io/gorm"
)

type Idc struct {
	City    string `gorm:"type:varchar(30);not null" json:"city"`
	Name    string `gorm:"type:varchar(30);not null" json:"name"`
	Cabinet string `gorm:"type:varchar(30);not null" json:"cabinet"`
	gorm.Model
	Cabinet_number string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
}

// 查询IDC是否存在
func CheckIdc(name string) (code int) {
	var idc Server
	db.Select("id").Where("name = ?", name).First(&idc)
	if idc.ID > 0 {
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

func GetIDC(name string ) int64{
	var idc Server
	db.Select("id").Where("name = ?", name).First(&idc)
	if idc.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST //2001
	}
	return errmsg.SUCCSE
}

func GetIDCs(pageSize int, pageNum int) ([]Idc, int64) {
	var svc []Idc
	var total int64
	err = db.Unscoped().Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}
//网络拓扑展示
//根据机房名查询name和机柜号查询对应服务器名中对应的idc和机柜号和所属用户，形成网络拓扑
func Network_topology(id int, name,cabinet_number,user string ) int {
	//
	var idc Idc
	db.Where("id =?",id).First(Idc{})
	// 多连接及参数
	db.Model(&Server{}).Select("idc.name,idc.cabinet_number,server.name,user,location,label,cpu,memory,disk,state").Joins("left join server on server.idc = idc.name").Scan(&idc)
	db.Joins("JOIN server ON server.idc = idc.name AND server.cabinet_number = ?",cabinet_number ).Find(&idc)
	//select server.name,idc.name,server.user from server,idc where server.idc=idc.name and idc.cabinet_number=server.cabinet_number and server.user="汪洋";
	return errmsg.SUCCSE
}

func DeleteIDC(id int) int {
	var idc Idc
	err = db.Debug().Unscoped().Where("id = ? ", id).Delete(&idc).Error
	if err != nil {
		fmt.Println(err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

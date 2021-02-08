package model

import (
	"cmdb/middleware"
	"cmdb/utils/errmsg"
	"fmt"

	"gorm.io/gorm"
)

func CreateServer(data *Server) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
func LastCabintID() int {
	var data = Cabinet{}
	db.Last(&data)
	return int(data.ID)
}

func LastIdcID() int {
	var data = Idc{}
	db.Last(&data)
	return int(data.ID)
}
func LastServeID() int {
	var data = Idc{}
	db.Last(&data)
	return int(data.ID)
}

func (servers *Servers) BatchCreateServer() int {
	//func BatchCreateServer(servers *Servers)  int {
	//fmt.Println(servers.Servers)
	err := db.Create(&servers.Servers).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//批量创建
func BatchCreateServer2(servers *[]Server) int {
	err := db.Create(&servers).Error
	if err != nil {
		middleware.SugarLogger.Errorf("创建错误%s", err)
		fmt.Println(err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func BatchUpdateServer(servers *[]Server) int {
	err := db.Debug().Model(Server{}).Updates(servers).Error
	if err != nil {
		fmt.Println(err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 查询服务器是否存在
func CheckServer(name string) (int, int) {
	var svc Server
	db.Select("id").Where("name = ?", name).First(&svc)
	if svc.ID > 0 {
		return svc.ServerID, errmsg.ERROR_DEVICE_EXIST //2001
	}
	return svc.ServerID, errmsg.SUCCSE
}

//批量检查服务器名是否存在BatchCheckServer(data []model.Names)
func BatchCheckServer(data []string) (code int) {
	var svc Server
	db.Where("name IN ?", data).Find(&svc)
	if svc.ID > len(data) {
		return errmsg.ERROR_ALL_DEVICE_EXIST
	} else if svc.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST
	}
	return errmsg.SUCCSE
}

func BatchCheckServerID(data []int) (code int) {
	var svc Server
	db.Find(&svc, data)
	if svc.ID >= len(data) {
		fmt.Println(svc.ID)
		return errmsg.ERROR_ALL_DEVICE_EXIST
	} else if svc.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST
	}
	return errmsg.SUCCSE
}

func GetServers(pageSize int, pageNum int) ([]Server, int64) {
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
	maps["private_ip_address"] = data.PrivateIpAddress
	maps["public_ip_address"] = data.PublicIpAddress
	maps["label"] = data.Label
	maps["cpu"] = data.Cpu
	maps["memory"] = data.Memory
	maps["disk"] = data.Disk
	maps["user"] = data.User
	maps["state"] = data.State
	//fmt.Println(maps)
	err = db.Model(&servers).Where("id=?", id).Updates(maps).Error
	if err != nil {
		middleware.SugarLogger.Errorf("插入错误%s", err)
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

func GetServerInfo(id int) ([]Server, int) {
	var svc []Server
	var total int
	err := db.Where("ID = ?", id).First(&svc).Error
	if err != nil {
		return nil, 0
	}
	return svc, total
}

//查询所有客户
func GetOwnedUser() ([]Server, int64) {
	var svc []Server
	var total int64
	//err :=  db.Find(&svc).Distinct("user").Error
	db.Distinct("user").Find(&svc)
	db.Model(&svc).Count(&total)
	if err != nil {
		return nil, 0
	}
	return svc, total
}

//查询对应城市的所有服务器
func GetIdcServers(pageSize, pageNum int, name string) ([]Server, int64) {
	var svc []Server
	var total int64
	err = db.Unscoped().Where(" idc = ?", name).Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}

//查询对应城市所对应机柜的所有服务器
func GetCabinetServers(pageSize, pageNum int, name, cabinet_number string) ([]Server, int64) {
	var svc []Server
	var total int64
	err = db.Unscoped().Where("idc = ? AND  cabinet_number = ?", name, cabinet_number).Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}

func UpdateID(idc_name, city_name, cabinet_number, name string, idc_id, server_id, cabinet_number_id int) int {
	var idc = Idc{}
	var cabinet = Cabinet{}
	var server = Server{}

	var servers = make(map[string]interface{})
	servers["server_id"] = server_id
	servers["cabinet_number_id"] = cabinet_number_id
	servers["idc_id"] = idc_id

	var idcs = make(map[string]interface{})
	idcs["cabinet_number_id"] = cabinet_number_id
	idcs["idc_id"] = idc_id
	idcs["city"] = city_name

	var cabinets = make(map[string]interface{})
	cabinets["cabinet_number_id"] = cabinet_number_id
	cabinets["idc_id"] = idc_id

	db.Model(&idc).Where("idc_name =?", idc_name).Updates(idcs)
	db.Model(&cabinet).Where("cabinet_number =?", cabinet_number).Updates(cabinets)
	db.Model(&server).Where("name =?", name).Updates(servers)
	return errmsg.SUCCSE
}

func InsertID(idc_name, city_name, cabinet_number, name string, idc_id, server_id, cabinet_number_id int) {
	var idc = Idc{}
	var cabinet = Cabinet{}
	var server = Server{}

	var servers = make(map[string]interface{})
	servers["server_id"] = server_id
	servers["cabinet_number_id"] = cabinet_number_id
	servers["idc_id"] = idc_id

	var idcs = make(map[string]interface{})
	idcs["cabinet_number_id"] = cabinet_number_id
	idcs["idc_id"] = idc_id
	idcs["city"] = city_name
	idcs["idc_name"] = idc_name

	var cabinets = make(map[string]interface{})
	cabinets["cabinet_number_id"] = cabinet_number_id
	cabinets["idc_id"] = idc_id
	cabinets["cabinet_number"] = cabinet_number

	db.Model(&idc).Create(idcs)
	db.Model(&cabinet).Create(cabinets)
	db.Model(&server).Where("name =?", name).Updates(servers)
}

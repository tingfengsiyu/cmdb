package model

import (
	"cmdb/middleware"
	"cmdb/utils/errmsg"
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
	db.Order("cabinet_number_id desc").Find(&data).Limit(1)
	return int(data.Cabinet_NumberID)
}

func LastIdcID() int {
	var data = Idc{}
	db.Order("idc_id desc").Find(&data).Limit(1)
	return int(data.IDC_ID)
}
func LastServeID() int {
	var data = Server{}
	db.Unscoped().Debug().Order("server_id desc").Find(&data).Limit(1)
	return int(data.ServerID)
}

//func (servers *Servers) BatchCreateServer() int {
//	err := db.Create(&servers.Servers).Error
//	if err != nil {
//		return errmsg.ERROR
//	}
//	return errmsg.SUCCSE
//}

type ScanServers struct {
	ID               int    `gorm:"primary_key;auto_increment;int" json:"id"`
	City             string `gorm:"type:varchar(30);not null" json:"city" validate:"required,min=4"`
	IDC_Name         string `gorm:"type:varchar(30);not null" json:"idc_name" validate:"required,min=4"`
	Cabinet_Number   string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Name             string `gorm:"type:varchar(30);not null" json:"name" validate:"required,min=4"`
	Models           string `gorm:"type:varchar(30);not null" json:"models" validate:"required,min=4"`
	Location         string `gorm:"type:varchar(30);not null" json:"location" validate:"required,min=4"`
	PrivateIpAddress string `gorm:"type:varchar(30);not null" json:"private_ip_address" validate:"required,min=16"`
	PublicIpAddress  string `gorm:"type:varchar(30);not null" json:"public_ip_address" `
	Label            string `gorm:"type:varchar(30);not null" json:"label" validate:"required,min=4"`
	Cluster          string `gorm:"type:varchar(30);not null" json:"cluster" validate:"required,min=4"`
	LabelIpAddress   string `gorm:"type:varchar(30);not null" json:"label_ip_address" validate:"required,min=4"`
	Cpu              string `gorm:"type:varchar(30);not null" json:"cpu" validate:"required,min=3"`
	Memory           string `gorm:"type:varchar(30);not null" json:"memory" validate:"required,min=3"`
	Disk             string `gorm:"type:varchar(30);not null" json:"disk" validate:"required,min=3"`
	User             string `gorm:"type:varchar(30);not null" json:"user" validate:"required,min=4"`
	State            string `gorm:"type:varchar(10);not null" json:"state" validate:"required,min=4"`
	IDC_ID           int    `gorm:"type:int;not null" json:"idc_id" validate:"required,min=4"`
	Cabinet_NumberID int    `gorm:"type:int;not null" json:"cabinet_number_id" validate:"required,min=4"`
}

//批量创建
func BatchCreateServer(servers *[]Server) int {
	err := db.Create(&servers).Error
	if err != nil {
		middleware.SugarLogger.Errorf("创建错误%s", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func BatchUpdateServer(servers map[string]interface{}, ID int) int {
	err := db.Debug().Model(Server{}).Where("id=?", ID).Updates(servers).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 查询服务器是否存在
func CheckServer(name string) (int, int) {
	var svc Server
	db.Select("id").Where("name = ?", name).First(&svc)
	if svc.ID > 0 {
		return svc.ID, errmsg.ERROR_DEVICE_EXIST //2001
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
	//select name,models,idc_name,city,cabinet.cabinet_number from  server  left join cabinet on  cabinet.cabinet_number_id=server.cabinet_number_id left join idc on idc.idc_id =server.idc_id;
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
	maps["label_ip_address"] = data.LabelIpAddress
	maps["label"] = data.Label
	maps["cluster"] = data.Cluster
	maps["cpu"] = data.Cpu
	maps["memory"] = data.Memory
	maps["disk"] = data.Disk
	maps["user"] = data.User
	maps["state"] = data.State
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
		middleware.SugarLogger.Errorf("删除错误%s", err)
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

func InsertServerID(name string, idc_id, server_id, cabinet_number_id int) {
	var server = Server{}
	var servers = make(map[string]interface{})
	servers["server_id"] = server_id
	servers["cabinet_number_id"] = cabinet_number_id
	servers["idc_id"] = idc_id
	db.Model(&server).Where("name =?", name).Updates(servers)
}

func GenerateServerID(hostNames []string) []int {
	//server_id
	var server_ids = make([]int, 0)
	number := 0
	for k, v := range hostNames {
		server_id, _ := CheckServer(v)
		if server_id == 0 {
			if k == 0 {
				id := LastServeID()
				number = id + 1
			} else {
				number = number + 1
			}
			server_ids = append(server_ids, number)
		} else {
			server_ids = append(server_ids, server_id)
		}
	}
	return server_ids
}

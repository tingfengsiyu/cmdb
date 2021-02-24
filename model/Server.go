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

func (servers *Servers) BatchCreateServer() int {
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
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func BatchUpdateServer(servers *[]Server) int {
	err := db.Debug().Model(Server{}).Updates(servers).Error
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

	//检查 cabinet_number_id 是否已存在，存在相同则不创建
	err = db.Unscoped().Debug().Where("idc_name = ? AND  cabinet_number_id = ?", idc_name, cabinet_number_id).Find(&idc).Error
	if err != nil {
		middleware.SugarLogger.Errorf("查询idc错误%s", err)
	}
	if idc.Cabinet_NumberID == 0 {
		db.Model(&idc).Create(idcs)
	}

	err = db.Where("idc_id = ? AND  cabinet_number_id = ?", idc_id, cabinet_number_id).Find(&cabinet).Error
	if err != nil {
		middleware.SugarLogger.Errorf("查询cabinet错误%s", err)
	}
	if cabinet.Cabinet_NumberID == 0 {
		db.Model(&cabinet).Create(cabinets)
	}
	db.Model(&server).Where("name =?", name).Updates(servers)

	//更新数据到监控表,  监控状态 0 未安装 ；1 已安装 ；2 已运行；3 已被监控
	var prometheustarget = MonitorPrometheus{}
	var prometheus = make(map[string]interface{})
	prometheus["server_id"] = server_id
	prometheus["node_export_port"] = 9100
	prometheus["process_export_port"] = 9256
	prometheus["script_export_port"] = 9172
	prometheus["node_export_status"] = 0
	prometheus["process_export_status"] = 0
	prometheus["script_export_status"] = 0

	db.Model(&prometheustarget).Create(prometheus)
}

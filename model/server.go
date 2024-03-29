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
	db.Order("server_id desc").Find(&data).Limit(1)
	return int(data.ServerID)
}

//func (servers *Servers) BatchCreateServer() int {
//	err := db.Create(&servers.Servers).Error
//	if err != nil {
//		return errmsg.ERROR
//	}
//	return errmsg.SUCCSE
//}

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
	var svc []Server
	db.Debug().Find(&svc, data)
	if len(svc) >= len(data) {
		return errmsg.ERROR_ALL_DEVICE_EXIST
	} else if len(svc) > 0 {
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

func GetClusters() ([]Cluster, int) {
	var cluster []Cluster
	var svc []Server
	//err = db.Distinct("cluster,count(cluster) as count").Group("cluster").Find(&svc).Error
	err = db.Model(&Server{}).Select("cluster,count(cluster) as count").Group("cluster").Find(&cluster).Error
	for _, v := range svc {
		cluster = append(cluster, Cluster{
			ID:      v.ID,
			Cluster: v.Cluster,
		})
	}

	if err != nil {
		return nil, 0
	}
	return cluster, len(cluster)
}
func GetCluster(cluster string) ([]Server, int64) {
	var svc []Server
	var total int64
	err = db.Where("cluster = ?", cluster).Find(&svc).Error
	db.Model(&svc).Count(&total)
	if err != nil {
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
	err = db.Unscoped().Where("id = ? ", id).Delete(&servers).Error
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
func GetIdcServers(pageSize, pageNum int, name string) ([]ScanServers, int64) {
	var scan []ScanServers
	var total int64
	err = db.Model(&Server{}).Select("select distinct server.id, city,idc_name,cabinet_number,name,"+
		"models,location,private_ip_address,public_ip_address,label,cluster,label_ip_address,cpu,memory,disk,user,state,server.idc_id,server.cabinet_number_id").Joins("left join idc on "+
		"idc.idc_id =server.idc_id and idc_name= ?", name).Scan(&scan).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return scan, total
}

//查询对应城市所对应机柜的所有服务器
func GetCabinetServers(pageSize, pageNum int, name, cabinetNumber string) ([]ScanServers, int64) {
	var total int64
	var scan []ScanServers
	err := db.Raw("select distinct server.id, city,idc_name,cabinet_number,name,models,location,private_ip_address,public_ip_address,label,cluster,label_ip_address,cpu,memory,disk,user,state,"+
		"server.idc_id,server.cabinet_number_id from  server  left join cabinet  on  cabinet.cabinet_number_id=server.cabinet_number_id  and idc_name = ? left join idc on idc.idc_id =server.idc_id and cabinet_number=? "+
		"limit ?,?", name, cabinetNumber, pageSize, pageNum).Scan(&scan).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return scan, total
}

func InsertServerID(name string, idcId, serverId, cabinetNumberId int) {
	var server = Server{}
	var servers = make(map[string]interface{})
	servers["server_id"] = serverId
	servers["cabinet_number_id"] = cabinetNumberId
	servers["idc_id"] = idcId
	db.Model(&server).Where("name =?", name).Updates(servers)
}

func GenerateServerID(hostNames []string) []int {
	//server_id
	var serverIds = make([]int, 0)
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
			serverIds = append(serverIds, number)
		} else {
			serverIds = append(serverIds, server_id)
		}
	}
	return serverIds
}

func CheckClusterIp(ipaddress string) int {
	var svc Server
	if err := db.Select("id").Where("private_ip_address = ?", ipaddress).First(&svc).Error; err != nil {
		middleware.SugarLogger.Error(err)
	}
	return svc.ID
}
func CheckClusterName(cluster string) int {
	var svc Server
	db.Select("id").Where("cluster = ?", cluster).First(&svc)
	return svc.ID
}

func BatchCheckClusterIps(ipaddress []string) bool {
	var svc []Server
	db.Select("private_ip_address").Where("private_ip_address IN ?", ipaddress).Find(&svc)
	if len(svc) != len(ipaddress) {
		return false
	}
	return true
}

func UpdateClusterName(id int, ipaddress, clustername string) {
	err := db.Model(Server{}).Where("id = ?", id).Updates(Server{PrivateIpAddress: ipaddress, Cluster: clustername})
	if err != nil {
		middleware.SugarLogger.Error(err)
	}
}

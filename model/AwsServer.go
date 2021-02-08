package model

import (
	"cmdb/utils/errmsg"

	"gorm.io/gorm"
)

func GetAws(pageSize int, pageNum int) ([]AwsServer, int64) {
	var svc []AwsServer
	var total int64
	err = db.Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}

func EditAws(id int, data *AwsServer) int {
	var servers AwsServer
	var maps = make(map[string]interface{})
	maps["name"] = servers.InstanceId
	maps["model"] = servers.Name
	maps["location"] = servers.LaunchTime
	maps["ipaddress"] = servers.ExpiredTime
	maps["label"] = servers.PublicIpAddress
	maps["cpu"] = servers.PrivateIpAddress
	maps["memory"] = servers.InstanceType
	maps["disk"] = servers.State
	maps["cabinet_number"] = servers.InstanceId
	maps["idc"] = servers.InstanceId

	//InstanceId,Name,CreateTime,ExpiredTime,PublicIpAddress,PrivateIpAddress,OsName,InstanceType,State,Regions,AvailabilityZones,Disk,Arch,User,Cloud
	err = db.Model(&servers).Where("id=?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

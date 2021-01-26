package model

import (
	"cmdb/utils/errmsg"
	"gorm.io/gorm"
	"time"
)

type AwsServer struct {
	gorm.Model
	InstanceId       string    `gorm:"type:varchar(30);not null" json:"instanceid"`
	Name             string    `gorm:"type:varchar(30);not null" json:"label"`
	LaunchTime       time.Time `gorm:"type:datetime;" json:"launchtime"`
	ExpiredTime      time.Time `gorm:"type:datetime;" json:"expiredtime"`
	PublicIpAddress  string    `gorm:"type:varchar(30);not null" json:"publicip"`
	PrivateIpAddress string    `gorm:"type:varchar(30);not null" json:"privateip"`
	ImageName        string    `gorm:"type:varchar(30);not null" json:"imagename"`
	InstanceType     string    `gorm:"type:varchar(30);not null" json:"instancetype"`
	State            string    `gorm:"type:varchar(30);not null" json:"state"`
	Region           string    `gorm:"type:varchar(30);not null" json:"region"`
	AvailabilityZone string    `gorm:"type:varchar(30);not null" json:"availabilityZone"`
	Disk             string    `gorm:"type:varchar(30);not null" json:"disk"`
	Architecture     string    `gorm:"type:varchar(30);not null" json:"architecture"`
	Owner            string    `gorm:"type:varchar(30);not null" json:"owner"`
	Cloud            string    `gorm:"type:varchar(30);not null" json:"cloud"`
}

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

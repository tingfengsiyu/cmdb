package model

import (
	"cmdb/middleware"
	"cmdb/utils/errmsg"
	"gorm.io/gorm"
	"time"
)

type CloudInstance struct {
	gorm.Model
	InstanceId             string `gorm:"type:varchar(30);not null" json:"InstanceId" xml:"InstanceId"`
	HostName               string `gorm:"type:varchar(30);not null" json:"HostName" xml:"HostName"`
	Status                 string `gorm:"type:varchar(30);not null" json:"Status" xml:"Status"`
	CPU                    int    `gorm:"type:int;not null" json:"CPU" xml:"CPU"`
	Memory                 int    `gorm:"type:int;not null" json:"Memory" xml:"Memory"`
	OSName                 string `gorm:"type:varchar(50);not null" json:"OSName" xml:"OSName"`
	RegionId               string `gorm:"type:varchar(30);not null" json:"RegionId" xml:"RegionId"`
	InstanceType           string `gorm:"type:varchar(30);not null" json:"InstanceType" xml:"InstanceType"`
	OsType                 string `gorm:"type:varchar(30);not null" json:"OsType" xml:"OsType"`
	InternetMaxBandwidthIn int    `gorm:"type:varchar(30);not null" json:"InternetMaxBandwidthIn" xml:"InternetMaxBandwidthIn"`
	StartTime              time.Time  `gorm:"type:datetime(3);null" json:"StartTime" xml:"StartTime"`
	ExpiredTime            time.Time `gorm:"type:datetime(3);null" json:"ExpiredTime" xml:"ExpiredTime"`
	InstanceCreationTime   time.Time `gorm:"type:datetime(3);null" json:"InstanceCreationTime" xml:"InstanceCreationTime"`
	LocalStorageCapacity   int64  `gorm:"type:varchar(30);not null" json:"LocalStorageCapacity" xml:"LocalStorageCapacity"`
	InnerIpAddress         string `gorm:"type:varchar(30);not null" json:"InnerIpAddress" xml:"InnerIpAddress"`
	PublicIpAddress        string `gorm:"type:varchar(30);not null" json:"PublicIpAddress" xml:"PublicIpAddress"`
	Cloud string `gorm:"type:varchar(10);not null" json:"Cloud" xml:"Cloud"`
}

func BatchAddAliEcs(servers []CloudInstance) int {
	err := db.Debug().Create(&servers).Error
	if err != nil {
		middleware.SugarLogger.Debugf("批量插入错误%s", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

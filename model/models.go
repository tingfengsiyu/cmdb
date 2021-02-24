package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

//IDC机房信息表
type Idc struct {
	gorm.Model
	City             string `gorm:"type:varchar(30);not null" json:"city" validate:"required,min=4"`
	IDC_Name         string `gorm:"type:varchar(30);not null" json:"idc_name" validate:"required,min=4"`
	IDC_ID           int    `gorm:"type:int;not null" json:"idc_id" validate:"required,min=1"`
	Cabinet_NumberID int    `gorm:"type:int;not null" json:"cabinet_number_id" validate:"required,min=4"`
}

//机柜表
type Cabinet struct {
	gorm.Model
	IDC_ID           int    `gorm:"type:int;not null" json:"idc_id"`
	Cabinet_NumberID int    `gorm:"type:int;not null" json:"cabinet_number_id"`
	Cabinet_Number   string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
}

//服务器表
type Server struct {
	ID               int    `gorm:"primary_key;auto_increment;int" json:"id"`
	Name             string `gorm:"type:varchar(30);not null" json:"name" validate:"required,min=4"`
	Models           string `gorm:"type:varchar(30);not null" json:"models" validate:"required,min=4"`
	Location         string `gorm:"type:varchar(30);not null" json:"location" validate:"required,min=4"`
	PrivateIpAddress string `gorm:"type:varchar(30);not null" json:"private_ip_address" validate:"required,min=16"`
	PublicIpAddress  string `gorm:"type:varchar(30);not null" json:"public_ip_address" `
	Label            string `gorm:"type:varchar(30);not null" json:"label" validate:"required,min=4"`
	Cpu              string `gorm:"type:varchar(30);not null" json:"cpu" validate:"required,min=3"`
	Memory           string `gorm:"type:varchar(30);not null" json:"memory" validate:"required,min=3"`
	Disk             string `gorm:"type:varchar(30);not null" json:"disk" validate:"required,min=3"`
	gorm.Model
	User             string `gorm:"type:varchar(30);not null" json:"user" validate:"required,min=4"`
	State            string `gorm:"type:varchar(10);not null" json:"state" validate:"required,min=4"`
	ServerID         int    `gorm:"type:int;not null" json:"server_id" validate:"required,min=1"`
	IDC_ID           int    `gorm:"type:int;not null" json:"idc_id" validate:"required,min=4"`
	Cabinet_NumberID int    `gorm:"type:int;not null" json:"cabinet_number_id" validate:"required,min=4"`
}

//监控表
type MonitorPrometheus struct {
	gorm.Model
	ServerID            int `gorm:"unique;type:int;not null" json:"server_id" validate:"required,min=1"`
	NodeExportPort      int `gorm:"type:int;not null" json:"node_export_port" validate:"required,min=4"`
	ProcessExportPort   int `gorm:"type:int;not null" json:"process_export_port" validate:"required,min=4"`
	ScriptExportPort    int `gorm:"type:int;not null" json:"script_export_port" validate:"required,min=4"`
	NodeExportStatus    int `gorm:"type:int;not null" json:"node_export_status" validate:"required,min=1"`
	ProcessExportStatus int `gorm:"type:int;not null" json:"process_export_status" validate:"required,min=1"`
	ScriptExportStatus  int `gorm:"type:int;not null" json:"script_export_status" validate:"required,min=1"`
}

type Servers struct {
	Servers []Server `json: "servers"`
}

type Idcs struct {
	Idc_name       []string `json:"idc_name"`
	City           []string `json:"city"`
	Cabinet_Number []string `json:"cabinet_number"`
}

type Asset struct {
	Servers []Server `json: "servers"`
	Idcs    Idcs     `json: "idcs"`
}
type Assets struct {
	Asset Asset
}

/*
type Router struct {
	ID             uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name           string    `gorm:"type:varchar(30);not null" json:"name"`
	Models         string    `gorm:"type:varchar(30);not null" json:"model"`
	Location       string    `gorm:"type:varchar(30);not null" json:"location"`
	Ipaddress      string    `gorm:"type:varchar(30);not null" json:"ipaddress"`
	Create_time    time.Time `gorm:"type:datetime;" json:"create_time"`
	Update_time    time.Time `gorm:"type:datetime;" json:"update_time"`
	Cabinet_number string    `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Idc            string    `gorm:"type:varchar(30);not null" json:"idc"`
	gorm.Model
}

type Switch struct {
	Name           string `gorm:"type:varchar(30);not null" json:"name"`
	Models         string `gorm:"type:varchar(30);not null" json:"model"`
	Location       string `gorm:"type:varchar(30);not null" json:"location"`
	Ipaddress      string `gorm:"type:varchar(30);not null" json:"ipaddress"`
	Cabinet_number string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Idc            string `gorm:"type:varchar(30);not null" json:"idc"`
	gorm.Model
}
*/
type CloudInstance struct {
	gorm.Model
	InstanceId             string    `gorm:"type:varchar(30);not null" json:"InstanceId" xml:"InstanceId"`
	HostName               string    `gorm:"type:varchar(30);not null" json:"HostName" xml:"HostName"`
	Status                 string    `gorm:"type:varchar(30);not null" json:"Status" xml:"Status"`
	CPU                    int       `gorm:"type:int;not null" json:"CPU" xml:"CPU"`
	Memory                 int       `gorm:"type:int;not null" json:"Memory" xml:"Memory"`
	OSName                 string    `gorm:"type:varchar(50);not null" json:"OSName" xml:"OSName"`
	RegionId               string    `gorm:"type:varchar(30);not null" json:"RegionId" xml:"RegionId"`
	InstanceType           string    `gorm:"type:varchar(30);not null" json:"InstanceType" xml:"InstanceType"`
	OsType                 string    `gorm:"type:varchar(30);not null" json:"OsType" xml:"OsType"`
	InternetMaxBandwidthIn int       `gorm:"type:varchar(30);not null" json:"InternetMaxBandwidthIn" xml:"InternetMaxBandwidthIn"`
	StartTime              time.Time `gorm:"type:datetime(3);null" json:"StartTime" xml:"StartTime"`
	ExpiredTime            time.Time `gorm:"type:datetime(3);null" json:"ExpiredTime" xml:"ExpiredTime"`
	InstanceCreationTime   time.Time `gorm:"type:datetime(3);null" json:"InstanceCreationTime" xml:"InstanceCreationTime"`
	LocalStorageCapacity   int64     `gorm:"type:varchar(30);not null" json:"LocalStorageCapacity" xml:"LocalStorageCapacity"`
	PrivateIpAddress       string    `gorm:"type:varchar(30);not null" json:"PrivateIpAddress" xml:"PrivateIpAddress"`
	PublicIpAddress        string    `gorm:"type:varchar(30);not null" json:"PublicIpAddress" xml:"PublicIpAddress"`
	Cloud                  string    `gorm:"type:varchar(10);not null" json:"Cloud" xml:"Cloud"`
}

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

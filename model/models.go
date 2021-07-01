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
	Name             string `gorm:"type:varchar(42);not null;unique" json:"name" validate:"required,min=4"`
	Models           string `gorm:"type:varchar(30);not null" json:"models" validate:"required,min=4"`
	Location         string `gorm:"type:varchar(30);not null" json:"location" validate:"required,min=4"`
	PrivateIpAddress string `gorm:"type:varchar(30);not null;unique" json:"private_ip_address" validate:"required,min=16"`
	PublicIpAddress  string `gorm:"type:varchar(30);not null" json:"public_ip_address" `
	Label            string `gorm:"type:varchar(30);not null" json:"label" binding:"required" validate:"required,min=4"`
	Cluster          string `gorm:"type:varchar(30);not null" json:"cluster" binding:"required" validate:"required,min=4"`
	LabelIpAddress   string `gorm:"type:varchar(30);not null" json:"label_ip_address" validate:"required,min=4"`
	Cpu              string `gorm:"type:varchar(30);not null" json:"cpu" validate:"required,min=3"`
	Memory           string `gorm:"type:varchar(30);not null" json:"memory" validate:"required,min=3"`
	Disk             string `gorm:"type:varchar(30);not null" json:"disk" validate:"required,min=3"`
	gorm.Model
	User             string `gorm:"type:varchar(30);not null" json:"user" validate:"required,min=4"`
	State            string `gorm:"type:varchar(10);not null" json:"state" validate:"required,min=4"`
	ServerID         int    `gorm:"type:int;not null;unique" json:"server_id" validate:"required,min=1"`
	IDC_ID           int    `gorm:"type:int;not null" json:"idc_id" validate:"required,min=4"`
	Cabinet_NumberID int    `gorm:"type:int;not null" json:"cabinet_number_id" validate:"required,min=4"`
}

//监控表
type MonitorPrometheus struct {
	gorm.Model
	ServerID             int `gorm:"unique;type:int;not null" json:"server_id" validate:"required,min=1"`
	NodeExportPort       int `gorm:"type:int;DEFAULT:9100" json:"node_export_port" validate:"required,min=4"`
	ProcessExportPort    int `gorm:"type:int;DEFAULT:9256" json:"process_export_port" validate:"required,min=4"`
	ScriptExportPort     int `gorm:"type:int;DEFAULT:9172" json:"script_export_port" validate:"required,min=4"`
	NodeExportStatus     int `gorm:"type:int;DEFAULT:0" json:"node_export_status" validate:"required,min=1"`
	ProcessExportStatus  int `gorm:"type:int;DEFAULT:0" json:"process_export_status" validate:"required,min=1"`
	ScriptExportStatus   int `gorm:"type:int;DEFAULT:0" json:"script_export_status" validate:"required,min=1"`
	DisableNodeExport    int `gorm:"type:int;DEFAULT:1" json:"disable_node_export" `
	DisableProcessExport int `gorm:"type:int;DEFAULT:1" json:"disable_process_export"`
	DisableScriptExport  int `gorm:"type:int;DEFAULT:1" json:"disable_script_export"`

	//监控状态 0 未安装 ；1 已安装 ；2 已运行；3 已被监控,disable 0 禁用监控 1 启用监控
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
	gorm.Model
	Name           string    `gorm:"type:varchar(30);not null" json:"name"`
	Models         string    `gorm:"type:varchar(30);not null" json:"model"`
	Location       string    `gorm:"type:varchar(30);not null" json:"location"`
	Ipaddress      string    `gorm:"type:varchar(30);not null" json:"ipaddress"`
	Cabinet_number string    `gorm:"type:varchar(30);not null" json:"cabinet_number"`
	Idc            string    `gorm:"type:varchar(30);not null" json:"idc"`

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
type Cluster struct {
	ID      int    `gorm:"primary_key;auto_increment" json:"id"`
	Cluster string `gorm:"type:varchar(30);not null" json:"cluster" binding:"required" validate:"required,min=4"`
}

type ScanServers struct {
	ID               int    `gorm:"primary_key;auto_increment;int" json:"id"`
	Name             string `gorm:"type:varchar(42);not null" json:"name" validate:"required,min=4"`
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
	City             string `gorm:"type:varchar(30);not null" json:"city" validate:"required,min=4"`
	IDC_Name         string `gorm:"type:varchar(30);not null" json:"idc_name" validate:"required,min=4"`
	Cabinet_Number   string `gorm:"type:varchar(30);not null" json:"cabinet_number"`
}

type OpsRecords struct {
	gorm.Model
	User    string `gorm:"type:varchar(30);not null" json:"user"`
	Object  string `gorm:"type:varchar(30);not null" json:"object"`
	Action  string `gorm:"type:varchar(30);not null" json:"action"`
	State   int    `gorm:"type:int;not null;default:2" json:"state"`
	Success string `gorm:"type:varchar(1000);not null" json:"success"`
	Error   string `gorm:"type:varchar(1000);not null" json:"error"`
}

type BatchIpStruct struct {
	SourceStartIp     string `json:"source_start_ip" validate:"required,min=10,max=12" `
	SourceGateway     string `json:"source_gateway" validate:"required,min=10,max=10" `
	SourceEndNumber   string `json:"source_end_number" validate:"required,gte=2"  `
	TargetStartIP     string `json:"target_start_ip" validate:"required,min=10,max=12" `
	TargetGateway     string `json:"target_gateway" validate:"required,min=10,max=10" `
	TargetClusterName string `json:"target_cluster_name" validate:"required,min=4,max=50"`
}

type UpdateClusterStruct struct {
	SourceStartIp     string `json:"source_start_ip" validate:"required,min=10,max=12" `
	SourceEndNumber   string `json:"source_end_number" validate:"required,gte=2"  `
	TargetClusterName string `json:"target_cluster_name" validate:"required,min=4,max=50"`
}

type OsInitStruct struct {
	InitUser     string `json:"init_user" validate:"required,min=10,max=10" `
	InitPass     string `json:"init_pass" validate:"required,min=4,max=50"`
	Role         string `json:"role" validate:"required,min=4,max=10"`
	StorageMount StorageMountStruct
}

type StorageMountStruct struct {
	InitStartIP       string `json:"init_start_ip" validate:"required,min=10,max=12" `
	InitEndNumber     string `json:"init_end_number" validate:"required,gte=2" `
	StorageStartIP    string `json:"storage_start_ip" validate:"required,min=4,max=50"`
	StorageStopnumber string `json:"storage_stop_number" validate:"required,min=1,max=3"`
	Operating         string `json:"operating" validate:"required,min=1,max=3"`
}

type ansibleStruct struct {
	PrivateIpAddress string `json:"private_ip_address"`
	Label            string `json:"label"`
	Cluster          string `json:"cluster"`
}

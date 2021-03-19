package utils

import (
	"cmdb/middleware"
	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
	JwtKey   string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	AccessKey string
	SecretKey string

	LogFile           string
	LogLevel          string
	ErrorLogFile      string
	KubeFile          string
	PrometheusConfDir string
	WorkerUser        string
	WorkerPass        string
	WorkerSudoPass    string
	StorageUser       string
	StoragePass       string
	StorageSudoPass   string
	RootPass          string
	RootPub           string
	SR_File_Max_Bytes string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		middleware.SugarLogger.Errorf("配置文件读取错误，请检查文件路径: %s", err)
		panic(1)
	}
	LoadServer(file)
	LoadData(file)
	LoadAli(file)
	//	LoadAws(file)
	Loadk8s(file)
	LoadPrometheus(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
	LogFile = file.Section("server").Key("LogFile").MustString("xxxx/cmdb")
	ErrorLogFile = file.Section("server").Key("ErrorLogFile").MustString("xxxx/cmdb")
	WorkerUser = file.Section("server").Key("WorkerUser").MustString("xxxx/cmdb")
	WorkerPass = file.Section("server").Key("WorkerPass").MustString("xxxx/cmdb")
	WorkerSudoPass = file.Section("server").Key("WorkerSudoPass").MustString("xxxx/cmdb")
	StorageUser = file.Section("server").Key("StorageUser").MustString("xxxx/cmdb")
	StoragePass = file.Section("server").Key("StoragePass").MustString("xxxx/cmdb")
	StorageSudoPass = file.Section("server").Key("StorageSudoPass").MustString("xxxx/cmdb")
	RootPass = file.Section("server").Key("RootPass").MustString("xxxx/cmdb")
	RootPub = file.Section("server").Key("RootPub").MustString("xxxx/cmdb")
	SR_File_Max_Bytes = file.Section("server").Key("SR_File_Max_Bytes").MustString("")
}
func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("debug")
	DbHost = file.Section("database").Key("DbHost").MustString("47.s.197.46")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("caichangen")
	DbName = file.Section("database").Key("DbName").MustString("cmdb")
}

func LoadAli(file *ini.File) {
	AccessKey = file.Section("ali").Key("AccessKey").MustString("12344x")
	SecretKey = file.Section("ali").Key("SecretKey").MustString("xxxxxxxxxx")
}

func LoadPrometheus(file *ini.File) {
	PrometheusConfDir = file.Section("prometheus").Key("prometheus_config_dir").MustString("12344x")
}
func LoadAws(file *ini.File) {

	AccessKey = file.Section("aws").Key("AccessKey").MustString("xxxxxxxxxx")
	SecretKey = file.Section("aws").Key("SecretKey").MustString("xxxxxxxxxx")

}
func Loadk8s(file *ini.File) {
	KubeFile = file.Section("k8s").Key("KubeFile").MustString("xxxxxxxxxx")
}

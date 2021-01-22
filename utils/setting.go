package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	//"github.com/sirupsen/logrus"
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

	LogFile  string
	LogLevel string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
		panic(1)
	}
	LoadServer(file)
	LoadData(file)
	LoadAli(file)
	LoadAws(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
	LogFile = file.Section("server").Key("logfile").MustString("xxxx/cmdb")
}
func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("debug")
	DbHost = file.Section("database").Key("DbHost").MustString("47.104.197.46")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("caichangen")
	DbName = file.Section("database").Key("DbName").MustString("cmdb")
}

func LoadAli(file *ini.File) {
	AccessKey = file.Section("ali").Key("AccessKey").MustString("xxxxxxxxxx")
	SecretKey = file.Section("ali").Key("SecretKey").MustString("xxxxxxxxxx")
}

func LoadAws(file *ini.File) {
	AccessKey = file.Section("aws").Key("AccessKey").MustString("xxxxxxxxxx")
	SecretKey = file.Section("aws").Key("SecretKey").MustString("xxxxxxxxxx")
}

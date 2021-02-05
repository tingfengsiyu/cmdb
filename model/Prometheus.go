package model

import (
	"cmdb/middleware"
	"encoding/json"
	"os"

	"gorm.io/gorm"
)

//prometheus json配置文件格式
type PrometheusTarget struct {
	Targets []string `json:"targets"`
	Labels  Labels   `json:"labels"`
}

type Labels struct {
	Idc string `json:"idc"`
}

type MonitorPrometheus struct {
	gorm.Model
	ServerID int `gorm:"type:int;not null" json:"server_id"`
	//PrivateIpAddress    string `gorm:"type:varchar(30);not null" json:"private_ip_address"`
	NodeExportPort      int `gorm:"type:int;not null" json:"NodeExportPort"`
	ProcessExportPort   int `gorm:"type:int;not null" json:"ProcessExportPort"`
	ScriptExportPort    int `gorm:"type:int;not null" json:"ScriptExportPort"`
	NodeExportStatus    int `gorm:"type:int;not null" json:"NodeExportStatus"`
	ProcessExportStatus int `gorm:"type:int;not null" json:"ProcessExportStatus"`
	ScriptExportStatus  int `gorm:"type:int;not null" json:"ScriptExportStatus"`
}

func ReadPrometheus() {

}

func WritePrometheus() {

}

func DeletePrometheus() {

}
func InstallAgent() {

}

func ReadJsonfile(file string, data []PrometheusTarget) ([]PrometheusTarget, error) {
	filePtr, err := os.Open(file)
	var read = make([]PrometheusTarget, 0)
	if err != nil {
		middleware.SugarLogger.Errorf("Open file failed [Err:%s]", err.Error())
	}
	defer filePtr.Close()
	decode := json.NewDecoder(filePtr)
	err = decode.Decode(&read)
	if err != nil {
		middleware.SugarLogger.Errorf("Decoder failed", err.Error())
	}
	return read, err
}

func WriteJsonfile(file string, data []PrometheusTarget) error {
	// var Targets = make([]string, 0)
	// var b = make([]PrometheusTarget, 0)
	// var label = Labels{}
	// Targets = append(Targets, "172.16.0.1:9100")
	// label = Labels{
	// 	Idc: "This is  a  tests",
	// }
	// b = append(b, PrometheusTarget{Targets: Targets,
	// 	Labels: label,
	// })
	filePtr, err := os.Create(file)
	if err != nil {
		middleware.SugarLogger.Errorf("Create file failed", err.Error())
		return err
	}
	defer filePtr.Close()
	datas, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		middleware.SugarLogger.Errorf("Encoder failed", err.Error())
	}
	filePtr.Write(datas)
	return err
}

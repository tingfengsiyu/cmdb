package model

import (
	"cmdb/middleware"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"sync"
	"time"
)

//prometheus json配置文件格式
type PrometheusTarget struct {
	Targets []string `json:"targets"`
	Labels  Labels   `json:"labels"`
}
type Labels struct {
	Idc string `json:"idc"`
}

type ScanMonitorPrometheus struct {
	ServerID            int    `json:"server_id"`
	PrivateIpAddress    string `json:"private_ip_address"`
	NodeExportPort      string `json:"node_export_port"`
	ProcessExportPort   string `json:"process_export_port"`
	ScriptExportPort    string `json:"script_export_port"`
	NodeExportStatus    int    `json:"node_export_status"`
	ProcessExportStatus int    `json:"process_export_status"`
	ScriptExportStatus  int    `json:"script_export_status"`
	Label               string `json:"label"`
}

var wg sync.WaitGroup

func Croninit() {
	go func() {
		crontab := cron.New()
		crontab.AddFunc("0 */3 * * * *", CheckAgentStatus) // 每隔3分钟 定时执行 CheckAgentStatus 函数
		crontab.Start()
	}()
}

//并发检测监控agent运行状态
func CheckAgentStatus() {
	monitorPrometheus := Prometheus_server()
	var msgChan = make(chan []ScanMonitorPrometheus, len(monitorPrometheus))
	msgChan <- monitorPrometheus
	client := &http.Client{Timeout: 200 * time.Millisecond}
	wg.Add(1)
	go HHH(client, msgChan)
	wg.Wait()
}

func HHH(client *http.Client, msgChan chan []ScanMonitorPrometheus) {
	defer wg.Done()
	var monitorPrometheus, _ = <-msgChan
	for _, v := range monitorPrometheus {
		nodeExporterUrl := fmt.Sprintf("http://%s:%s", "172.22.0.20", v.NodeExportPort)
		processExporterUrl := fmt.Sprintf("http://%s:%s", "172.22.0.20", v.ProcessExportPort)
		scriptExporterUrl := fmt.Sprintf("http://%s:%s", "172.22.0.20", v.ScriptExportPort)
		nodeStatusCode := HttpCheckExporter(client, nodeExporterUrl)
		processStatusCode := HttpCheckExporter(client, processExporterUrl)
		scriptStatusCode := HttpCheckExporter(client, scriptExporterUrl)
		var monitor = make(map[string]interface{})
		monitor["node_export_status"] = nodeStatusCode
		monitor["process_export_status"] = processStatusCode
		monitor["script_export_status"] = scriptStatusCode
		db.Model(&MonitorPrometheus{}).Where("server_id =?", v.ServerID).Updates(monitor)
	}
}
func HttpCheckExporter(client *http.Client, url string) int {
	resp, err := client.Get(url)
	if err != nil {
		middleware.SugarLogger.Errorf("url %s", err)
		return 0
	}
	if resp.StatusCode == 200 {
		return 2
	}
	defer resp.Body.Close()
	return 0
}

func Prometheus_server() []ScanMonitorPrometheus {
	var svc []ScanMonitorPrometheus
	errs := db.Debug().Unscoped().Model(&MonitorPrometheus{}).Select("monitor_prometheus.server_id,private_ip_address,node_export_port,process_export_port,script_export_port,node_export_status,process_export_status,script_export_status,label  ").Joins("left join server on server.server_id=monitor_prometheus.server_id").Scan(&svc)
	if errs != nil {
		middleware.SugarLogger.Errorf("", errs)
		return svc
	}
	return svc
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

package model

import (
	"cmdb/middleware"
	"cmdb/utils"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//prometheus json配置文件格式
type PrometheusTarget struct {
	Targets []string `json:"targets"`
	Labels  Labels   `json:"labels"`
}
type Labels struct {
	//Idc     string `json:"idc"`
	Cluster string `json:"cluster"`
}

type ScanMonitorPrometheus struct {
	ServerID             int    `json:"server_id"`
	PrivateIpAddress     string `json:"private_ip_address"`
	NodeExportPort       string `json:"node_export_port"`
	ProcessExportPort    string `json:"process_export_port"`
	ScriptExportPort     string `json:"script_export_port"`
	NodeExportStatus     int    `json:"node_export_status"`
	ProcessExportStatus  int    `json:"process_export_status"`
	ScriptExportStatus   int    `json:"script_export_status"`
	Label                string `json:"label"`
	Cluster              string `json:"cluster"`
	DisableNodeExport    int    `json:"disable_node_export"`
	DisableProcessExport int    `json:"disable_process_export"`
	DisableScriptExport  int    `json:"disable_script_export"`
}

//并发检测监控agent运行状态
func CheckAgentStatus() {

	monitorPrometheus := PrometheusServer()
	client := &http.Client{Timeout: 200 * time.Millisecond}
	if err != nil {
		middleware.SugarLogger.Errorf("Open file failed [Err:%s]", err.Error())
	}
	for _, v := range monitorPrometheus {
		go func(client *http.Client, tmp ScanMonitorPrometheus) {
			ipaddress := strings.Trim(tmp.PrivateIpAddress, " ")
			nodeExporterUrl := "http://" + ipaddress + ":" + tmp.NodeExportPort
			processExporterUrl := "http://" + ipaddress + ":" + tmp.ProcessExportPort
			scriptExporterUrl := "http://" + ipaddress + ":" + tmp.ScriptExportPort
			nodeStatusCode := HttpCheckExporter(client, nodeExporterUrl)
			processStatusCode := HttpCheckExporter(client, processExporterUrl)
			scriptStatusCode := HttpCheckExporter(client, scriptExporterUrl)
			var monitor = make(map[string]interface{})
			monitor["node_export_status"] = nodeStatusCode
			monitor["process_export_status"] = processStatusCode
			monitor["script_export_status"] = scriptStatusCode
			db.Model(&MonitorPrometheus{}).Where("server_id =?", tmp.ServerID).Updates(monitor)
		}(client, v)

	}
}

func HttpCheckExporter(client *http.Client, url string) int {
	resp, err := client.Get(url)
	if err != nil {
		//middleware.SugarLogger.Errorf("%s", err)
		return 0
	}
	if resp.StatusCode == 200 {
		return 2
	}
	defer resp.Body.Close()
	return 0
}

func PrometheusServer() []ScanMonitorPrometheus {
	var svc []ScanMonitorPrometheus
	err := db.Model(&Server{}).Select("monitor_prometheus.server_id,private_ip_address,node_export_port,process_export_port,script_export_port,node_export_status," +
		"process_export_status,script_export_status,label,cluster,disable_node_export,disable_process_export,disable_script_export").Joins("left join" +
		" monitor_prometheus on server.server_id=monitor_prometheus.server_id").Order("cluster").Scan(&svc).Error
	if err != nil {
		middleware.SugarLogger.Errorf("", err)
		return svc
	}
	return svc
}

func WritePrometheus() {
	monitorPrometheus := PrometheusServer()
	var node = make([]string, 0)
	var script = make([]string, 0)
	maps := make(map[string][]string, 0)
	//sort.Slice(monitorPrometheus,func(i, j int) bool { return monitorPrometheus[i].Cluster < monitorPrometheus[j].Cluster })
	for _, v := range monitorPrometheus {
		ipaddress := strings.TrimSpace(v.PrivateIpAddress)
		if maps[v.Cluster+".json"] == nil {
			node = []string{}
		}
		if maps[v.Cluster+".yaml"] == nil {
			script = []string{}
		}
		if v.NodeExportStatus == 2 && v.DisableNodeExport == 1 {
			node = append(node, ipaddress+":"+v.NodeExportPort)
			maps[v.Cluster+".json"] = append(node)
		}
		if v.ProcessExportStatus == 2 && v.DisableProcessExport == 1 {
			node = append(node, ipaddress+":"+v.ProcessExportPort)
			maps[v.Cluster+".json"] = append(node)
		}
		if v.ScriptExportStatus == 2 && v.DisableScriptExport == 1 {
			script = append(script, ipaddress+":"+v.ScriptExportPort)
			maps[v.Cluster+".yaml"] = append(script)
		}
	}
	for k, v := range maps {
		if maps[k] != nil {
			//WriteJsonfile(strings.Split(k, ".")[0], utils.PrometheusConfDir+"/"+k, v)
			func(cluster, file string, targets []string) {
				var tmp = make([]PrometheusTarget, 0)
				tmp = append(tmp, PrometheusTarget{Targets: targets,
					Labels: Labels{
						Cluster: cluster,
					},
				})
				filePtr, err := os.Create(file)
				if err != nil {
					middleware.SugarLogger.Errorf("Create file failed", err.Error())
				}

				datas, err := json.MarshalIndent(tmp, "", "  ")
				if err != nil {
					middleware.SugarLogger.Errorf("Encoder failed", err.Error())
				}
				filePtr.Write(datas)
				defer filePtr.Close()
			}(strings.Split(k, ".")[0], utils.PrometheusConfDir+"/"+k, v)
		}
	}
}

func InsertPrometheusID(server_id int) {
	//更新数据到监控表,  监控状态 0 未安装 ；1 已安装 ；2 已运行；3 已被监控
	var prometheustarget = MonitorPrometheus{}
	var prometheus = make(map[string]interface{})
	prometheus["server_id"] = server_id
	db.Model(&prometheustarget).Create(prometheus)
}

func ReadJsonfile(filePtr *os.File, targets string) (int, error) {
	var read = make([]PrometheusTarget, 0)
	decode := json.NewDecoder(filePtr)
	err = decode.Decode(&read)
	if err != nil {
		middleware.SugarLogger.Errorf("Decoder failed", err.Error())
	}
	code := 0
	for _, v := range read {
		for _, target := range v.Targets {
			if target == targets {
				code = 1
				return code, nil
			}
		}
	}
	return code, err
}

type alert struct {
	Status string `json:"status"`
	Data   Data
}
type Data struct {
	Data   string `json:"data"`
	Alerts []Alerts
}
type Alerts struct {
	Value  string `json:"value"`
	Labels Label
}
type Label struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Cluster   string `json:"cluster"`
	Device    string `json:"device"`
}
type AlertInfo struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Cluster   string `json:"cluster"`
	Value     string `json:"value"`
}

func PrometheusAlerts() []AlertInfo {
	client := &http.Client{Timeout: 200 * time.Millisecond}
	resp, _ := client.Get("http://" + utils.PrometheusAddr + "/api/v1/alerts")
	body, _ := ioutil.ReadAll(resp.Body)
	var a alert
	var alertinfo = []AlertInfo{}
	_ = make([]string, 0)
	json.Unmarshal([]byte(body), &a)
	for _, v := range a.Data.Alerts {
		value, _ := strconv.ParseFloat(v.Value, 64)
		number := big.NewRat(1, 1)
		number.SetFloat64(value)
		alertinfo = append(alertinfo, AlertInfo{
			v.Labels.Alertname, v.Labels.Instance, v.Labels.Cluster, ""})
	}
	return alertinfo
}

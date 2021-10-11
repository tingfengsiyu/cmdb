package script

import (
	logger "cmdb/middleware"
	"cmdb/model"
	"cmdb/utils"
	"cmdb/utils/errmsg"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func UpdateHostName(c *gin.Context) {
	c.Copy()
	go model.UpdateHostName()
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "执行shell中",
		},
	)
}

func ShellInit(c *gin.Context) {
	var shellinit = model.OsInitStruct{}
	if err := c.ShouldBindJSON(&shellinit); err != nil {
		//c.String(400, "格式不符合要求", err.Error())
		errors(c, "格式不符合要求"+err.Error())
		return
	}
	var cmd string
	if shellinit.Role == "lotus-worker" {
		cmd = "/data/ops/script/osinit.sh " + shellinit.StorageMount.InitStartIP + " " + shellinit.StorageMount.InitEndNumber + " " + shellinit.InitUser + " " +
			shellinit.InitPass + " " + shellinit.Role + " " + shellinit.StorageMount.StorageStartIP + " " + shellinit.StorageMount.StorageStopnumber
	} else if shellinit.Role == "lotus-storage" {
		cmd = "/data/ops/script/osinit.sh " + shellinit.StorageMount.InitStartIP + " " + shellinit.StorageMount.InitEndNumber + " " + shellinit.InitUser + " " +
			shellinit.InitPass + " " + shellinit.Role
	} else {
		errors(c, "role错误 lotus-worker|lotus-storage")
		return
	}
	flag := ScriptIpVerfiy(shellinit.StorageMount.InitEndNumber, shellinit.StorageMount.InitStartIP)
	if !flag {
		errors(c, "集群ip不存在数据库，请确认ip")
		return
	}
	c.Copy()
	user := user(c)
	//osinit.sh  initStartIP initStopNumber  initUser initPass  Role  storageStartIP  storageStopnumber
	tmp := model.OpsRecords{
		User: user,
		Object: shellinit.StorageMount.InitStartIP + "-" + shellinit.StorageMount.InitEndNumber + "初始化用户名：" + shellinit.InitUser + "初始化密码" + shellinit.InitPass +
			"角色：" + shellinit.Role + "存储ip" + shellinit.StorageMount.StorageStartIP + "-" + shellinit.StorageMount.StorageStopnumber,
		Action: "系统初始化",
	}
	id := model.InsertRecords(tmp)
	go model.ExecLocalShell(id, cmd)
	//model.GenerateAnsibleHosts()
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "初始化系统中",
		},
	)
}

func BatchIp(c *gin.Context) {
	var batchip = model.BatchIpStruct{}
	if err := c.ShouldBindJSON(&batchip); err != nil {
		errors(c, "格式不符合要求"+err.Error())
		return
	}
	var ids = make([]int, 0)
	var ips = make([]string, 0)
	tmpEndNumber, _ := strconv.Atoi(batchip.SourceEndNumber)
	if len(strings.Split(batchip.SourceStartIp, ".")) != 4 || len(strings.Split(batchip.TargetStartIP, ".")) != 4 {
		errors(c, "源ip或目标ip输入错误")
	}
	sourceStartIpNumber, _ := strconv.Atoi(strings.Split(batchip.SourceStartIp, ".")[3])
	t := strings.Split(batchip.SourceStartIp, ".")
	startPrefix := t[0] + "." + t[1] + "." + t[2] + "."
	targetStartNumber, _ := strconv.Atoi(strings.Split(batchip.TargetStartIP, ".")[3])
	t = strings.Split(batchip.TargetStartIP, ".")
	targetPrefix := t[0] + "." + t[1] + "." + t[2] + "."

	if id := model.CheckClusterName(batchip.TargetClusterName); id <= 0 {
		errors(c, "集群名不存在数据库，请确认集群名")
		return
	}
	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		tmp := strconv.Itoa(i)
		id := model.CheckClusterIp(startPrefix + tmp)
		if id <= 0 {
			errors(c, "集群ip不存在数据库，请确认ip")
			return
		}
		ids = append(ids, id)
		ips = append(ips, targetPrefix+strconv.Itoa(targetStartNumber))
		id = model.CheckClusterIp(targetPrefix + strconv.Itoa(targetStartNumber))
		if id > 0 {
			errors(c, "目标集群ip已存在数据库，冲突，请确认ip")
			return
		}
		targetStartNumber += 1
	}
	user := user(c)
	tmp := model.OpsRecords{
		User:   user,
		Object: "源操作：" + batchip.SourceStartIp + "-" + batchip.SourceEndNumber + "源网关" + batchip.SourceGateway + "\n目标操作：" + batchip.TargetStartIP + "目标网关" + batchip.TargetGateway + "目标集群" + batchip.TargetClusterName,
		Action: "修改ip",
	}
	id := model.InsertRecords(tmp)
	for k, v := range ids {
		model.UpdateClusterName(v, ips[k], batchip.TargetClusterName)
	}
	model.GenerateAnsibleHosts()
	model.GenerateClustersHosts()
	//sh batchip.sh sourceIP sourceGateway sourceEndNumber targetStartIP targetGateway
	cmd := "/data/ops/script/batchip.sh " + batchip.SourceStartIp + " " + batchip.SourceGateway + " " +
		batchip.SourceEndNumber + " " + batchip.TargetStartIP + " " + batchip.TargetGateway
	go model.ExecLocalShell(id, cmd)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "批量修改IP中",
		},
	)
}

func StorageMount(c *gin.Context) {
	var storagemount = model.StorageMountStruct{}
	if err := c.ShouldBindJSON(&storagemount); err != nil {
		errors(c, "格式不符合要求"+err.Error())
		return
	}
	flag := ScriptIpVerfiy(storagemount.InitEndNumber, storagemount.InitStartIP)
	if !flag {
		errors(c, "集群ip不存在数据库，请确认ip")
		return
	}
	user := user(c)
	tmp := model.OpsRecords{
		User:   user,
		Object: "worker: " + storagemount.InitStartIP + "-" + storagemount.InitEndNumber + " 存储: " + storagemount.StorageStartIP + "-" + storagemount.StorageStopnumber + " worker操作: " + storagemount.Operating,
		Action: "挂载存储",
	}
	id := model.InsertRecords(tmp)
	c.Copy()
	// batchStorage-mount.sh  sourceIP  sourceEndNumber storageStartIP storageEndNumber operating
	cmd := "/data/ops/script/batchStorage-mount.sh  " + storagemount.InitStartIP + " " + storagemount.InitEndNumber + " " + storagemount.StorageStartIP + " " + storagemount.StorageStopnumber + " " + storagemount.Operating
	go model.ExecLocalShell(id, cmd)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "存储生成并挂载中",
		},
	)
}

func WritePrometheus(c *gin.Context) {
	model.WritePrometheus()
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "write ok!!!",
		},
	)
}

func InstallMointorAgent(c *gin.Context) {
	clustername := c.Query("clustername")
	if clustername == "" {
		c.String(400, "clustername不能为空")
		return
	}
	user := user(c)
	tmp := model.OpsRecords{
		User:   user,
		Object: "cluster: " + clustername,
		Action: "安装监控agent",
	}
	id := model.InsertRecords(tmp)
	c.Copy()
	model.GenerateAnsibleHosts()
	model.GenerateClustersHosts()
	model.ScanHardWareInfo()
	go model.ExecLocalShell(id, "/data/ops/script/monitoragent.sh "+clustername)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "agent安装中",
		},
	)
}
func AnsiblePlaybook(c *gin.Context) {
	type Playbook struct {
		ContinuousIp    string   `json:"continuous_ip"`    //连续ip
		DiscontinuousIp string   `json:"discontinuous_ip"` //不连续ip
		Group           []string `json:"group"`
		Variable        string   `json:"variable"`
		Tag             string   `json:"tag"  binding:"required"`
		FileName        string   `json:"filename" binding:"required"`
	}
	var ansiblePlaybook Playbook
	err := c.ShouldBindJSON(&ansiblePlaybook)
	if err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  404,
				"message": "必选字段不存在",
			},
		)
		return
	}
	//Generate ip list

	user := user(c)
	var group, cmd string
	ips := make([]string, 0)
	for _, v := range ansiblePlaybook.Group {
		group = group + " " + v
	}
	//连续ip
	tmpips := strings.Split(ansiblePlaybook.ContinuousIp, ",")
	if len(tmpips) >= 1 {
		for _, s := range tmpips {
			tmp := strings.Split(s, "-")
			if len(tmp) >= 2 {
				t := strings.Split(tmp[0], ".")
				prefix := t[0] + "." + t[1] + "." + t[2] + "."
				tmpEnd, _ := strconv.Atoi((tmp[1]))
				tmpStart, _ := strconv.Atoi(t[3])
				for i := tmpStart; i <= tmpEnd; i++ {
					ip := strings.TrimSpace(prefix + strconv.Itoa(i))
					if len(ip) > 2 {
						server, _ := model.NetworkTopology(0, "", "", "", "", ip)
						for _, v := range server {
							str := v.PrivateIpAddress + " roles=" + v.Label + " hostname=" + v.Name
							ips = append(ips, str)
						}
					}
				}
			}
		}
	}
	//不连续ip
	tmpip := strings.Split(ansiblePlaybook.DiscontinuousIp, ",")
	if len(tmpip) >= 1 {
		for _, s := range tmpip {
			if len(strings.TrimSpace(s)) > 2 {
				server, _ := model.NetworkTopology(0, "", "", "", "", strings.TrimSpace(s))
				for _, v := range server {
					str := v.PrivateIpAddress + " roles=" + v.Label + " hostname=" + v.Name
					ips = append(ips, str)
				}
			}
		}
	}
	if len(ips)+len(group) == 0 {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  404,
				"message": "操作ip不存在或组未选中",
			},
		)
		return
	}
	id := model.InsertRecords(model.OpsRecords{
		User:   user,
		Object: ansiblePlaybook.ContinuousIp + "," + ansiblePlaybook.DiscontinuousIp + "  " + ansiblePlaybook.Variable,
		Action: "执行ansible" + ansiblePlaybook.FileName + " " + ansiblePlaybook.Tag,
	})

	model.GenerateAnsibleHosts()
	model.GenerateClustersHosts()
	model.AppendAnsibleHost(ips)
	if group == "" {
		cmd = "cd  " + utils.AnsiblePlaybookDir + "; ansible-playbook --limit   tmplotus  -i " + utils.AnsibleHosts + "-tmplotus  -t " + ansiblePlaybook.Tag + " " + ansiblePlaybook.FileName + " -e " + "'" + ansiblePlaybook.Variable + "'"
	} else {
		cmd = "cd  " + utils.AnsiblePlaybookDir + "; ansible-playbook --limit    " + "'" + group + "'" + " -t " + ansiblePlaybook.Tag + " " + ansiblePlaybook.FileName + " -e " + "'" + ansiblePlaybook.Variable + "'"
	}
	c.Copy()
	go model.ExecLocalShell(id, cmd)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "任务执行中",
		},
	)
}

func GenerateAnsibleHosts(c *gin.Context) {
	var code int
	if err := model.GenerateAnsibleHosts(); err != nil {
		code = 4003
	} else {
		code = 200
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
func AnsibleHosts(c *gin.Context) {
	maps := model.AllHosts()
	code := 200
	type cluster struct {
		Id    int    `json:"id"`
		Group string `json:"cluster"`
	}
	var host = []cluster{}
	for k, _ := range maps {
		host = append(host, cluster{Id: 0, Group: k})
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    host,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
func AnsiblePlaybookFiles(c *gin.Context) {
	code := 200
	type cluster struct {
		Id       int    `json:"id"`
		FileName string `json:"filename"`
	}
	var yamls = []cluster{}
	files, err := ioutil.ReadDir(utils.AnsiblePlaybookDir)
	if err != nil {
		logger.SugarLogger.Error(err)
	}

	for _, file := range files {
		if ok, _ := regexp.MatchString(".yml|.yaml", file.Name()); ok {
			yamls = append(yamls, cluster{0,
				file.Name()})
		}
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    yamls,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func GenerateClustersHosts(c *gin.Context) {
	code := 200
	model.GenerateClustersHosts()
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func ScriptIpVerfiy(initNumber, initStartIp string) bool {
	tmpEndNumber, _ := strconv.Atoi(initNumber)
	sourceStartIpNumber, _ := strconv.Atoi(strings.Split(initStartIp, ".")[3])
	t := strings.Split(initStartIp, ".")
	var ipLists = make([]string, 0)
	startPrefix := t[0] + "." + t[1] + "." + t[2] + "."
	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		tmp := strconv.Itoa(i)
		ipLists = append(ipLists, startPrefix+tmp)
	}
	flag := model.BatchCheckClusterIps(ipLists)
	return flag
}

func UpdateCluster(c *gin.Context) {
	var cluster = model.UpdateClusterStruct{}
	if err := c.ShouldBindJSON(&cluster); err != nil {
		errors(c, "格式不符合要求"+err.Error())
		return
	}
	var ids = make([]int, 0)
	var ips = make([]string, 0)
	if len(strings.Split(cluster.SourceStartIp, ".")) != 4 {
		errors(c, "源ip输入错误")
	}
	tmpEndNumber, _ := strconv.Atoi(cluster.SourceEndNumber)
	sourceStartIpNumber, _ := strconv.Atoi(strings.Split(cluster.SourceStartIp, ".")[3])
	t := strings.Split(cluster.SourceStartIp, ".")
	startPrefix := t[0] + "." + t[1] + "." + t[2] + "."
	if id := model.CheckClusterName(cluster.TargetClusterName); id <= 0 {
		errors(c, "集群名不存在数据库，请确认集群名")
		return
	}
	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		tmp := strconv.Itoa(i)
		id := model.CheckClusterIp(startPrefix + tmp)
		if id <= 0 {
			errors(c, "集群ip不存在数据库，请确认ip")
			return
		}
		ids = append(ids, id)
		ips = append(ips, startPrefix+strconv.Itoa(sourceStartIpNumber))
		sourceStartIpNumber += 1
	}

	for k, v := range ids {
		model.UpdateClusterName(v, ips[k], cluster.TargetClusterName)
	}
	user := user(c)
	tmp := model.OpsRecords{
		User:   user,
		Object: "源操作：" + cluster.SourceStartIp + "-" + cluster.SourceEndNumber + "\n目标集群: " + cluster.TargetClusterName,
		Action: "修改机器所属集群",
		State:  1,
	}
	model.InsertRecords(tmp)

	model.GenerateAnsibleHosts()
	//追加生成 ansible hosts  worker
	model.GenerateClustersHosts()
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "ok",
		},
	)

}
func ExecWebShell(c *gin.Context) {
	tmpfile := utils.LogFile + "cmd"
	var cluster = model.UpdateClusterStruct{}
	if err := c.ShouldBindJSON(&cluster); err != nil {
		errors(c, "格式不符合要求"+err.Error())
	}
	var ips = make([]string, 0)
	if len(strings.Split(cluster.SourceStartIp, ".")) != 4 {
		errors(c, "源ip输入错误")
	}
	tmpEndNumber, _ := strconv.Atoi(cluster.SourceEndNumber)
	sourceStartIpNumber, _ := strconv.Atoi(strings.Split(cluster.SourceStartIp, ".")[3])
	t := strings.Split(cluster.SourceStartIp, ".")
	startPrefix := t[0] + "." + t[1] + "." + t[2] + "."
	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		tmp := strconv.Itoa(i)
		id := model.CheckClusterIp(startPrefix + tmp)
		if id <= 0 {
			errors(c, "集群ip不存在数据库，请确认ip")
			return
		}
		ips = append(ips, startPrefix+strconv.Itoa(sourceStartIpNumber))
		sourceStartIpNumber += 1
	}
	user := user(c)
	var tmp = model.OpsRecords{
		User:   user,
		Object: "源操作：" + cluster.SourceStartIp + "-" + cluster.SourceEndNumber + "\n执行命令: " + cluster.Cmd[0:10],
		Action: "执行命令",
		State:  1,
	}
	id := model.InsertRecords(tmp)
	model.GenerateAnsibleHosts()
	file, err := os.OpenFile(tmpfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.SugarLogger.Errorf("写入文件错误!!!%s", err)
	}
	defer file.Close()
	file.WriteString("#!/bin/bash\n")
	file.WriteString(cluster.Cmd)
	go func() {
		cmd := "scp.sh " + " host-tmpworker " + tmpfile + " /tmp/cmd"
		model.ExecLocalShell(0, cmd)
		cmd = "execshell.sh " + " host-tmpworker " + " /tmp/cmd"
		model.ExecLocalShell(id, cmd)
	}()

	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "shell执行中",
		},
	)

}

func errors(c *gin.Context, str string) {
	c.JSON(
		http.StatusOK, gin.H{
			"status":  400,
			"message": str,
		},
	)
}

func user(c *gin.Context) string {
	user, _ := c.Get("username")
	return user.(string)
}

func PrometheusAlerts(c *gin.Context) {
	alertinfo := model.PrometheusAlerts()
	sort.Slice(alertinfo, func(i, j int) bool { return alertinfo[i].Cluster < alertinfo[j].Cluster })
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   alertinfo,
			"total":  len(alertinfo),
		},
	)
}

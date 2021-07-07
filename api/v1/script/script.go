package script

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func UpdateHostName(c *gin.Context) {
	c.Copy()
	go model.UpdateHostName()
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   "执行shell中",
		},
	)
}

func ShellInit(c *gin.Context) {
	var shellinit = model.OsInitStruct{}
	if err := c.ShouldBindJSON(&shellinit); err != nil {
		c.String(400, "格式不符合要求", err.Error())
		return
	}
	var cmd string
	if shellinit.Role == "lotus-worker" {
		cmd = "osinit.sh " + shellinit.StorageMount.InitStartIP + " " + shellinit.StorageMount.InitEndNumber + " " + shellinit.InitUser + " " +
			shellinit.InitPass + " " + shellinit.Role + " " + shellinit.StorageMount.StorageStartIP + " " + shellinit.StorageMount.StorageStopnumber
	} else if shellinit.Role == "lotus-storage" {
		cmd = "osinit.sh " + shellinit.StorageMount.InitStartIP + " " + shellinit.StorageMount.InitEndNumber + " " + shellinit.InitUser + " " +
			shellinit.InitPass + " " + shellinit.Role
	} else {
		c.String(400, "role错误 lotus-worker|lotus-storage")
		return
	}
	flag := ScriptIpVerfiy(shellinit.StorageMount.InitEndNumber, shellinit.StorageMount.InitStartIP)
	if !flag {
		c.String(400, "集群ip不存在数据库,请确认ip已录入")
		return
	}
	c.Copy()
	//osinit.sh  initStartIP initStopNumber  initUser initPass  Role  storageStartIP  storageStopnumber
	tmp := model.OpsRecords{
		Object: shellinit.StorageMount.InitStartIP + "-" + shellinit.StorageMount.InitEndNumber + "初始化用户名：" + shellinit.InitUser + "初始化密码" + shellinit.InitPass +
			"角色：" + shellinit.Role + "存储ip" + shellinit.StorageMount.StorageStartIP + "-" + shellinit.StorageMount.StorageStopnumber,
		Action: "系统初始化",
	}
	id := model.InsertRecords(tmp)
	go model.ExecLocalShell(id, cmd)
	//model.GenerateAnsibleHosts()
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   "初始化系统中",
		},
	)
}

func BatchIp(c *gin.Context) {
	var batchip = model.BatchIpStruct{}
	if err := c.ShouldBindJSON(&batchip); err != nil {
		c.String(400, "格式不符合要求", err.Error())
		return
	}
	var ids = make([]int, 0)
	var ips = make([]string, 0)
	tmpEndNumber, _ := strconv.Atoi(batchip.SourceEndNumber)
	sourceStartIpNumber, _ := strconv.Atoi(strings.Split(batchip.SourceStartIp, ".")[3])
	t := strings.Split(batchip.SourceStartIp, ".")
	startPrefix := t[0] + "." + t[1] + "." + t[2] + "."
	targetStartNumber, _ := strconv.Atoi(strings.Split(batchip.TargetStartIP, ".")[3])
	t = strings.Split(batchip.TargetStartIP, ".")
	targetPrefix := t[0] + "." + t[1] + "." + t[2] + "."

	if id := model.CheckClusterName(batchip.TargetClusterName); id <= 0 {
		c.String(400, "集群名不存在数据库，请确认后修改ip")
		return
	}
	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		tmp := strconv.Itoa(i)
		id := model.CheckClusterIp(startPrefix + tmp)
		if id <= 0 {
			c.String(400, "集群ip不存在数据库，请确认后修改ip")
			return
		}
		ids = append(ids, id)
		ips = append(ips, targetPrefix+strconv.Itoa(targetStartNumber))
		targetStartNumber += 1
	}
	tmp := model.OpsRecords{
		Object: "源操作：" + batchip.SourceStartIp + "-" + batchip.SourceEndNumber + "源网关" + batchip.SourceGateway + "\n目标操作：" + batchip.TargetStartIP + "目标网关" + batchip.TargetGateway + "目标集群" + batchip.TargetClusterName,
		Action: "修改ip",
	}
	id := model.InsertRecords(tmp)
	//sh batchip.sh sourceIP sourceGateway sourceEndNumber targetStartIP targetGateway
	cmd := "batchip.sh " + batchip.SourceStartIp + " " + batchip.SourceGateway + " " +
		batchip.SourceEndNumber + " " + batchip.TargetStartIP + " " + batchip.TargetGateway
	go model.ExecLocalShell(id, cmd)
	for k, v := range ids {
		model.UpdateClusterName(v, ips[k], batchip.TargetClusterName)
	}
	model.GenerateAnsibleHosts()
	model.AppendAnsibleHosts(ips, batchip.TargetClusterName)
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   "批量修改IP中",
		},
	)
}

func StorageMount(c *gin.Context) {
	var storagemount = model.StorageMountStruct{}
	if err := c.ShouldBindJSON(&storagemount); err != nil {
		c.String(400, "格式不符合要求", err.Error())
		return
	}
	flag := ScriptIpVerfiy(storagemount.InitEndNumber, storagemount.InitStartIP)
	if !flag {
		c.String(400, "集群ip不存在数据库,请确认ip已录入")
		return
	}
	tmp := model.OpsRecords{
		Object: "worker: " + storagemount.InitStartIP + "-" + storagemount.InitEndNumber + " 存储: " + storagemount.StorageStartIP + "-" + storagemount.StorageStopnumber + " worker操作: " + storagemount.Operating,
		Action: "挂载存储",
	}
	id := model.InsertRecords(tmp)
	c.Copy()
	// batchStorage-mount.sh  sourceIP  sourceEndNumber storageStartIP storageEndNumber operating
	cmd := "batchStorage-mount.sh  " + storagemount.InitStartIP + " " + storagemount.InitEndNumber + " " + storagemount.StorageStartIP + " " + storagemount.StorageStopnumber + " " + storagemount.Operating
	go model.ExecLocalShell(id, cmd)
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   "存储生成并挂载中",
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
	tmp := model.OpsRecords{
		Object: "cluster: " + clustername,
		Action: "安装监控agent",
	}
	id := model.InsertRecords(tmp)
	c.Copy()
	model.GenerateAnsibleHosts()
	go model.ExecLocalShell(id, "monitoragent.sh "+clustername)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "agent安装中",
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

func GenerateClustersHosts(c *gin.Context) {
	var code int
	clusters, _ := model.GetClusters()
	for _, v := range clusters {
		var ips []string
		if err := model.SyncTargetHosts(ips, v.Cluster); err != nil {
			code = 4003
		} else {
			code = 200
		}
	}

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
		c.String(400, "格式不符合要求", err.Error())
		return
	}
	var ids = make([]int, 0)
	var ips = make([]string, 0)
	tmpEndNumber, _ := strconv.Atoi(cluster.SourceEndNumber)
	sourceStartIpNumber, _ := strconv.Atoi(strings.Split(cluster.SourceStartIp, ".")[3])
	t := strings.Split(cluster.SourceStartIp, ".")
	startPrefix := t[0] + "." + t[1] + "." + t[2] + "."
	if id := model.CheckClusterName(cluster.TargetClusterName); id <= 0 {
		c.String(400, "集群名不存在数据库，请确认后修改ip")
		return
	}
	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		tmp := strconv.Itoa(i)
		id := model.CheckClusterIp(startPrefix + tmp)
		if id <= 0 {
			c.String(400, "集群ip不存在数据库，请确认后修改ip")
			return
		}
		ids = append(ids, id)
		ips = append(ips, startPrefix+strconv.Itoa(sourceStartIpNumber))
		sourceStartIpNumber += 1
	}

	for k, v := range ids {
		model.UpdateClusterName(v, ips[k], cluster.TargetClusterName)
	}
	tmp := model.OpsRecords{
		Object: "源操作：" + cluster.SourceStartIp + "-" + cluster.SourceEndNumber + "\n目标集群: " + cluster.TargetClusterName,
		Action: "修改机器所属集群",
		State:  1,
	}
	model.InsertRecords(tmp)

	model.GenerateAnsibleHosts()
	model.AppendAnsibleHosts(ips, cluster.TargetClusterName)
	//追加生成 ansible hosts  worker

	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   "ok",
		},
	)

}

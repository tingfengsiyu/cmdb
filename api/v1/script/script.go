package script

import (
	"cmdb/middleware"
	"cmdb/model"
	"cmdb/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func OsInit(c *gin.Context) {
	c.Copy()
	monitorPrometheus := model.PrometheusServer()
	for _, v := range monitorPrometheus {
		go func(v model.ScanMonitorPrometheus) {
			var user, passwd, sudopasswd string
			if v.Label == "lotus-miner" || v.Label == "lotus-worker" {
				user = utils.WorkerUser
				passwd = utils.WorkerPass
				sudopasswd = utils.WorkerPass
			} else if v.Label == "lotus-storage" {
				user = utils.StorageUser
				passwd = utils.StoragePass
				sudopasswd = utils.StorageSudoPass
			} else {
				user = utils.WorkerUser
				passwd = utils.WorkerPass
				sudopasswd = utils.WorkerPass
			}
			sshdConfig := "sudo sed -i 's@PasswordAuthentication no@PasswordAuthentication yes@g' /etc/ssh/sshd_config"
			updatePass := fmt.Sprintf("sudo echo root:%s | chpasswd", utils.RootPass)
			updatePubKey := fmt.Sprintf("sudo grep ops /root/.ssh/authorized_keys || sudo sed -i '1i %s' /root/.ssh/authorized_keys ", utils.RootPub)
			outs, err := model.SshCommands(user, passwd, v.PrivateIpAddress+":"+"22", sudopasswd, sshdConfig, updatePass, updatePubKey)
			if err != nil {
				middleware.SugarLogger.Errorf("ssh commands  %s ", err)
			}
			middleware.SugarLogger.Infof("%s ", string(outs))
		}(v)

	}
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"data":   "系统初始化中",
		},
	)
}

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
	if shellinit.Role == "lotus-worker" || shellinit.Role == "lotus-storage" {
	} else {
		c.String(400, "role错误 lotus-worker|lotus-storage")
		return
	}
	c.Copy()
	//osinit.sh  initStartIP initStopNumber  initUser initPass  Role  storageStartIP  storageStopnumber
	cmd := "/root/ops/osinit.sh " + shellinit.StorageMount.InitStartIP + " " + shellinit.StorageMount.InitEndNumber + " " + shellinit.InitUser + " " +
		shellinit.InitPass + " " + shellinit.Role + " " + shellinit.StorageMount.StorageStartIP + " " + shellinit.StorageMount.StorageStopnumber
	go model.ExecLocalShell(cmd)
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
	targetStartNumber, _ := strconv.Atoi(strings.Split(batchip.TargetStartIP, ".")[3])
	targetPrefix := strings.Replace(batchip.TargetStartIP, strings.Split(batchip.TargetStartIP, ".")[3], "", -1)

	for i := sourceStartIpNumber; i <= tmpEndNumber; i++ {
		id := model.CheckClusterName(batchip.SourceStartIp, batchip.TargetClusterName)
		if id <= 0 {
			c.String(400, "集群名和ip不存在数据库，请确认后修改ip")
			return
		}
		ids = append(ids, id)
		ips = append(ips, targetPrefix+strconv.Itoa(targetStartNumber))
		targetStartNumber += 1
	}

	//sh /root/ops/batchip.sh sourceIP sourceGateway sourceEndNumber targetStartIP targetGateway
	cmd := "/root/ops/batchip.sh " + batchip.SourceStartIp + " " + batchip.SourceGateway + " " +
		batchip.SourceEndNumber + " " + batchip.TargetStartIP + " " + batchip.TargetGateway
	model.ExecLocalShell(cmd)
	for k, v := range ids {
		model.UpdateClusterName(v, ips[k], batchip.TargetClusterName)
	}
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
	c.Copy()
	// batchStorage-mount.sh  sourceIP  sourceEndNumber storageStartIP storageEndNumber
	cmd := "/root/ops/batchStorage-mount.sh  " + storagemount.InitStartIP + " " + storagemount.InitEndNumber + " " + storagemount.StorageStartIP + " " + storagemount.StorageStopnumber
	go model.ExecLocalShell(cmd)
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
	c.Copy()
	go model.InstallAgent(clustername)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "agent安装中",
		},
	)
}

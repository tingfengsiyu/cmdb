package script

import (
	"cmdb/middleware"
	"cmdb/model"
	"cmdb/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OsInit(c *gin.Context) {
	c.Copy()
	monitorPrometheus := model.Prometheus_server()
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
			outs, err := model.SshCommands(user, passwd, v.PrivateIpAddress, sudopasswd, sshdConfig, updatePass, updatePubKey)
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

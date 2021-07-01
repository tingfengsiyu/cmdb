package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加服务器
func AddServer(c *gin.Context) {
	var data model.Server
	_ = c.ShouldBindJSON(&data)
	_, code := model.CheckServer(data.Name)
	if code == errmsg.SUCCSE {
		model.CreateServer(&data)
	}
	if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func BatchAddServers(c *gin.Context) {
	var code int
	var hostNames, idcNames, cabinetNumbers, citys []string
	assets := []model.ScanServers{}
	servers := []model.Server{}
	if err := c.ShouldBindJSON(&assets); err != nil {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			},
		)
	}
	for _, v := range assets {
		hostNames = append(hostNames, v.Name)
		idcNames = append(idcNames, v.IDC_Name)
		cabinetNumbers = append(cabinetNumbers, v.Cabinet_Number)
		citys = append(citys, v.City)
		servers = append(servers, model.Server{
			Name:             v.Name,
			Models:           v.Models,
			Location:         v.Location,
			PrivateIpAddress: v.PrivateIpAddress,
			PublicIpAddress:  v.PublicIpAddress,
			Label:            v.Label,
			Cluster:          v.Cluster,
			LabelIpAddress:   v.LabelIpAddress,
			Cpu:              v.Cpu,
			Memory:           v.Memory,
			Disk:             v.Disk,
			User:             v.User,
			State:            v.State,
		})
	}

	code = model.BatchCheckServer(hostNames)
	code = addServerVerify(code, servers, idcNames, cabinetNumbers, hostNames, citys)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func GetServer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	data, code := model.GetServerInfo(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   1,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
func GetCluster(c *gin.Context) {
	cluster := c.Query("cluster")
	data, code := model.GetCluster(cluster)
	c.JSON(
		http.StatusOK, gin.H{
			"status": code,
			"data":   data,
			"total":  len(data),
		},
	)
}

func Cron(c *gin.Context) {
	c.Copy()
	go model.CheckAgentStatus()
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"total":  1,
		},
	)
}

func GetServers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetServers(pageSize, pageNum)
	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func GetClusters(c *gin.Context) {

	data, total := model.GetClusters()
	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func UpdateServer(c *gin.Context) {
	var data model.Server
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code := model.EditServer(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func BatchUpdateServers(c *gin.Context) {
	var code int
	assets := []model.ScanServers{}
	var idcNames, citys, cabinetNumbers, hostNames []string
	var CabinetID []int
	_ = c.ShouldBindJSON(&assets)
	var IDS []int
	for _, v := range assets {
		IDS = append(IDS, v.ID)
	}
	code = model.BatchCheckServerID(IDS)
	if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	} else if code == errmsg.ERROR_ALL_DEVICE_EXIST {

		for _, v := range assets {
			idcNames = append(idcNames, v.IDC_Name)
			citys = append(citys, v.City)
			cabinetNumbers = append(cabinetNumbers, v.Cabinet_Number)
			hostNames = append(hostNames, v.Name)
			CabinetID = append(CabinetID, v.Cabinet_NumberID)
			var maps = make(map[string]interface{})
			maps["name"] = v.Name
			maps["models"] = v.Models
			maps["location"] = v.Location
			maps["private_ip_address"] = v.PrivateIpAddress
			maps["public_ip_address"] = v.PublicIpAddress
			maps["label_ip_address"] = v.LabelIpAddress
			maps["label"] = v.Label
			maps["cluster"] = v.Cluster
			maps["cpu"] = v.Cpu
			maps["memory"] = v.Memory
			maps["disk"] = v.Disk
			maps["user"] = v.User
			maps["state"] = v.State
			code = model.BatchUpdateServer(maps, v.ID)
		}

		idcIds := model.GenerateIDCID(idcNames)

		//生成cabinetNumberId
		cabinetNumberIds := model.GenerateCabinetID(cabinetNumbers, idcIds)

		//生成server_id
		serverIds := model.GenerateServerID(hostNames)
		//插入对应id
		for k, _ := range serverIds {
			if len(idcNames)-1 < k {
				idcNames = append(idcNames, idcNames[0])
			}
			if len(citys)-1 < k {
				citys = append(citys, citys[0])
			}
			if len(cabinetNumbers)-1 < k {
				cabinetNumbers = append(cabinetNumbers, cabinetNumbers[0])
			}
			if len(idcIds)-1 < k {
				idcIds = append(idcIds, idcIds[0])
			}
			cabinetNumberId, _ := model.CheckCabinetNumber(cabinetNumbers[k], idcIds[k])
			if cabinetNumberId == 0 {
				model.InsertCabinetID(cabinetNumbers[k], idcIds[k], cabinetNumberIds[k])
			}
			idc_id, _ := model.CheckIdcName(idcNames[k])
			if idc_id == 0 {
				model.InsertIdcID(idcNames[k], citys[k], idcIds[k], cabinetNumberIds[k])
			}
			model.InsertServerID(hostNames[k], idcIds[k], serverIds[k], cabinetNumberIds[k])

		}
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func DeleteServer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteServer(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func GetIdcServers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	idcname := c.Query("idcname")
	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetIdcServers(pageSize, pageNum, idcname)
	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func GetCabinetServers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	idcname := c.Query("idcname")
	cabinet_number := c.Query("cabinet_number")
	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetCabinetServers(pageSize, pageNum, idcname, cabinet_number)
	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func GetUser(c *gin.Context) {
	data, total := model.GetOwnedUser()
	var Users = make([]string, 0)
	for _, v := range data {
		Users = append(Users, v.User)
	}

	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    Users,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func addServerVerify(code int, servers []model.Server, idcNames, cabinetNumbers, hostNames, citys []string) int {
	if code == errmsg.SUCCSE {
		code = model.BatchCreateServer(&servers)
		//检查不存在后执行

		//生成idc_id
		idcIds := model.GenerateIDCID(idcNames)

		//生成cabinetNumberId
		cabinetNumberIds := model.GenerateCabinetID(cabinetNumbers, idcIds)

		//生成server_id
		serverIds := model.GenerateServerID(hostNames)

		//插入对应id
		for k, _ := range serverIds {
			if len(idcNames)-1 < k {
				idcNames = append(idcNames, idcNames[0])
			}
			if len(citys)-1 < k {
				citys = append(citys, citys[0])
			}
			if len(cabinetNumbers)-1 < k {
				cabinetNumbers = append(cabinetNumbers, cabinetNumbers[0])
			}
			if len(idcIds)-1 < k {
				idcIds = append(idcIds, idcIds[0])
			}
			model.InsertIdcID(idcNames[k], citys[k], idcIds[k], cabinetNumberIds[k])
			model.InsertServerID(hostNames[k], idcIds[k], serverIds[k], cabinetNumberIds[k])
			model.InsertCabinetID(cabinetNumbers[k], idcIds[k], cabinetNumberIds[k])
			model.InsertPrometheusID(serverIds[k])
		}
	} else if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	} else if code == errmsg.ERROR_ALL_DEVICE_EXIST {
		code = errmsg.ERROR_ALL_DEVICE_EXIST
	}
	return code
}

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
	var hostNames = make([]string, 0)
	assets := model.Assets{}
	_ = c.ShouldBindJSON(&assets)
	for _, v := range assets.Asset.Servers {
		hostNames = append(hostNames, v.Name)
	}
	idcNames := assets.Asset.Idcs.Idc_name
	cabinetNumbers := assets.Asset.Idcs.Cabinet_Number
	citys := assets.Asset.Idcs.City
	code = model.BatchCheckServer(hostNames)
	if code == errmsg.SUCCSE {
		code = model.BatchCreateServer(&assets.Asset.Servers)
		//检查不存在后执行
		//生成idc_id
		idc_ids := model.GenerateIDCID(idcNames)

		//生成cabinet_number_id
		cabinet_number_ids := model.GenerateCabinetID(cabinetNumbers, idc_ids)

		//生成server_id
		server_ids := model.GenerateServerID(hostNames)

		//插入对应id
		for k, _ := range server_ids {
			if len(idcNames)-1 < k {
				idcNames = append(idcNames, idcNames[0])
			}
			if len(citys)-1 < k {
				citys = append(citys, citys[0])
			}
			if len(cabinetNumbers)-1 < k {
				cabinetNumbers = append(cabinetNumbers, cabinetNumbers[0])
			}
			if len(idc_ids)-1 < k {
				idc_ids = append(idc_ids, idc_ids[0])
			}
			model.InsertIdcID(idcNames[k], citys[k], idc_ids[k], cabinet_number_ids[k])
			model.InsertServerID(hostNames[k], idc_ids[k], server_ids[k], cabinet_number_ids[k])
			model.InsertCabinetID(cabinetNumbers[k], idc_ids[k], cabinet_number_ids[k])
			model.InsertPrometheusID(server_ids[k])
		}
	} else if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	} else if code == errmsg.ERROR_ALL_DEVICE_EXIST {
		code = errmsg.ERROR_ALL_DEVICE_EXIST
	}
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
	var IDCID, CabinetID []int
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
			IDCID = append(CabinetID, v.IDC_ID)
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

		idc_ids := model.GenerateIDCID(idcNames)

		//生成cabinet_number_id
		cabinet_number_ids := model.GenerateCabinetID(cabinetNumbers, idc_ids)

		//生成server_id
		server_ids := model.GenerateServerID(hostNames)
		//插入对应id
		for k, _ := range server_ids {
			if len(idcNames)-1 < k {
				idcNames = append(idcNames, idcNames[0])
			}
			if len(citys)-1 < k {
				citys = append(citys, citys[0])
			}
			if len(cabinetNumbers)-1 < k {
				cabinetNumbers = append(cabinetNumbers, cabinetNumbers[0])
			}
			if len(idc_ids)-1 < k {
				idc_ids = append(idc_ids, idc_ids[0])
			}
			model.UpdateIdcID(idcNames[k], citys[k], idc_ids[k], cabinet_number_ids[k], IDCID[k])
			model.UpdateCabinetID(cabinetNumbers[k], idc_ids[k], cabinet_number_ids[k], CabinetID[k])
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

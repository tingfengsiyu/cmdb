package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"fmt"
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
	var cabinet_number_ids, idc_ids, server_ids []int
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
		code = model.BatchCreateServer2(&assets.Asset.Servers)
		//检查不存在后执行
		number := 0

		//生成idc_id
		for k, v := range idcNames {
			idc_id, _ := model.Check_Idc_Name(v)
			if idc_id == 0 {
				if k == 0 {
					id := model.LastIdcID()
					number = id + 1
				} else {
					number = number + 1
				}
				idc_ids = append(idc_ids, number)
			} else {
				idc_ids = append(idc_ids, idc_id)
			}
		}

		//生成cabinet_number_id
		number = 0
		for k, v := range cabinetNumbers {
			if len(idc_ids)-1 < k {
				idc_ids = append(idc_ids, idc_ids[0])
			}
			cabinet_number_id, _ := model.Check_Cabinet_Number(v, idc_ids[k])
			if cabinet_number_id == 0 {
				if k == 0 {
					id := model.LastCabintID()
					number = id + 1
				} else {
					number = number + 1
				}
				cabinet_number_ids = append(cabinet_number_ids, number)
			} else {
				cabinet_number_ids = append(cabinet_number_ids, cabinet_number_id)
			}
		}

		//server_id
		number = 0
		for k, v := range hostNames {
			server_id, _ := model.CheckServer(v)
			if server_id == 0 {
				if k == 0 {
					id := model.LastServeID()
					fmt.Println(id, number)
					number = id + 1
				} else {
					number = number + 1
				}
				server_ids = append(server_ids, number)
			} else {
				server_ids = append(server_ids, server_id)
			}

		}
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
			model.InsertID(idcNames[k], citys[k], cabinetNumbers[k], hostNames[k], idc_ids[k], server_ids[k], cabinet_number_ids[k])
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
	servers := []model.Server{}
	_ = c.ShouldBindJSON(&servers)
	var IDS = make([]int, 0)
	for _, v := range servers {
		IDS = append(IDS, v.ID)
	}
	code := model.BatchCheckServerID(IDS)
	if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	} else if code == errmsg.ERROR_ALL_DEVICE_EXIST {
		code = model.BatchUpdateServer(&servers)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
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

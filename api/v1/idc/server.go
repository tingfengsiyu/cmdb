package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 添加服务器
func AddServer(c *gin.Context) {
	var data model.Server
	_ = c.ShouldBindJSON(&data)
	code := model.CheckServer(data.Name)
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
	//server := model.Servers{}
	//_ = c.ShouldBindJSON(&server)
	servers := []model.Server{}
	_ = c.ShouldBindJSON(&servers)
	//var serverNames = make([]model.Names,0)
	var serverNames = make([]string,0)
	for _,v := range servers {
		//serverNames = append(serverNames, model.Names{v.Name})
		serverNames = append(serverNames, v.Name)
	}
	code := model.BatchCheckServer(serverNames)
	if code == errmsg.SUCCSE {
		//servers.BatchCreateServer()
		code = model.BatchCreateServer2(&servers)
	}else if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	}else if code == errmsg.ERROR_ALL_DEVICE_EXIST {
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
	//fmt.Println(&data)
	//fmt.Println(id)
	code := model.EditServer(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func BatchUpdateServers(c *gin.Context) {
	servers := []model.Server{}
	_ = c.ShouldBindJSON(&servers)
	var IDS = make([]int,0)
	for _,v := range servers {
		IDS = append(IDS, v.ID)
	}
	code := model.BatchCheckServerID(IDS)
	if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	}else if code == errmsg.ERROR_ALL_DEVICE_EXIST {
		code = model.BatchUpdateServer(&servers)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func DeleteServer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	code := model.DeleteServer(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func GetIdcServers(c *gin.Context){
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

	data, total := model.GetIdcServers(pageSize, pageNum,idcname)
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

func GetCabinetServers(c *gin.Context){
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

	data, total := model.GetCabinetServers(pageSize, pageNum,idcname,cabinet_number)
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

func GetUser(c *gin.Context){
	data, total := model.GetOwnedUser()
	//type User struct {
	//	User  string `json:"user"`
	//}
	//var Users = make([]User,0)
	var Users = make([]string,0)
	for _,v :=range data{
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

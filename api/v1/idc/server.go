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

func BatchAddServer(c *gin.Context) {
	server := model.Servers{}
	_ = c.ShouldBindJSON(&server)
	servers := []model.Server{}
	_ = c.ShouldBindJSON(&servers)

	//code := model.BatchCheckServer()
	code := 000
	if code == errmsg.SUCCSE {
		server.BatchCreateServer()
	}
	if code == errmsg.ERROR_DEVICE_EXIST {
		code = errmsg.ERROR_DEVICE_EXIST
	}
	code = model.BatchCreateServer2(servers)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func GetServerInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := model.GetUser(id)
	maps["username"] = data.Username
	maps["role"] = data.Role
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

func UpdateServers(c *gin.Context) {
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
	data := []model.Server{}
	_ = c.ShouldBindJSON(&data)
	code := model.BatchUpdateServer(data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func DeleteServers(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	code := model.DeleteServer(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

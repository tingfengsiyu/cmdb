package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "hello",
		})
}

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

	data, total := model.GetServer(pageSize, pageNum)
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

func DeleteServers(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	//fmt.Println(id)
	code := model.DeleteServer(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

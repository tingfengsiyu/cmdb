package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateIdc(c *gin.Context) {
	var data model.Idc
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code := model.EditIdc(id, &data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
func DeleteIdc(c *gin.Context) {
	var data model.Idc
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteIDC(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}
func GetIDCs(c *gin.Context) {
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

	data, total := model.GetIDCs(pageSize, pageNum)
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

func Network_topology(c *gin.Context) {
	//var data []model.Server
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	cabinet_number := c.Param("cabinet_number")
	user := c.Param("user")
	data, code := model.Network_topology(id, name, cabinet_number, user)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

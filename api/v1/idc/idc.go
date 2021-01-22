package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 添加IDC
func AddIdc(c *gin.Context) {
	var data model.Idc
	_ = c.ShouldBindJSON(&data)
	code := model.CreateIdc(&data)
	if code == errmsg.SUCCSE {
		model.CreateIdc(&data)
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
func UpdateIdc(c *gin.Context) {
	var data model.Idc
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	fmt.Println(&data)
	fmt.Println(id)
	code := model.EditIdc(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

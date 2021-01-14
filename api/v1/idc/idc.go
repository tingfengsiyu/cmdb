package idc

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
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
package user

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"cmdb/utils/validator"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加用户
func AddTermUser(c *gin.Context) {
	var data model.TermUser
	var msg string
	_ = c.ShouldBindJSON(&data)
	fmt.Println(data)
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": msg,
			},
		)
		c.Abort()
	}

	code = model.CheckTermUser(data.Username)
	if code == errmsg.SUCCSE {
		model.CreateTermUser(&data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		code = errmsg.ERROR_USERNAME_USED
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 查询单个用户
func GetTermUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetTermUser(id)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   1,
			"message": errmsg.GetErrMsg(code),
		},
	)

}

// 查询用户列表
func GetTermUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetTermUsers(username, pageSize, pageNum)

	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 编辑用户
func EditTermUser(c *gin.Context) {
	var data model.TermUser
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code = model.CheckUpTermUser(id, data.Username)
	if code == errmsg.SUCCSE {
		model.EditTermUser(id, &data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		c.Abort()
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 删除用户
func DeleteTermUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code = model.DeleteTermUser(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

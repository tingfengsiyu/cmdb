package user

import (
	"cmdb/middleware"
	"cmdb/model"
	"cmdb/utils/errmsg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var token string
	var code int

	formData, code = model.CheckLogin(formData.Username, formData.Password)

	if code == errmsg.SUCCSE {
		token, code = middleware.SetToken(formData.Username, formData.Role)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    formData.Username,
		"id":      formData.ID,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}

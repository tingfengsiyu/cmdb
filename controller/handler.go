package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"instances/cloud"
	"instances/db"
	Auth "instances/middleware"
	"instances/middleware/config"
	"net/http"
	"strconv"
)


func IndexHandler(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{
		"msg":"data",
	})
}

func SyncDB(c *gin.Context){
	cloud.SyncAwsECS()
	c.JSON(http.StatusOK,gin.H{
		"msg":"ok",
	})
}
func EcsListAllHandler(c *gin.Context){
	var (
		pageNum int
		pageSize int
	)
	tmpPageNum, _ := strconv.ParseInt(c.DefaultQuery("pagenum","0"),10,64)
	pageNum = int(tmpPageNum)
	tmpPageSize ,_ :=strconv.ParseInt(c.DefaultQuery("pageSize","25"),10,64)
	pageSize =int(tmpPageSize)
	fmt.Println(pageNum,pageSize)

	data,err := db.QueryAllEcs()
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg": err,
		})
		return
	}else {
		c.JSON(http.StatusOK,gin.H{
			"code":200,
			"msg": data,
		})
	}

}

func AuthHandler(c *gin.Context){
	// 用户发送用户名和密码过来
	var user config.UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	if user.Username == "q1mi" && user.Password == "q1mi123" {
		// 生成Token
		tokenString, _ := Auth.GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

func HomeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	fmt.Println(c.Get("username"))
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}
package router

import (
	"cmdb/utils"

	"github.com/gin-gonic/gin"

	"cmdb/api/v1/idc"
	"cmdb/api/v1/user"
	"cmdb/middleware"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(middleware.Log())
	r.Use(gin.Recovery())
	router := r.Group("api/v1/idc")
	router.Use(middleware.JwtToken())
	{
		router.GET("/", idc.Hello)
		router.POST("createserver", idc.AddServer)
		router.POST("createidc", idc.AddIdc)
		router.GET("getservers", idc.GetServers)
		router.DELETE("deleteservers/:id", idc.DeleteServers)
		router.PUT("editservers/:id", idc.UpdateServers)
		router.PUT("editidc/:id", idc.UpdateIdc)
	}
	cloud := r.Group("api/v1/cloud")
	{
		cloud.GET("/", idc.Hello)
	}
	auser := r.Group("api/v1/user")
	{
		auser.POST("adduser", user.AddUser)
		auser.GET("getusers", user.GetUsers)
		auser.PUT("edituser/:id", user.EditUser)
		auser.DELETE("deleteuser/:id", user.DeleteUser)
		auser.POST("login", user.Login)
	}
	_ = r.Run(utils.HttpPort)
}

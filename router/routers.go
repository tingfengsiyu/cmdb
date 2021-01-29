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
		router.POST("createserver", idc.AddServer)
		router.POST("batchcreateserver", idc.BatchAddServer)
		router.POST("batchupdateserver", idc.BatchUpdateServers)
		router.POST("createidc", idc.AddIdc)
		router.GET("getservers", idc.GetServers)
		router.GET("getserver", idc.GetServerInfo)
		router.DELETE("deleteservers/:id", idc.DeleteServers)
		router.PUT("editservers/:id", idc.UpdateServers)
		router.PUT("editidc/:id", idc.UpdateIdc)
	}
	cloud := r.Group("api/v1/cloud")
	cloud.Use(middleware.JwtToken())
	{
		cloud.GET("/", idc.GetServers)
	}
	u := r.Group("api/v1/user")
	router.Use(middleware.JwtToken())
	{
		u.POST("adduser", user.AddUser)
		u.GET("getusers", user.GetUsers)
		u.GET("getuser/:id", user.GetUserInfo)
		u.PUT("edituser/:id", user.EditUser)
		u.DELETE("deleteuser/:id", user.DeleteUser)
		//u.POST("login", user.Login)
	}
	//k := r.Group("api/v1/k8s")
	//{
	//	//k.GET("listpod",k8s.InitConfig)
	//}
	r.GET("api/v1/user/login", user.Login)
	_ = r.Run(utils.HttpPort)
}

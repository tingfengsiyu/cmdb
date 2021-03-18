package router

import (
	"cmdb/api/v1/prometheus"
	"cmdb/utils"

	"github.com/gin-gonic/gin"

	"cmdb/api/v1/cloud"
	"cmdb/api/v1/idc"
	"cmdb/api/v1/script"
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
		router.GET("getidcs", idc.GetIDCs)
		router.PUT("editidc/:id", idc.UpdateIdc)
		router.DELETE("deleteidc/:id", idc.DeleteIdc)

		router.POST("createserver", idc.AddServer)
		router.POST("batchcreateserver", idc.BatchAddServers)
		router.POST("batchupdateserver", idc.BatchUpdateServers)
		router.POST("cron", idc.Cron)
		router.POST("osinit", script.OsInit)
		router.POST("updatehostname", script.UpdateHostName)

		router.GET("WritePrometheus", prometheus.WritePrometheus)
		router.GET("getservers", idc.GetServers)
		router.GET("getserver/:id", idc.GetServer)
		router.GET("getuser", idc.GetUser)
		router.GET("getidcserver", idc.GetIdcServers)
		router.GET("getcabinetserver", idc.GetCabinetServers)
		router.DELETE("deleteserver/:id", idc.DeleteServer)
		router.PUT("editservers/:id", idc.UpdateServer)
		router.GET("getnetwork_topology", idc.Network_topology)

		router.POST("uploadexcel", idc.UploadExcel)
		router.GET("exportcsv", idc.ExportCsv)

	}
	clouds := r.Group("api/v1/cloud")
	clouds.Use(middleware.JwtToken())
	{
		clouds.GET("/sync", cloud.Sync)
	}
	u := r.Group("api/v1/user")
	router.Use(middleware.JwtToken())
	{
		u.POST("adduser", user.AddUser)
		u.GET("getusers", user.GetUsers)
		u.GET("getuser/:id", user.GetUserInfo)
		u.PUT("edituser/:id", user.EditUser)
		u.DELETE("deleteuser/:id", user.DeleteUser)
	}
	//k := r.Group("api/v1/k8s")
	//{
	//	//k.GET("listpod",k8s.InitConfig)
	//}
	r.GET("api/v1/user/login", user.Login)
	_ = r.Run(utils.HttpPort)
}

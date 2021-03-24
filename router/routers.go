package router

import (
	"cmdb/api/v1/cloud"
	"cmdb/api/v1/idc"
	"cmdb/api/v1/script"
	"cmdb/api/v1/user"
	"cmdb/middleware"
	"cmdb/utils"
	"github.com/gin-gonic/gin"
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

		router.POST("osinit", script.OsInit)
		router.POST("shellosinit", script.ShellInit)
		router.POST("storagemount", script.StorageMount)
		router.PUT("batchip", script.BatchIp)
		router.POST("updatehostname", script.UpdateHostName)
		router.POST("writeprometheus", script.WritePrometheus)
		router.GET("installmointoragent", script.InstallMointorAgent)
		router.GET("generateansiblehosts", script.GenerateAnsibleHosts)

		router.GET("getservers", idc.GetServers)
		router.GET("getserver/:id", idc.GetServer)
		router.GET("getuser", idc.GetUser)
		router.GET("getidcserver", idc.GetIdcServers)
		router.GET("getcabinetserver", idc.GetCabinetServers)
		router.DELETE("deleteserver/:id", idc.DeleteServer)
		router.PUT("editservers/:id", idc.UpdateServer)
		router.GET("getnetworktopology", idc.Networktopology)

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

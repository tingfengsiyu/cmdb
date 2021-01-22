package router

import (
	"cmdb/utils"
	"github.com/gin-gonic/gin"
	//"cmdb/api/v1/cloud"
	"cmdb/api/v1/idc"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	//auth :=  r.Group("api/v1")
	router := r.Group("api/v1/idc")
	{
		router.GET("/", idc.Hello)
		router.POST("createserver", idc.AddServer)
		router.POST("createidc", idc.AddIdc)
		router.GET("getservers", idc.GetServers)
		router.DELETE("deleteservers/:id", idc.DeleteServers)
		router.PUT("editservers/:id", idc.UpdateServers)
		router.PUT("editidc/:id", idc.UpdateIdc)
	}
	routers := r.Group("api/v1/cloud")
	{
		routers.GET("/", idc.Hello)
	}
	_ = r.Run(utils.HttpPort)
}

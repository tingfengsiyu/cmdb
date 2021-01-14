package router

import (
	"github.com/gin-gonic/gin"
	"cmdb/utils"
	//"cmdb/api/v1/cloud"
	"cmdb/api/v1/idc"

)

func InitRouter()  {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	//auth :=  r.Group("api/v1")
	router := r.Group("api/v1")
	{
		router.GET("/",idc.Hello)
		router.POST("idc/createserver",idc.AddServer)
		router.POST("idc/createidc",idc.AddIdc)
		router.GET("idc/getservers",idc.GetServers)
	}
	_ = r.Run(utils.HttpPort)
}
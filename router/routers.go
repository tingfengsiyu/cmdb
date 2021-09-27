package router

import (
	"cmdb/api/v1/cloud"
	"cmdb/api/v1/idc"
	"cmdb/api/v1/script"
	"cmdb/api/v1/user"
	"cmdb/api/v1/ws"
	"cmdb/middleware"
	"cmdb/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()

	r.LoadHTMLGlob("web/dist/index.html")
	r.Static("static", "web/dist/static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.StaticFile("/favicon.ico", "web/dist/favicon.ico")
	r.Use(middleware.Log())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	router := r.Group("api/v1/idc")

	router.Use(middleware.JwtToken())
	{
		router.GET("getidcs", idc.GetIDCs)
		router.PUT("editidc/:id", idc.UpdateIdc)
		router.DELETE("deleteidc/:id", idc.DeleteIdc)

		router.POST("createserver", idc.AddServer)
		router.POST("batchcreateserver", idc.BatchAddServers)
		router.POST("batchupdateserver", idc.BatchUpdateServers)

		router.POST("shellosinit", script.ShellInit)
		router.POST("storagemount", script.StorageMount)
		router.PUT("batchip", script.BatchIp)
		router.PUT("updatecluster", script.UpdateCluster)
		router.PUT("execshell", script.ExecWebShell)
		router.POST("updatehostname", script.UpdateHostName)
		router.POST("writeprometheus", script.WritePrometheus)
		router.GET("prometheusalerts", script.PrometheusAlerts)
		router.GET("installmointoragent", script.InstallMointorAgent)
		router.POST("generateansiblehosts", script.GenerateAnsibleHosts)
		router.POST("generateclustershosts", script.GenerateClustersHosts)

		router.GET("getservers", idc.GetServers)
		router.GET("getclusters", idc.GetClusters)
		router.GET("getserver/:id", idc.Networktopology)
		router.GET("getcluster", idc.GetCluster)
		router.GET("getuser", idc.GetUser)
		router.GET("getidcserver", idc.GetIdcServers)
		router.GET("getcabinetserver", idc.GetCabinetServers)
		router.DELETE("deleteserver/:id", idc.DeleteServer)
		router.PUT("editservers/:id", idc.UpdateServer)
		router.GET("getnetworktopology", idc.Networktopology)

		router.GET("opsrecords", idc.Records)
		router.GET("opsrecord/:id", idc.Records)

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
	//term_user 新增终端用户和对应机器权限配置
	t := r.Group("api/v1/term")
	router.Use(middleware.JwtToken())
	{
		t.POST("adduser", user.AddTermUser)
		t.GET("getusers", user.GetTermUsers)
		t.GET("getuser/:id", user.GetTermUserInfo)
		t.PUT("edituser/:id", user.EditTermUser)
		t.GET("permissions", user.ServerPermissions)
		t.POST("addpermissions", user.AddPermissions)
		t.DELETE("deletePermission/:id", user.DeletePermission)
		t.PUT("editpermission/:id", user.EditPermission)
		t.DELETE("deleteuser/:id", user.DeleteTermUser)
		t.GET("consolelog", ws.GetTermLogs)
	}

	w := r.Group("api/v1/ws")
	w.Use(middleware.WsJwtToken())
	{
		w.GET("console/:id", ws.WsSsh)
	}
	//k := r.Group("api/v1/k8s")
	//{
	//	//k.GET("listpod",k8s.InitConfig)
	//}
	r.POST("api/v1/login", user.Login)
	_ = r.Run(utils.HttpPort)
}

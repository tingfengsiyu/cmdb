package cloud

import (
	"cmdb/middleware"
	"cmdb/model"
	"cmdb/utils"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func Sync(c *gin.Context) {
	client, err := ecs.NewClientWithAccessKey(
		"cn-beijing", // 地域ID
		utils.AccessKey,
		utils.SecretKey)
	if err != nil {
		// 异常处理
		middleware.SugarLogger.Debugf("认证错误%s", err)
	}
	// 创建API请求并设置参数
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"

	response, err := client.DescribeInstances(request)
	/* 等价于 request.PageSize = "10"
	PageNumber
	PageSize
	MaxResults
	request.PageSize = requests.NewInteger(10)
	*/
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(100)
	//request.MaxResults = requests.NewInteger(100)
	// 发起请求并处理异常
	if err != nil {
		middleware.SugarLogger.Debugf("请求错误%s", err)
	}
	var instanceStruct = make([]model.CloudInstance, 0)
	//fmt.Println(response.TotalCount)
	for _, tmp := range response.Instances.Instance {
		StartTime, _ := time.Parse("2006-01-02T15:04Z", tmp.StartTime)
		CreationTime, _ := time.Parse("2006-01-02T15:04Z", tmp.CreationTime)
		ExpiredTime, _ := time.Parse("2006-01-02T15:04Z", tmp.ExpiredTime)
		var PublicIpAddress string
		if len(tmp.PublicIpAddress.IpAddress) > 0 {
			PublicIpAddress = tmp.PublicIpAddress.IpAddress[0]
		}
		instanceStruct = append(instanceStruct, model.CloudInstance{
			Model:                  gorm.Model{},
			InstanceId:             tmp.InstanceId,
			HostName:               tmp.HostName,
			Status:                 tmp.Status,
			CPU:                    tmp.Cpu,
			Memory:                 tmp.Memory,
			OSName:                 tmp.OSName,
			RegionId:               tmp.RegionId,
			InstanceType:           tmp.InstanceType,
			OsType:                 tmp.OSType,
			InternetMaxBandwidthIn: tmp.InternetMaxBandwidthIn,
			StartTime:              StartTime.Local(),
			ExpiredTime:            ExpiredTime.Local(),
			InstanceCreationTime:   CreationTime.Local(),
			LocalStorageCapacity:   tmp.LocalStorageCapacity,
			InnerIpAddress:         tmp.NetworkInterfaces.NetworkInterface[0].PrimaryIpAddress,
			PublicIpAddress:        PublicIpAddress,
			Cloud:                  "Ali",
		})
		//fmt.Println(tmp.InstanceId,
		//	tmp.HostName,
		//	tmp.OSName,
		//	tmp.CreationTime,
		//	tmp.ExpiredTime,
		//	tmp.Cpu,
		//	tmp.Memory,
		//	tmp.LocalStorageCapacity,
		//	tmp.Status,
		//	tmp.PublicIpAddress.IpAddress[0],
		//	tmp.InnerIpAddress.IpAddress,
		//	tmp.RdmaIpAddress,
		//	tmp.RegionId,
		//	)
	}
	code := model.BatchAddAliEcs(instanceStruct)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  "ok",
	})
}

func SyncDB(c *gin.Context) {
	//cloud.SyncAwsECS()
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
func EcsListAllHandler(c *gin.Context) {
	var (
		pageNum  int
		pageSize int
	)
	tmpPageNum, _ := strconv.ParseInt(c.DefaultQuery("pagenum", "0"), 10, 64)
	pageNum = int(tmpPageNum)
	tmpPageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "25"), 10, 64)
	pageSize = int(tmpPageSize)
	fmt.Println(pageNum, pageSize)

	////data,err := db.QueryAllEcs()
	//if err != nil{
	//	c.JSON(http.StatusOK,gin.H{
	//		"code":500,
	//		"msg": err,
	//	})
	//	return
	//}else {
	//	c.JSON(http.StatusOK,gin.H{
	//		"code":200,
	//		"msg": data,
	//	})
	//}
}

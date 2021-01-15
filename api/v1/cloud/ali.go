package cloud

import (
	"cmdb/utils"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SyncAliECS() {
	ecsClient, err := ecs.NewClientWithAccessKey(
		"cn-beijing",    // 地域ID
		utils.AccessKey, // 您的Access Key ID
		utils.SecretKey) // 您的Access Key Secret
	if err != nil {
		// 异常处理
		fmt.Println("认证错误", err)
	}
	// 创建API请求并设置参数
	request := ecs.CreateDescribeInstancesRequest()
	// 等价于 request.PageSize = "10"
	/*request.Filter3Key("ExpiredStartTime")
	request.Filter3Value("INSTANCE_EXPIRED_START_TIME_IN_UTC_STRING")
	request.Filter4Key("ExpiredEndTime")
	request.Filter4Value("INSTANCE_EXPIRE_END_TIME_IN_UTC_STRING")
	request.PageSize = requests.NewInteger(10)

	*/
	// 发起请求并处理异常
	response, err := ecsClient.DescribeInstances(request)
	if err != nil {
		// 异常处理
		panic(err)
	}
	tmp := response.Instances.Instance[0]
	fmt.Println(tmp.InstanceId,
		tmp.InstanceName,
		tmp.ExpiredTime,
		tmp.PublicIpAddress,
		tmp.InnerIpAddress,
		tmp.InstanceType)
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

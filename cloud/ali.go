package cloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)


func SyncAliECS() {
	ecsClient, err := ecs.NewClientWithAccessKey(
		"cn-beijing",             // 地域ID
		"LTAIeT4TFO30pcQp",         // 您的Access Key ID
		"sn3h6qh3O7P7dHBlR7J8KlcprOp8V9")        // 您的Access Key Secret
	if err != nil {
		// 异常处理
		panic(err)
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
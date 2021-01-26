package cloud

import (
	"cmdb/model"
	_ "cmdb/model"
	"cmdb/utils/errmsg"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SyncAwsECS(c *gin.Context) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	var data model.AwsServer
	// Create new EC2 client
	svc := ec2.New(sess)
	ss := &ec2.DescribeInstancesInput{}
	// Call to get detailed information on each instance
	result, err := svc.DescribeInstances(ss)
	if err != nil {
		panic("get ec2 error")
	} else {
		for _, value := range result.Reservations {
			//fmt.Println(value)
			var Disk string
			for _, v := range value.Instances {
				for _, diskID := range v.BlockDeviceMappings {
					tmpDisk := *diskID.Ebs.VolumeId
					Disk = Disk + " " + tmpDisk
				}
				/*fmt.Println(*v.InstanceId,
				*v.Tags[0].Value,*v.LaunchTime,"0",
				*v.PublicIpAddress, *v.PrivateIpAddress,*v.ImageId,*v.InstanceType,*v.State.Name,
				"Regions",*v.Placement.AvailabilityZone,Disk,*v.Architecture,*v.NetworkInterfaces[0].OwnerId,"aws")
				*/
				//fmt.Println(*v.ImageId,*v.InstanceId)
				//imageName := GetImageName(svc, *v.ImageId)
				instanceid, _ := strconv.Atoi(c.Param("instanceid"))
				code := model.EditAws(instanceid, &data)
				fmt.Println(code)
				//sqlStr := "insert into ecs (InstanceId,Name,CreateTime,ExpiredTime,PublicIpAddress,PrivateIpAddress,OsName,InstanceType,State,Regions,AvailabilityZones,Disk,Arch,User,Cloud )values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
				//_, err = db.DB.Exec(sqlStr, *v.InstanceId,
				//	*v.Tags[0].Value, *v.LaunchTime, "0",
				//	*v.PublicIpAddress, *v.PrivateIpAddress, imageName, *v.InstanceType, *v.State.Name,
				//	"Regions", *v.Placement.AvailabilityZone, Disk, *v.Architecture, *v.NetworkInterfaces[0].OwnerId, "aws")
				if err != nil {
					fmt.Println("插入ECS实例失败")
					return
				}

			}

		}
	}
}

func GetImageName(svc *ec2.EC2, ImageId string) (imageName string) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{
			aws.String(ImageId),
		},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		imageName = ""
	}
	//fmt.Println(result)
	if len(result.Images) == 0 {
		imageName = ""
	} else {
		imageName = *result.Images[0].PlatformDetails

	}
	return imageName
}

func GetDiskName(svc *ec2.EC2) (DiskSize string) {
	input := &ec2.DiskInfo{}

	fmt.Println(input)
	return
}

func EcsListAll(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total := model.GetAws(pageSize, pageNum)
	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func EditAws(c *gin.Context){
	var data model.AwsServer
	instanceid, _ := strconv.Atoi(c.Param("instanceid"))
	_ = c.ShouldBindJSON(&data)
	fmt.Println(&data)
	fmt.Println(instanceid)
	code := model.EditAws(instanceid, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
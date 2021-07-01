package idc

import (
	"bytes"
	"cmdb/middleware"
	"cmdb/model"
	"cmdb/utils"
	"cmdb/utils/errmsg"
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strconv"
)

func UpdateIdc(c *gin.Context) {
	var data model.Idc
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code := model.EditIdc(id, &data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func DeleteIdc(c *gin.Context) {
	var data model.Idc
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteIDC(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

func GetIDCs(c *gin.Context) {
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

	data, total := model.GetIDCs(pageSize, pageNum)
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

func Networktopology(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println(c.Param("id"))
	name := c.Query("name")
	cabinet_number := c.Query("cabinet_number")
	user := c.Query("user")
	cluster := c.Query("cluster")
	private_ip_address := c.Query("private_ip_address")
	data, code := model.NetworkTopology(id, name, cabinet_number, user, cluster, private_ip_address)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   len(data),
		"message": errmsg.GetErrMsg(code),
	})
}

//判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func UploadExcel(c *gin.Context) {
	files, err := c.FormFile("file")
	if err != nil {
		c.String(400, "文件格式错误")
		return
	}
	SR_File_Max_Bytes, _ := strconv.ParseInt(utils.SR_File_Max_Bytes, 10, 64)
	if files.Size > SR_File_Max_Bytes {
		c.String(400, "文件大小超过", SR_File_Max_Bytes)
		return
	}

	// 设置文件需要保存的指定位置并设置保存的文件名字
	dst := path.Join("logs", files.Filename)
	// 上传文件到指定的路径
	a := c.SaveUploadedFile(files, dst)
	if a != nil {
		c.String(400, "文件格式错误")
		return
	}
	xlsx, err := excelize.OpenFile(dst)
	if err != nil {
		middleware.SugarLogger.Errorf("打开文件错误%s", err)
	}
	rows := xlsx.GetRows("Sheet" + "1")

	servers := []model.Server{}
	var idcNames = make([]string, 0)
	var cabinetNumbers = make([]string, 0)
	var citys = make([]string, 0)
	var hostNames = make([]string, 0)
	for key, row := range rows {
		if key > 0 {
			servers = append(servers, model.Server{Name: row[1],
				Models:           row[2],
				Location:         row[3],
				PrivateIpAddress: row[4],
				PublicIpAddress:  row[5],
				Label:            row[6],
				Cluster:          row[7],
				LabelIpAddress:   row[8],
				Cpu:              row[9],
				Memory:           row[10],
				Disk:             row[11],
				User:             row[12],
				State:            row[13],
			})
			hostNames = append(hostNames, row[1])
			citys = append(citys, row[14])
			idcNames = append(idcNames, row[15])
			cabinetNumbers = append(cabinetNumbers, row[16])
		}
	}

	code := model.BatchCheckServer(hostNames)
	code = addServerVerify(code, servers, idcNames, cabinetNumbers, hostNames, citys)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func ExportCsv(c *gin.Context) {
	bytesBuffer := &bytes.Buffer{}
	bytesBuffer.WriteString("xEFxBBxBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码
	data, _ := model.NetworkTopology(0, "", "", "", "", "")
	writer := csv.NewWriter(bytesBuffer)

	writer.Write([]string{"id", "主机名", "型号", "位置U数", "私有地址", "公网地址", "角色标签", "集群名", "机房标签ip", "cpu", "内存", "磁盘", "用户", "状态已上架", "城市", "机房名", "机柜名"})
	for _, v := range data {
		writer.Write([]string{strconv.Itoa(v.ID), v.Name, v.Models, v.Location, v.PrivateIpAddress, v.PublicIpAddress,
			v.Label, v.Cluster, v.LabelIpAddress, v.Cpu, v.Memory, v.Disk, v.User, v.State, v.City, v.IDC_Name, v.Cabinet_Number})
	}

	writer.Flush() // 此时才会将缓冲区数据写入
	// 设置下载的文件名
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	// 设置文件类型以及输出数据
	c.Data(http.StatusOK, "text/csv", bytesBuffer.Bytes())
}

func Records(c *gin.Context) {
	action := c.Query("action")
	data, err := model.GetRecords(action)
	var code int
	if err != nil {
		code = 400
	} else {
		code = 200
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status": code,
			"data":   data,
			"total":  len(data),
		},
	)
}

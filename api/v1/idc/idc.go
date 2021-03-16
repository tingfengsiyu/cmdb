package idc

import (
	"bufio"
	"cmdb/model"
	"cmdb/utils"
	"cmdb/utils/errmsg"
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"io"
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

func Network_topology(c *gin.Context) {
	//var data []model.Server
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	cabinet_number := c.Param("cabinet_number")
	user := c.Param("user")
	data, code := model.Network_topology(id, name, cabinet_number, user)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

func DownloadReadFile(c *gin.Context) {
	//http下载地址 csv
	csvFileUrl := c.PostForm("file_name")
	res, err := http.Get(csvFileUrl)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer res.Body.Close()
	//读取csv
	reader := csv.NewReader(bufio.NewReader(res.Body))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			c.String(400, err.Error())
			return
		}
		//line 就是每一行的内容
		fmt.Println(line)
		//line[0] 就是第几列
		fmt.Println(line[0])
	}
}

func DownloadWriteFile(c *gin.Context) {
	//写文件
	var filename = "./output1.csv"
	if !checkFileIsExist(filename) {
		file, err := os.Create(filename) //创建文件
		if err != nil {
			c.String(400, err.Error())
			return
		}
		buf := bufio.NewWriter(file) //创建新的 Writer 对象
		buf.WriteString("\xEF\xBB\xBF")
		buf.Flush()
		defer file.Close()
	}
	//返回文件流
	c.File(filename)
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
		fmt.Println(err)
		os.Exit(1)
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

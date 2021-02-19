package model

import (
	"cmdb/middleware"
	"cmdb/utils/errmsg"
	"fmt"
	"gorm.io/gorm"
)

// 查询IDC是否存在
func CheckIdc(idc_name string) (code int) {
	var idc Idc
	db.Select("id").Where("idc_name = ?", idc_name).First(&idc)
	if idc.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST //2001
	}
	return errmsg.SUCCSE
}

func Check_Cabinet_Number(cabinet_number string, idc_id int) (int, int) {

	var data = Cabinet{}
	err := db.Unscoped().Debug().Select("cabinet_number_id").Where("idc_id = ? and cabinet_number = ? ", idc_id, cabinet_number).First(&data).Error
	if err != nil {
		middleware.SugarLogger.Errorf("sql查询错误%s", err)
		return int(data.ID), errmsg.ERROR
	}

	return data.Cabinet_NumberID, errmsg.SUCCSE
}

//检查idc name是否存在
func Check_Idc_Name(idc_name string) (int, int) {
	var data = Idc{}
	err := db.Unscoped().Debug().Select("idc_id").Where("idc_name = ?", idc_name).First(&data).Error
	fmt.Println(data.IDC_ID, data.ID)
	if err != nil {
		middleware.SugarLogger.Errorf("sql查询错误%s", err)
		return data.IDC_ID, errmsg.ERROR
	}
	return data.IDC_ID, errmsg.SUCCSE
}

func CreateIdc(data *Idc) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func EditIdc(id int, data *Idc) int {
	var idc Idc
	var maps = make(map[string]interface{})
	maps["name"] = data.IDC_Name
	maps["city"] = data.City
	err = db.Model(&idc).Where("id=?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetIDC(name string) int64 {
	var idc Server
	db.Select("id").Where("name = ?", name).First(&idc)
	if idc.ID > 0 {
		return errmsg.ERROR_DEVICE_EXIST //2001
	}
	return errmsg.SUCCSE
}

func GetIDCs(pageSize int, pageNum int) ([]Idc, int64) {
	var svc []Idc
	var total int64
	err = db.Unscoped().Find(&svc).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&svc).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return svc, total
}

type result struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

//网络拓扑展示
//根据机房名查询name和机柜号查询对应服务器名中对应的idc和机柜号和所属用户，形成网络拓扑
func Network_topology(id int, name, cabinet_number, user string) ([]result, int) {
	var svc []result
	errs := db.Debug().Unscoped().Model(&Server{}).Select("server.name,label,ipaddress,server.cabinet_number").Joins("left join idc on server.idc=idc.name").Scan(&svc)
	//select server.name,label,ipaddress ,idc.name  from server  left join idc  on  server.idc=idc.name;
	// 多连接及参数
	middleware.SugarLogger.Errorf("查询错误%s", errs)
	return svc, errmsg.SUCCSE
}

func DeleteIDC(id int) int {
	var idc Idc
	err = db.Debug().Unscoped().Where("id = ? ", id).Delete(&idc).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

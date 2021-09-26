package model

import (
	"cmdb/middleware"
	"cmdb/utils/errmsg"
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

func CheckCabinetNumber(cabinet_number string, idc_id int) (int, int) {

	var data = Cabinet{}
	err := db.Select("cabinet_number_id").Where("idc_id = ? and cabinet_number = ? ", idc_id, cabinet_number).First(&data).Error
	if err != nil {
		middleware.SugarLogger.Errorf("sql查询错误%s", err)
		return int(data.ID), errmsg.ERROR
	}

	return data.Cabinet_NumberID, errmsg.SUCCSE
}

//检查idc name是否存在
func CheckIdcName(idc_name string) (int, int) {
	var data = Idc{}
	err := db.Select("idc_id").Where("idc_name = ?", idc_name).First(&data).Error
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

//网络拓扑展示
//根据机房名查询name和机柜号查询对应服务器名中对应的idc和机柜号和所属用户，形成网络拓扑
func NetworkTopology(id int, name, cabinet_number, user, cluster, private_ip_address string) ([]ScanServers, int) {
	var scan []ScanServers
	if id != 0 {
		err = db.Raw("select distinct server.id, city,idc_name,cabinet_number,name,models,location,private_ip_address,public_ip_address,label,cluster,label_ip_address,cpu,"+
			"memory,disk,user,state,gpu,mac from  server  left join cabinet on  cabinet.cabinet_number_id=server.cabinet_number_id "+
			"left join idc on idc.idc_id =server.idc_id where name like ? and private_ip_address like ? and  cabinet_number like ? and  user like ? and  cluster  like ?  and   server.id=?", "%"+name+"%", "%"+private_ip_address+"%", "%"+cabinet_number+"%", "%"+user+"%", cluster+"%", id).Scan(&scan).Error

	} else {
		err = db.Raw("select distinct server.id, city,idc_name,cabinet_number,name,models,location,private_ip_address,public_ip_address,label,cluster,label_ip_address,cpu,"+
			"memory,disk,user,state,gpu,mac from  server  left join cabinet on  cabinet.cabinet_number_id=server.cabinet_number_id "+
			"left join idc on idc.idc_id =server.idc_id where name like ? and private_ip_address like ? and  cabinet_number like ? and  user like ? and  cluster  like ?  ", "%"+name+"%", "%"+private_ip_address+"%", "%"+cabinet_number+"%", "%"+user+"%", cluster+"%").Scan(&scan).Error

	}
	if err != nil {
		middleware.SugarLogger.Errorf("查询错误%s", err)
	}
	return scan, errmsg.SUCCSE
}

func DeleteIDC(id int) int {
	var idc Idc
	err = db.Debug().Unscoped().Where("id = ? ", id).Delete(&idc).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GenerateIDCID(idcNames []string) []int {
	var nameMap = make(map[string]interface{})
	number := 0
	var idcIds = make([]int, 0)
	for k, v := range idcNames {
		idcId, _ := CheckIdcName(v)
		if idcId == 0 {
			if nameMap[v] == nil {
				if k == 0 {
					//第一次查询数据库总记录
					id := LastIdcID()
					number = id + 1
				} else if k != 0 {
					number = number + 1
				}
				nameMap[v] = number
				idcIds = append(idcIds, number)
			} else if nameMap[v] != nil {
				number = nameMap[v].(int)
			}
			idcIds = append(idcIds, number)
		} else {
			idcIds = append(idcIds, idcId)
		}
	}
	return idcIds
}

func GenerateCabinetID(cabinetNumbers []string, idcIds []int) []int {
	var numberMap = make(map[string]interface{})
	number := 0
	var cabinetNumberIds = make([]int, 0)
	for k, v := range cabinetNumbers {
		if len(idcIds)-1 < k {
			idcIds = append(idcIds, idcIds[0])
		}
		cabinetNumberId, _ := CheckCabinetNumber(v, idcIds[k])
		if cabinetNumberId == 0 {
			if numberMap[v] == nil {
				if k == 0 {
					//第一次查询数据库总记录
					id := LastCabintID()
					number = id + 1
				} else if k != 0 {
					number = number + 1
				}
				numberMap[v] = number
				cabinetNumberIds = append(cabinetNumberIds, number)
			} else if numberMap[v] != nil {
				t := numberMap[v].(int)
				cabinetNumberIds = append(cabinetNumberIds, t)
			}

		} else {
			cabinetNumberIds = append(cabinetNumberIds, cabinetNumberId)
		}
	}
	return cabinetNumberIds
}

func InsertIdcID(idcName, cityName string, idcId, cabinetNumberId int) {
	var idc = Idc{}
	var idcs = make(map[string]interface{})
	idcs["cabinet_number_id"] = cabinetNumberId
	idcs["idc_id"] = idcId
	idcs["city"] = cityName
	idcs["idc_name"] = idcName
	////检查 cabinet_number_id 是否已存在，存在相同则不创建
	err = db.Where("idc_name = ? AND  cabinet_number_id = ?", idcName, cabinetNumberId).Find(&idc).Error
	//err = db.Unscoped().Debug().Where("idc_name = ? AND  cabinet_number_id = ?", idc_name, cabinet_number_id).Find(&idc).Error
	if idc.Cabinet_NumberID == 0 {
		//idc_id, _ := Check_Idc_Name(idc_name)
		//if idc_id == 0 {
		db.Model(&idc).Create(idcs)
		//}
	}
}

func InsertCabinetID(cabinetNumber string, idcId, cabinetNumberId int) {
	var cabinet = Cabinet{}
	var cabinets = make(map[string]interface{})
	cabinets["cabinet_number_id"] = cabinetNumberId
	cabinets["idc_id"] = idcId
	cabinets["cabinet_number"] = cabinetNumber

	//err = db.Where("idc_id = ? AND  cabinet_number_id = ?", idc_id, cabinet_number_id).Find(&cabinet).Error
	tmpcabinetNumberId, _ := CheckCabinetNumber(cabinetNumber, idcId)
	if tmpcabinetNumberId == 0 {
		//if cabinet.Cabinet_NumberID == 0 {
		//cabinet_number_id, _ := Check_Cabinet_Number(cabinet_number, idc_id)
		//if cabinet_number_id == 0 {
		db.Model(&cabinet).Create(cabinets)
		//}
	}

}

func GetRecords(action string) ([]OpsRecords, error) {
	var records []OpsRecords
	if action != "" {
		err = db.Where("action LIKE ?", action+"%").Find(&records).Error
	} else {
		err = db.Find(&records).Error
	}
	return records, err
}

func InsertRecords(records OpsRecords) int {
	if err := db.Create(&records).Error; err != nil {
		middleware.SugarLogger.Errorf("插入错误", err)
	}

	db.Select("id").Where("object = ?", records.Object).First(&records)
	return int(records.ID)
}

func UpdateRecords(id, state int, success, errors string) {
	err := db.Exec("UPDATE ops_records SET error=?,state=?,success=? where id=?", errors, state, success, id).Error
	if err != nil {
		middleware.SugarLogger.Errorf("sql update error", err)
	}
}

/*
func UpdateIdcID(idc_name, city_name string, idc_id, cabinet_number_id, ID int) {
	var idc = Idc{}
	var idcs = make(map[string]interface{})
	idcs["cabinet_number_id"] = cabinet_number_id
	idcs["idc_id"] = idc_id
	idcs["city"] = city_name
	idcs["idc_name"] = idc_name
	//检查 cabinet_number_id 是否已存在，存在相同则不创建
	err = db.Unscoped().Debug().Where("idc_name = ? AND  cabinet_number_id = ?", idc_name, cabinet_number_id).Find(&idc).Error
	if err != nil {
		middleware.SugarLogger.Errorf("查询idc错误%s", err)
	}
	if idc.Cabinet_NumberID == 0 {
		db.Model(&idc).Where("idc_id =?", ID).Updates(idcs)
	}
}

func UpdateCabinetID(cabinet_number string, idc_id, cabinet_number_id, ID int) {
	var cabinet = Cabinet{}
	var cabinets = make(map[string]interface{})
	cabinets["cabinet_number_id"] = cabinet_number_id
	cabinets["idc_id"] = idc_id
	cabinets["cabinet_number"] = cabinet_number

	err = db.Where("idc_id = ? AND  cabinet_number_id = ?", idc_id, cabinet_number_id).Find(&cabinet).Error
	if err != nil {
		middleware.SugarLogger.Errorf("查询cabinet错误%s", err)
	}
	if cabinet.Cabinet_NumberID == 0 {
		db.Model(&cabinet).Where("cabinet_number_id=?", ID).Updates(cabinets)
	}

}
*/

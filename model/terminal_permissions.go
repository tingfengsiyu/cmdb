package model

import (
	"cmdb/utils/errmsg"
)

func Permissions(server_id int, username string) (ScanTerm, int) {
	var scan []ScanTerm
	err := db.Debug().Raw(`select p.group,s.server_id, s.name,private_ip_address,u.username as user,t.username,t.password,t.protocol,t.port 
from user as u ,user_permissions  as p ,term_user as t ,server as s 
where u.id = p.user_id and s.server_id =p.server_id and p.term_user_id=t.id and  
s.id= ?  and u.username= ?`, server_id, username).Scan(&scan).Error
	if err != nil {
		return scan[0], errmsg.ERROR
	}
	return scan[0], errmsg.SUCCSE
}

func AllPermissions(id int, name, ip string) ([]ScanTerm, int) {
	var scan []ScanTerm
	var err error
	if id != 0 {
		err = db.Raw(`select p.group,p.id,s.server_id , s.name,private_ip_address,u.username as user,t.username,t.password,t.protocol,t.port 
from user as u ,user_permissions  as p ,term_user as t ,server as s where u.id = p.user_id and s.server_id =p.server_id and p.term_user_id=t.id and  s.id = ? and  s.name like   ? and private_ip_address  like ?`, id, name, ip).Scan(&scan).Error

	} else {
		err = db.Raw(`select p.group,p.id,s.server_id , s.name,private_ip_address,u.username as user,t.username,t.password,t.protocol,t.port 
from user as u ,user_permissions  as p ,term_user as t ,server as s where u.id = p.user_id and s.server_id =p.server_id and p.term_user_id=t.id and   s.name like   ? and private_ip_address  like ?`, "%"+name+"%", "%"+ip+"%").Scan(&scan).Error

	}

	if err != nil {
		return scan, errmsg.ERROR
	}
	return scan, errmsg.SUCCSE
}

func CreatePermissions(data *[]UserPermissions) int {
	err := db.Debug().Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func EditPermission(data *UserPermissions) {

}

func DeletePermission(id int) int {
	var p UserPermissions
	err := db.Where("id=?", id).Delete(&p).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteAllPermission(data *UserPermissions) {

}

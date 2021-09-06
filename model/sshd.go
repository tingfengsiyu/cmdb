package model

import "cmdb/middleware"

func Sshd() *[]ScanSshd {
	var svc []ScanSshd
	err := db.Model(&Server{}).Select("ssh_user.server_id,private_ip_address,name,ssh_username,ssh_password,port").Joins("inner  join " +
		"ssh_user on server.server_id=ssh_user.server_id").Limit(5).Scan(&svc).Error
	if err != nil {
		middleware.SugarLogger.Errorf("", err)
		return &svc
	}
	return &svc
}

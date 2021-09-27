package model

import (
	logger "cmdb/middleware"
)

func TermlogCreate(data TermLog) {
	err := db.Create(&data).Error
	if err != nil {
		logger.SugarLogger.Error(err)
	}
}

func GetTernLogs(user, ip string, pageSize int, pageNum int) ([]TermLog, int64) {
	var users []TermLog
	var total int64
	if user != "" {
		db.Select("id,user,term_user,updated_at,protocol,client_ip,log,private_ip_address").Where(
			"user LIKE ? and private_ip_address = ? ", user+"%", "%"+ip+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		db.Model(&users).Where(
			"username LIKE ?", user+"%",
		).Count(&total)
		return users, total
	}
	db.Select("id,user,term_user,updated_at,protocol,client_ip,log,private_ip_address").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	db.Model(&users).Count(&total)
	return users, total
}

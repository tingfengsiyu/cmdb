package model

import (
	logger "cmdb/middleware"
)

func SshlogCreate(data SshLog) {
	err := db.Create(&data).Error
	if err != nil {
		logger.SugarLogger.Error(err)
	}
}

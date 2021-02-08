package model

import (
	"cmdb/middleware"
	"cmdb/utils/errmsg"
)

func BatchAddAliEcs(servers []CloudInstance) int {
	err := db.Debug().Create(&servers).Error
	if err != nil {
		middleware.SugarLogger.Debugf("批量插入错误%s", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

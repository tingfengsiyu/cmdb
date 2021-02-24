package prometheus

import (
	"cmdb/model"
	"github.com/gin-gonic/gin"
)

func Listtargets() {

}

func Checktargets() {

}

func Updatetargets() {

}

func Deletetargets() {

}

func CheckAgentStatus(c *gin.Context) {
	go model.CheckAgentStatus()
}

package prometheus

import (
	"cmdb/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Checktargets() {

}

func Updatetargets() {

}

func Deletetargets() {

}

func WritePrometheus(c *gin.Context) {
	model.WritePrometheus()
	c.JSON(
		http.StatusOK, gin.H{
			"status":  200,
			"message": "write ok!!!",
		},
	)
}

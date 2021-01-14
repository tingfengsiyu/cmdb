package idc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(c *gin.Context){
	c.JSON(
		http.StatusOK,gin.H{
			"status": 200,
			"message": "hello",
		},)
}
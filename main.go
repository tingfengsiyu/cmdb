package main

import (
	"cmdb/middleware"
	"cmdb/model"
	"cmdb/router"
)

func main() {
	model.InitDb()
	err := middleware.InitLog()
	if err != nil {
		panic(err)
	}
	model.Croninit()
	router.InitRouter()
}

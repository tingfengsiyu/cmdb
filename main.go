package main

import (
	"cmdb/router"
	"cmdb/model"
)

func main() {
	model.InitDb()
	router.InitRouter()
}

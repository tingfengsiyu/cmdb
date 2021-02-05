package main

import (
	"fmt"
	"io/ioutil"
)

type PrometheusTarget struct {
	Targets []string `json:"targets"`
	Labels  string   `json:"labels"`
}

func main() {
	//var a = make([]PrometheusTarget, 0)
	// var Targets = make([]string, 0)
	// Targets = append(Targets, "172.16.0.1:9100")
	// a = append(a, PrometheusTarget{Targets: Targets,
	// 	Labels: "成都郫县",
	// })
	//fmt.Println(a)
	//a := & [] PrometheusTarget{}
	data, err := ioutil.ReadFile("logs/process.json")
	//errs := json.Unmarshal(data,&a)
	fmt.Println(data, err)
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PrometheusTarget struct {
	Targets []string `json:"targets"`
	Labels  Labels   `json:"labels"`
}

type Labels struct {
	Idc string `json:"idc"`
}

func main() {
	//decode

	ReadJsonfile()
	//WriteJsonfile()

}
func ReadJsonfile() {
	filePtr, err := os.Open("logs/node.json")
	var read = make([]PrometheusTarget, 0)
	if err != nil {
		fmt.Println("Open file failed [Err:%s]", err.Error())
		return
	}
	defer filePtr.Close()
	decode := json.NewDecoder(filePtr)
	err = decode.Decode(&read)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())
	} else {
		fmt.Println("Decoder success")
		fmt.Println(read)
	}
}
func WriteJsonfile() {
	var Targets = make([]string, 0)
	var b = make([]PrometheusTarget, 0)
	var label = Labels{}

	Targets = append(Targets, "172.16.0.1:9100")
	label = Labels{
		Idc: "This is  a  tests",
	}
	b = append(b, PrometheusTarget{Targets: Targets,
		Labels: label,
	})
	filePtr, err := os.Create("person_info2.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
	filePtr.Write(data)
}

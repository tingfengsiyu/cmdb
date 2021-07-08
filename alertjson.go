package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type AlertStatus struct {
	Status string `json:"status"`
	Data  string   `json:"data"`
	//Alert
}

type Alert struct {
	Alerts []Labels
}
type Labels struct {
	Alertname string `json:"alertname"`
	Instance string `json:"instance"`
}

func main() {
	//var read = make([]AlertStatus, 0)
	var read = AlertStatus{}
	filePtr, err := os.Open("b.json")
	if err != nil {
		fmt.Println("Open file failed [Err:%s]")
	}
	defer filePtr.Close()
	decode := json.NewDecoder(filePtr)
	err = decode.Decode(&read)
	if err != nil {
	fmt.Println("Decoder failed")
	}
}
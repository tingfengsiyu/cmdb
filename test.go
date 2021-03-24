package main

import (
	"fmt"
	"os"
)

type ansibleStruct struct {
	PrivateIpAddress string `json:"private_ip_address"`
	Label            string `json:"label"`
	Cluster          string `json:"cluster"`
}

type ansible2Struct struct {
	PrivateIpAddress string `json:"private_ip_address"`
	Group            string `json:"group"`
}

func main() {
	var a = make([]ansibleStruct, 0)
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.1", Label: "lotus-worker", Cluster: "chengdu-1"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.3", Label: "lotus-worker", Cluster: "chengdu-1"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.2", Label: "lotus-miner", Cluster: "chengdu-1"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.3", Label: "lotus-storage", Cluster: "chengdu-1"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.4", Label: "lotus-miner", Cluster: "chengdu-2"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.4", Label: "lotus-miner", Cluster: "chengdu-2"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.5", Label: "lotus-worker", Cluster: "chengdu-2"})
	a = append(a, ansibleStruct{PrivateIpAddress: "172.0.0.6", Label: "lotus", Cluster: "chengdu-2"})
	var worker, miner, storage, none []string
	var maps = make(map[string][]string, 0)
	for _, v := range a {
		if maps[v.Cluster+"-"+v.Label] == nil {
			miner = []string{}
			worker = []string{}
			storage = []string{}
			none = []string{}
		}
		if v.Label == "lotus-worker" {
			worker = append(worker, v.PrivateIpAddress)
			maps[v.Cluster+"-"+v.Label] = worker
		} else if v.Label == "lotus-storage" {
			storage = append(storage, v.PrivateIpAddress)
			maps[v.Cluster+"-"+v.Label] = storage
		} else if v.Label == "lotus-miner" {
			miner = append(miner, v.PrivateIpAddress)
			maps[v.Cluster+"-"+v.Label] = miner
		} else {
			none = append(none, v.PrivateIpAddress)
			maps[v.Cluster+"-"+v.Label] = none
		}

	}
	file, err := os.OpenFile("logs/hosts.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed ,err:", err)
		return
	}
	defer file.Close()
	for k, v := range maps {
		file.WriteString("[" + k + "]\n")
		for _, ip := range v {
			file.WriteString(ip + "\n")
		}
	}
	/*
		[chengdu-1-lotus-miner]
		 172.0.0.2
		[chengdu-2-lotus-miner]
		 172.0.0.4
		[chengdu-1-lotus-worker]
		 172.0.0.1
		 172.0.0.3
	*/
}

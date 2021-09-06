package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type alert struct {
	Status string `json:"status"`
	Data   Data
}
type Data struct {
	Data   string `json:"data"`
	Alerts []Alerts
}
type Alerts struct {
	Value  string `json:"value"`
	Labels Labels
}
type Labels struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Cluster   string `json:"cluster"`
	Device    string `json:"device"`
}

func main() {
	client := &http.Client{Timeout: 200 * time.Millisecond}
	resp, _ := client.Get("http://30.10.0.18:39090/api/v1/alerts")
	body, _ := ioutil.ReadAll(resp.Body)
	var a alert
	maps := make(map[string][]string, 0)
	_ = make([]string, 0)
	json.Unmarshal([]byte(body), &a)
	for _, v := range a.Data.Alerts {
		fmt.Println(v.Labels.Alertname, v.Labels.Instance, v.Labels.Cluster, v.Value)
		if maps[v.Labels.Alertname] == nil {

		}

	}
}

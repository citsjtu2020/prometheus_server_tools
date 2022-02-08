package prometheus_tools

import (
	"encoding/json"
	"io/ioutil"
)

type ProServer struct {
	Node map[string]string `json:"node"`
	Addr string `json:"addr"`
	Port int `json:"port"`
	NodeIP map[string]string `json:"node_ip"`
	Metrics []string `json:"metrics"`
}



func (ps *ProServer) Load(jsonfile string){
	data,err := ioutil.ReadFile(jsonfile)
	if err != nil{
		return
	}

	err = json.Unmarshal(data,ps)
	if err != nil{
		return
	}
}

package prometheus_tools

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"sync"
	"time"
)

func Gen_Influx_Addr(ip string,port int) string{
	address := fmt.Sprintf("http://%s:%v",ip,port)
	return address
}

func ConnInflux(address string,username string,password string) client.Client {
	//"http://127.0.0.1:8086"
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     address,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

//query
func QueryDB(cli client.Client, cmd string,database string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//Insert
func WritesContainerPoints(cli client.Client,database string,precision string,measurement string,data []ContainerResult) (bool,error){

	ok := false
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: precision,
	})
	if err != nil {
		log.Fatal(err)
	}
	//{"cpu": "ih-cpu"}
	var err0 error
	err0 = nil
	for _,v := range data{
		tags := make(map[string]string)
		tags["node"] = v.node
		tags["deploy"] = v.deploy
		tags["namespace"] = v.namespace
		tags["pod"] = v.pod
		tags["container"] = v.container
		//tags["metric"]
		fields := map[string]interface{}{
		"value": v.GetValuesi(0).value,
	}
 	//time.Now()
		//b1 := v.GetValuesi(0).timestamp-float64(int64(v.GetValuesi(0).timestamp)) > 0.0

		//if b1
	pt, err := client.NewPoint(
		measurement,
		tags,
		fields,
		time.Unix(int64(v.GetValuesi(0).timestamp),0),
	)
	if err != nil {
		//log.Fatal(err)
		err0 = err
	}
	bp.AddPoint(pt)
	}
	if err := cli.Write(bp); err != nil {
		//log.Fatal(err)
		return ok,err
	}else {
		ok = true
	}
	if err0 != nil{
		return ok,err0
	}else{
		return ok,nil
	}
}

func WritesContainerAccPoints(cli client.Client,database string,precision string,measurement string,
	data map[string]map[string]map[string]map[int64]*ContainerAccResult,
	wg *sync.WaitGroup,step int) (bool,error){
	defer func() {
		wg.Done()
	}()
	ok := false
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: precision,
	})
	if err != nil {
		log.Fatal(err)
	}
	//{"cpu": "ih-cpu"}
	var err0 error
	err0 = nil
	for k1,v1 := range data {
		for k2,v2 := range v1{
			for k3,v3 := range v2{
				for k4,_ := range v3{
					v := data[k1][k2][k3][k4]
					tags := make(map[string]string)
					tags["node"] = v.node
					tags["deploy"] = v.deploy
					tags["namespace"] = v.namespace
					tags["pod"] = v.pod
					tags["container"] = v.container
					tags["step"] = fmt.Sprintf("%v%s",step,precision)

		//tags["metric"]
					fields := make(map[string]interface{})
					for k5,v5 := range v.GetResults(){
						fields[k5] = v5
					}
					pt, err := client.NewPoint(
						measurement,
						tags,
						fields,
						time.Unix(int64(v.timestamp),0),
					)
					if err != nil {
					//log.Fatal(err)
						err0 = err
					}
					bp.AddPoint(pt)
				}
			}
		}
	}

	if err := cli.Write(bp); err != nil {
		//log.Fatal(err)
		return ok,err
	}else {
		ok = true
	}

	if err0 != nil{
		return ok,err0
	}else{
		return ok,nil
	}
}

func WritesContainerNetworkAccPoints(cli client.Client,database string,precision string,measurement string,
	data map[string]map[string]map[string]map[string]map[int64]*ContainerNetworkAccResult,
	wg *sync.WaitGroup,step int) (bool,error){
	defer func() {
		wg.Done()
	}()
	ok := false
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: precision,
	})
	if err != nil {
		log.Fatal(err)
	}
	//{"cpu": "ih-cpu"}
	var err0 error
	err0 = nil
	for k1,v1 := range data {
		for k2,v2 := range v1{
			for k3,v3 := range v2{
				for k31,v31 := range v3{
					for k4,_ := range v31{
						v := data[k1][k2][k3][k31][k4]
						tags := make(map[string]string)
						tags["node"] = v.node
						tags["deploy"] = v.deploy
						tags["namespace"] = v.namespace
						tags["pod"] = v.pod
						tags["inter"] = v.inter
						tags["container"] = v.container
						tags["step"] = fmt.Sprintf("%v%s",step,precision)
		//tags["metric"]
						fields := make(map[string]interface{})
						for k5,v5 := range v.GetResults(){
							fields[k5] = v5
						}
						pt, err := client.NewPoint(
							measurement,
							tags,
							fields,
							time.Unix(int64(v.timestamp),0),
						)
						if err != nil {
						//log.Fatal(err)
							err0 = err
						}
						bp.AddPoint(pt)
					}
				}
			}
		}
	}

	if err := cli.Write(bp); err != nil {
		//log.Fatal(err)
		return ok,err
	}else {
		ok = true
	}

	if err0 != nil{
		return ok,err0
	}else{
		return ok,nil
	}
}


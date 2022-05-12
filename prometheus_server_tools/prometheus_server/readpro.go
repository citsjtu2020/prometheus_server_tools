package main

import (
	"prometheus_server_tools/prometheus_controller"
	//"flag"
	"prometheus_server_tools/prometheus_tools"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gonum.org/v1/gonum/mat"
	"sync"
	"time"

	//"path/filepath"
)

var months map[string]int

func NewProServer() *prometheus_tools.ProServer {
	return &prometheus_tools.ProServer{}
}

func ProServerParse(jsonfile string) *prometheus_tools.ProServer {
	proserver := NewProServer()
	proserver.Load(jsonfile)
	//fmt.Println(*proserver)
	return proserver
}

func matPrint(X mat.Matrix){
	fa := mat.Formatted(X,mat.Prefix(""),mat.Squeeze())
	fmt.Printf("%v\n",fa)
}

func init() {
	//January Month = 1 + iota
	//	February
	//	March
	//	April
	//	May
	//	June
	//	July
	//	August
	//	September
	//	October
	//	November
	//	December
	months = make(map[string]int)
	months["January"] = 1
	months["February"] = 2
	months["March"] = 3
	months["April"] = 4
	months["May"] = 5
	months["June"] = 6
	months["July"] = 7
	months["August"] = 8
	months["September"] = 9
	months["October"] = 10
	months["November"] = 11
	months["December"] = 12
}





func main() {
	var wg sync.WaitGroup
	var wg_write sync.WaitGroup
	var mc *prometheus_controller.ManageController
	defer func() {
		if mc != nil{
			mc.Stop()
		}
	}()
	var timer *time.Timer
	var jsonfile *string
	var influxip *string
	var influxport *int
	var influxuser *string
	var influxpwd *string
	var influxbase *string
	var duration *int
	var unit *string
	var step *int
	//var serverport *int
	jsonfile = flag.String("file","pro_metric_config.json","(optional) absolute path to the kubeconfig file")
	//serverport = flag.Int("proport",9090,"prome")
	//
	influxip = flag.String("influxip","192.168.1.11","(optional) absolute path to the kubeconfig file")
	influxport = flag.Int("influxport",8086,"(optional) absolute path to the kubeconfig file")
	influxuser = flag.String("influxuser","admin","(optional) absolute path to the kubeconfig file")
	influxpwd = flag.String("influxpwd","admin","(optional) absolute path to the kubeconfig file")
	influxbase = flag.String("influxbase","prometheus","(optional) absolute path to the kubeconfig file")
	duration = flag.Int("duration",30,"(optional) duration of metrics")
	unit = flag.String("unit","s","(optional) unit of sampling")
	step = flag.Int("step",3,"(optional) step of sampling")

	prom_param := ProServerParse(*jsonfile)
	fmt.Println(prom_param.Addr)
	fmt.Println(prom_param.Port)
	fmt.Println(prom_param)
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UTC())
	fmt.Println(time.Now().Date())
	fmt.Println(time.Now().Hour())
	fmt.Println(time.Now().Minute())
	fmt.Println(time.Now().Second())
	//	return (t.unixSec())*1e9 + int64(t.nsec())
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano())
	//nowStr = now.Format("2006-January-02 03:04:05.999 pm")
	//	fmt.Println(nowStr)
	nowStr := time.Now().UTC().Format("2022-02-02T15:32:07.999Z")
	fmt.Println(nowStr)
	server_address := prometheus_tools.Generate_address(prom_param.Addr,prom_param.Port)
	fmt.Println(server_address)
	url := prometheus_tools.Gen_range_url(server_address,prom_param.Metrics[0],time.Now().Unix()-120,time.Now().Unix()-60,"s",5)
	//container_network_receive_bytes_total
	url2 := prometheus_tools.Gen_range_url(server_address,"container_network_receive_bytes_total",time.Now().Unix()-120,time.Now().Unix()-60,"s",5)
	raws := prometheus_tools.Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	//check(err)
	fmt.Println(err)
	datas,err := js.Get("data").Get("result").Array()
	//check(err)
	t := len(datas)
	fmt.Println(t)
	results := prometheus_tools.Container_Raw_Metric(url)
	fmt.Println(len(results))
	fmt.Println(results[0])

	results2 := prometheus_tools.Container_Network_Metric(url2)
	//results2 := prometheus_tools.Container_CPU_Usage(server_address,time.Now().Unix()-120,time.Now().Unix()-96,"s",3)
	fmt.Println(results2)
	inluxaddress := prometheus_tools.Gen_Influx_Addr(*influxip,*influxport)
	conn := prometheus_tools.ConnInflux(inluxaddress,*influxuser,*influxpwd)
	fmt.Println(conn)
	fmt.Println(influxbase)
	conn.Close()
	timer = prometheus_tools.NewGlobalTimer("s",10)
	kkk := make(map[string]int)
	fmt.Println()
	fmt.Println(len(kkk))
	kkk["c"] = 1
	fmt.Println(len(kkk))

	tt := make(chan int)
	//_,okchan := <- tt
	//fmt.Println(okchan)
	close(tt)
	_,okchan := <- tt
	fmt.Println(okchan)
	_,okchan = <- tt
	fmt.Println(okchan)

	mc = prometheus_controller.InitMC()
	mc.RegisterContainerController("cpu_util",1,prometheus_tools.Container_CPU_Usage,false)
	mc.RegisterContainerController("cpu_load",1,prometheus_tools.Container_Load,false)
	mc.RegisterContainerController("cpu_thro",1,prometheus_tools.Container_CPU_Throttled,false)
	mc.RegisterContainerController("memory_max",1,prometheus_tools.Container_Memory_Max,false)
	mc.RegisterContainerController("memory_mean",1,prometheus_tools.Container_Memory_Mean,false)
	mc.RegisterContainerNetworkController("network_err",1,prometheus_tools.Container_Network_Error,false)
	mc.RegisterContainerNetworkController("network_bytes",1,prometheus_tools.Container_Network_Bytes,false)
	mc.RegisterContainerNetworkController("network_packet",1,prometheus_tools.Container_Network_Packet,false)
	for key,get_func := range mc.Container_func{
		go prometheus_controller.Container_Result(server_address,int64(*duration),*unit,*step,mc.Container_res[key],mc.Global_time[key],&wg,get_func)
	}
	for key,get_func := range mc.Container_net_func{
		go prometheus_controller.Container_Network_Result(server_address,int64(*duration),*unit,*step,mc.Container_net_res[key],mc.Global_time[key],&wg,get_func)
	}

	//global_container_res := make(map[string])

	fmt.Println(mc.GetNumber())
	for{
		select {
		case <-timer.C:
			fmt.Println(time.Now().Unix())
			wg.Add(mc.GetNumber())
			for k,_ := range mc.Global_time{
				mc.Global_time[k] <- time.Now().Unix()
			}
			wg.Wait()
			container_res := make(map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult)
			for k,_ := range mc.Container_res{
				res,ok := <- mc.Container_res[k]
				if ok{
					//fmt.Println(len(res))
					for k1,v1 := range res{
						_,ok1 := container_res[k1]
						if !ok1{
							container_res[k1] = make(map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult)
						}
						for k2,v2 := range v1{
							_,ok2 := container_res[k1][k2]
							if !ok2{
								container_res[k1][k2] = make(map[string]map[int64]*prometheus_tools.ContainerAccResult)
							}
							for k3,v3 := range v2{
								_,ok3 := container_res[k1][k2][k3]
								if !ok3{
									container_res[k1][k2][k3] = make(map[int64]*prometheus_tools.ContainerAccResult)
								}
								for k4,v4 := range v3{
									_,ok4 := container_res[k1][k2][k3][k4]
									if !ok4{
										container_res[k1][k2][k3][k4] = v4
									}else{
										for k5,v5 := range v4.GetResults(){
											container_res[k1][k2][k3][k4].SetResultItem(k5,v5)
										}
									}
								}
							}
						}
					}
				}

			}
			container_net_res := make(map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
			for k,_ := range mc.Container_net_res{
				res,ok := <- mc.Container_net_res[k]
				if ok{
					for k1,v1 := range res{
						_,ok1 := container_net_res[k1]
						if !ok1{
							container_net_res[k1] = make(map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
						}
						for k2,v2 := range v1{
							_,ok2 := container_net_res[k1][k2]
							if !ok2{
								container_net_res[k1][k2] = make(map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
							}
							for k3,v3 := range v2{
								_,ok3 := container_net_res[k1][k2][k3]
								if !ok3{
									container_net_res[k1][k2][k3] = make(map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
								}
								for k4,v4 := range v3{
									_,ok4 := container_net_res[k1][k2][k3][k4]
									if !ok4{
										container_net_res[k1][k2][k3][k4] = make(map[int64]*prometheus_tools.ContainerNetworkAccResult)
									}
									for k41,v41 := range v4{
										_,ok41 := container_net_res[k1][k2][k3][k4][k41]
										if !ok41{
											container_net_res[k1][k2][k3][k4][k41] = v41
										}else{
											for k5,v5 := range v41.GetResults(){
											container_net_res[k1][k2][k3][k4][k41].SetResultItem(k5,v5)
											}
										}
									}
								}
							}
						}
					}
				}
			}
			wg_write.Add(2)
			conn1 := prometheus_tools.ConnInflux(inluxaddress,*influxuser,*influxpwd)
			fmt.Println(conn1)
			//fmt.Println(influxbase)
			conn2 := prometheus_tools.ConnInflux(inluxaddress,*influxuser,*influxpwd)
			fmt.Println(conn2)
	//conn.Close()
	//"s","cpuusage",results2
			go prometheus_tools.WritesContainerAccPoints(conn1,*influxbase,*unit,"container",container_res,&wg_write,*step)
			go prometheus_tools.WritesContainerNetworkAccPoints(conn2,*influxbase,*unit,"container_net",container_net_res, &wg_write,*step)
			wg_write.Wait()
			_ = conn1.Close()
			_ = conn2.Close()
			timeout := prometheus_tools.Gen_Timeout(*unit,*duration)
			timer.Reset(timeout)
		default:
			continue
		}
	}
	//
	////insert
	//ok,err := prometheus_tools.WritesContainerPoints(conn,*influxbase,"s","cpuusage",results2)
	//fmt.Println(ok)
	//fmt.Println(err)
	//url2 := prometheus_tools.Gen_range_url(server_address,"node_cpu_seconds_total",time.Now().Unix()-120,time.Now().Unix()-60,"s",5)
	//node_seconds_metric := "node_time_seconds"
	//url2 := prometheus_tools.Gen_range_url(server_address,node_seconds_metric,time.Now().Unix()-120,time.Now().Unix()-60,"s",5)

	//results2 := prometheus_tools.Node_Raw_Metric(url2)
	//fmt.Println(len(results2))
	//fmt.Println(results2[0])
	//dd := make(map[string]*prometheus_tools.NodeResult)
	//for k,_ := range results2{
	//	//res:= results2[k]
	//	dd[results2[k].GetNode()] = &results2[k]
	//	fmt.Println(&results2[k])
	//}
	//for k,_ := range dd{
	//	fmt.Println(dd[k])
	//}
	//kk,ok := dd["a"]
	//if kk != nil{
	//	fmt.Println("kkk")
	//}
	//fmt.Println(kk)
	//fmt.Println(ok)
	//results = make([]prometheus_tools.ContainerResult,0,t)
	//tmpresult := prometheus_tools.ContainerResult{}
	//
	//for _,data := range datas{
	//	//fmt.Println(data["metric"])
	//	//data3,_ := data.(map[string]interface{})
	//	//data3['']
	//	bytedata,_ := json.Marshal(data)
	//	data2,_ := simplejson.NewJson(bytedata)
	//	//fmt.Println(data2)
	//	container,_ := data2.Get("metric").Get("container").String()
	//	//fmt.Println(container)
	//	//fmt.Println(data2.Get("metric").Get("container"))
	//
	//	tmpresult.SetContainer(container)
	//	//fmt.Println(tmpresult.GetContainer())
	//	node,_ := data2.Get("metric").Get("instance").String()
	//	tmpresult.SetNode(node)
	//	tmppod,_ := data2.Get("metric").Get("pod").String()
	//	//tmpresult.pod = tmppod
	//	tmpresult.SetPod(tmppod)
	//	namespace,_ := data2.Get("metric").Get("namespace").String()
	//	tmpresult.SetNamespace(namespace)
	//	tmpdeploys := strings.Split(tmppod,"-")
	//	//tmpresult.deploy = tmpdeploys[0]
	//	tmpresult.SetDeploy(tmpdeploys[0])
	//	//fmt.Println(tmpresult.deploy)
	//	//fmt.Println(tmpresult.GetDeploy())
	//	//fmt.Println(tmpresult)
	//	results = append(results, tmpresult)
	//	}
	//	fmt.Println(len(results))
	//	fmt.Println(results[0])
	//2022-02-02 08:22:20.986032377 +0000 UTC
	//2021-07-24T15:57:07.000Z
	//2021-07-25T14:47:15.000Z

	//u := mat.NewVecDense(3,[]float64{1,2,3})
	//v := mat.NewVecDense(3,[]float64{4,5,6})
	//fmt.Println(u)
	//fmt.Println(v)
	//fmt.Println(u.AtVec(1))
	//fmt.Println(u.At(1,0))
	//u.SetVec(1,33.2)
	//fmt.Println(u)

	//w := mat.NewVecDense(3,nil)
	//w.AddVec(u,v)
	//
	//fmt.Println(w)
	//
	//w.AddScaledVec(u,2,v)
	//fmt.Println(w)
	//
	//w.SubVec(v,u)
	//fmt.Println("v - u:")
	//matPrint(w)
	//
	//w.ScaleVec(23,u)
	//fmt.Println("u*23:")
	//matPrint(w)

	//fmt.Println(mat.Dot(v,u))
	//
	//l := v.Len()
	//fmt.Println("length of V:",l)
	//
	//fmt.Println(mat.Norm(v,2)*mat.Norm(v,2))
	//matPrint(v)
	//matPrint(v.T())
	//
	//v2 := []float64{1,2,3,4,5,6,7,8,9,10,11,12}
	//A := mat.NewDense(3,4,v2)
	//matPrint(A)
}

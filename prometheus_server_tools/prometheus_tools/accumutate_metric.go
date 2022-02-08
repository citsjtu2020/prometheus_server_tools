package prometheus_tools

import (
	"strings"
)

//import "fmt"

func GenerateContainerMap(result[]ContainerResult) map[string]map[string]map[string]*ContainerResult{
	mapresult := make(map[string]map[string]map[string]*ContainerResult)
	for k, _:= range result{
		tmpnamespace := result[k].GetNamespace()
		_,ok := mapresult[tmpnamespace]
		if !ok{
			mapresult[tmpnamespace] = make(map[string]map[string]*ContainerResult)
		}
		tmppod := result[k].GetPod()
		_,ok = mapresult[tmpnamespace][tmppod]
		if !ok{
			mapresult[tmpnamespace][tmppod] = make(map[string]*ContainerResult)
		}
		tmpcontainer := result[k].GetContainer()
		//_,ok = mapresult[tmpnamespace][tmppod][tmpcontainer]
		//if !ok{
		mapresult[tmpnamespace][tmppod][tmpcontainer] = &result[k]
		//}
	}
	return mapresult
}

func GenerateContainerNetworkMap(result[]ContainerNetworkResult) map[string]map[string]map[string]map[string]*ContainerNetworkResult{
	mapresult := make(map[string]map[string]map[string]map[string]*ContainerNetworkResult)
	for k, _:= range result{
		tmpnamespace := result[k].GetNamespace()
		_,ok := mapresult[tmpnamespace]
		if !ok{
			mapresult[tmpnamespace] = make(map[string]map[string]map[string]*ContainerNetworkResult)
		}
		tmppod := result[k].GetPod()
		_,ok = mapresult[tmpnamespace][tmppod]
		if !ok{
			mapresult[tmpnamespace][tmppod] = make(map[string]map[string]*ContainerNetworkResult)
		}
		tmpcontainer := result[k].GetContainer()
		_,ok = mapresult[tmpnamespace][tmppod][tmpcontainer]
		if !ok{
			mapresult[tmpnamespace][tmppod][tmpcontainer] = make(map[string]*ContainerNetworkResult)
		}
		//_,ok = mapresult[tmpnamespace][tmppod][tmpcontainer]
		//if !ok{
		tmpinter := result[k].GetInfer()
		mapresult[tmpnamespace][tmppod][tmpcontainer][tmpinter] = &result[k]
		//}
	}
	return mapresult
}

func Container_CPU_Usage(address string,start int64,end int64,unit string,step int)([]ContainerAccResult){
	usage_metric := "container_cpu_usage_seconds_total"
	node_seconds_metric := "node_time_seconds"
	quota_metric := "container_spec_cpu_quota"
	cpu_perioid_metric := "container_spec_cpu_period"
	//
	usage_url := Gen_range_url(address,usage_metric,start,end,unit,step)
	node_second_url := Gen_range_url(address,node_seconds_metric,start,end,unit,step)
	quota_url := Gen_range_url(address,quota_metric,start,end,unit,step)
	cpu_perioid_url := Gen_range_url(address,cpu_perioid_metric,start,end,unit,step)

	usage_result := Container_Raw_Metric(usage_url)
	node_second_result := Node_Raw_Metric(node_second_url)
	quota_result := Container_Raw_Metric(quota_url)
	cpu_period_result := Container_Raw_Metric(cpu_perioid_url)
	//
	//fmt.Println(usage_result)
	//fmt.Println(quota_result)
	//fmt.Println(cpu_period_result)
	node_seconds := make(map[string]*NodeResult)
	var usage map[string]map[string]map[string]*ContainerResult
	var quota map[string]map[string]map[string]*ContainerResult
	var cpu_period map[string]map[string]map[string]*ContainerResult

	usage = GenerateContainerMap(usage_result)
	quota = GenerateContainerMap(quota_result)
	cpu_period = GenerateContainerMap(cpu_period_result)
	//for k,_ := range results2{
	//		//res:= results2[k]
	//		dd[results2[k].GetNode()] = &results2[k]
	//		fmt.Println(&results2[k])
	//	}
	for k,_ := range node_second_result{
		node_seconds[node_second_result[k].GetNode()] = &node_second_result[k]
	}
	results := make([]ContainerAccResult,0,len(usage_result))
	//cpu_quota := make([]ContainerResult,0,len(usage_result))
	for k1,v1 := range usage{
		_,ok11 := quota[k1]
		_,ok12 := cpu_period[k1]

		//wrongvalue1 := false
		if (!ok11 || !ok12){
			continue
		}else{
			for k2,v2 := range v1{
				_,ok21 := quota[k1][k2]
				_,ok22 := cpu_period[k1][k2]
				//fmt.Println(v2)
				if (!ok21 || !ok22){
					continue
				}else{
					for k3,v3 := range v2{
						tmpquota,ok31 := quota[k1][k2][k3]
						tmp_period,ok32 := cpu_period[k1][k2][k3]
						node_name := strings.Split(v3.node,"-")
						tmp_time,ok33 := node_seconds[node_name[len(node_name)-1]]
						if (!ok31 || !ok32 || !ok33){
							continue
						}
						tmputil := ContainerAccResult{}
						tmputil.node = v3.GetNode()
						tmputil.namespace = k1
						tmputil.pod = k2
						tmputil.container = k3
						tmputil.deploy = v3.GetDeploy()

						tmppoint := Rawpoint{}
						quotapoint := Rawpoint{}

						//tmppoint.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
						totaltime := (tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
						int_total_time := int(totaltime)
						if int_total_time == 0 && totaltime < 0.001 && totaltime > -0.001{
							tmppoint.value = 0
						}else{
							tmppoint.value = (v3.GetValuesi(len(v3.values)-1).value - v3.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
						}
						//quotapoint.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
						quotapoint.value = (compute_mean(tmpquota.values)/ compute_mean(tmp_period.values))

						tmputil.timestamp = int64(v3.GetValuesi(len(v3.values)-1).timestamp)
						tmputil.SetResultItem("cpu_util",tmppoint.value)
						tmputil.SetResultItem("cpu_quota",quotapoint.value)
						results = append(results,tmputil)
					}
				}
			}
		}

	}
	return results
}

//usage_data = pd.read_csv(os.path.join(
//                os.path.join(os.path.join(datapath, 'container_cpu_cfs_throttled_periods_total'),
//                             container_place['container_cpu_cfs_throttled_periods_total'][k]), k + ".csv"))
//            spec_quota = pd.read_csv(os.path.join(os.path.join(os.path.join(datapath, 'container_cpu_cfs_periods_total'),
//                                                               container_place['container_cpu_cfs_periods_total'][k]),
//
//
//                                                k + ".csv"))

func Container_CPU_Throttled(address string,start int64,end int64,unit string,step int)[]ContainerAccResult{
	usage_metric := "container_cpu_cfs_throttled_periods_total"
	quota_metric := "container_cpu_cfs_periods_total"
	//
	usage_url := Gen_range_url(address,usage_metric,start,end,unit,step)
	quota_url := Gen_range_url(address,quota_metric,start,end,unit,step)

	usage_result := Container_Raw_Metric(usage_url)
	quota_result := Container_Raw_Metric(quota_url)
	//
	//fmt.Println(usage_result)
	//fmt.Println(quota_result)
	//fmt.Println(cpu_period_result)
	var usage map[string]map[string]map[string]*ContainerResult
	var quota map[string]map[string]map[string]*ContainerResult

	usage = GenerateContainerMap(usage_result)
	quota = GenerateContainerMap(quota_result)

	cpu_thro := make([]ContainerAccResult,0,len(usage_result))

	for k1,v1 := range usage{
		_,ok11 := quota[k1]

		//wrongvalue1 := false
		if (!ok11){
			continue
		}else{
			for k2,v2 := range v1{
				_,ok21 := quota[k1][k2]
				if (!ok21){
					continue
				}else{
					for k3,v3 := range v2{
						tmpquota,ok31 := quota[k1][k2][k3]
						if (!ok31){
							continue
						}
						tmpthro := ContainerAccResult{}
						tmpthro.node = v3.GetNode()
						tmpthro.namespace = k1
						tmpthro.pod = k2
						tmpthro.container = k3
						tmpthro.deploy = v3.GetDeploy()
						tmppoint := Rawpoint{}
						tmpthro.timestamp = int64(v3.GetValuesi(len(v3.values)-1).timestamp)
						//((usage_data['value'][j] - usage_data['value'][j-7])/(time_data['value'][j] - time_data['value'][j-7]))
						// /(np.mean(spec_quota['value'][j-7:j+1])/np.mean(spec_period['value'][j-7:j+1]))
						tmppoint.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
						total_period := int(tmpquota.GetValuesi(len(tmpquota.values)-1).value - tmpquota.GetValuesi(0).value)
						float_period := (tmpquota.GetValuesi(len(tmpquota.values)-1).value - tmpquota.GetValuesi(0).value)
						if total_period == 0 && float_period < 0.001 && float_period > -0.001{
							tmppoint.value = 0
						}else{
							tmppoint.value = (v3.GetValuesi(len(v3.values)-1).value - v3.GetValuesi(0).value)/(tmpquota.GetValuesi(len(tmpquota.values)-1).value - tmpquota.GetValuesi(0).value)
						}
						//fmt.Println(len(v3.values))

						tmpthro.SetResultItem("cpu_thro",tmppoint.value)
						cpu_thro = append(cpu_thro,tmpthro)
					}
				}
			}
		}
	}
	return cpu_thro
}

func Container_Memory_Max(address string,start int64,end int64,unit string,step int)[]ContainerAccResult{
	usage_metric := "container_memory_working_set_bytes"
	//quota_metric := "container_cpu_cfs_periods_total"
	//
	usage_url := Gen_range_url(address,usage_metric,start,end,unit,step)

	usage_result := Container_Raw_Metric(usage_url)

	var usage map[string]map[string]map[string]*ContainerResult

	usage = GenerateContainerMap(usage_result)

	memory_max := make([]ContainerAccResult,0,len(usage_result))

	for k1,v1 := range usage{
		for k2,v2 := range v1{
			for k3,v3 := range v2{
				tmpmax := ContainerAccResult{}
				tmpmax.node = v3.GetNode()
				tmpmax.namespace = k1
				tmpmax.pod = k2
				tmpmax.container = k3
				tmpmax.deploy = v3.GetDeploy()
				tmpmax.timestamp = int64(v3.GetValuesi(len(v3.values)-1).timestamp)
				tmppoint := Rawpoint{}
				tmppoint.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
				tmppoint.value = (compute_max(v3.GetValues())/ 1024) / 1024
				tmpmax.SetResultItem("memory_max",tmppoint.value)
				memory_max = append(memory_max,tmpmax)
			}
		}
	}
	return memory_max
}

func Container_Memory_Mean(address string,start int64,end int64,unit string,step int)([]ContainerAccResult){
	usage_metric := "container_memory_working_set_bytes"
	quota_metric := "container_spec_memory_limit_bytes"

	//container_spec_memory_limit_bytes
	usage_url := Gen_range_url(address,usage_metric,start,end,unit,step)
	quota_url := Gen_range_url(address,quota_metric,start,end,unit,step)


	usage_result := Container_Raw_Metric(usage_url)
	quota_result := Container_Raw_Metric(quota_url)


	var usage map[string]map[string]map[string]*ContainerResult
	var quota map[string]map[string]map[string]*ContainerResult


	usage = GenerateContainerMap(usage_result)
	quota = GenerateContainerMap(quota_result)

	memory_mean := make([]ContainerAccResult,0,len(usage_result))
	for k1,v1 := range usage{
		_,ok11 := quota[k1]
		for k2,v2 := range v1{
			_,ok21 := quota[k1][k2]
			for k3,v3 := range v2{
				quota3,ok31 := quota[k1][k2][k3]
				tmpmean := ContainerAccResult{}
				tmpmean.node = v3.GetNode()
				tmpmean.namespace = k1
				tmpmean.pod = k2
				tmpmean.container = k3
				tmpmean.deploy = v3.GetDeploy()
				tmpmean.timestamp = int64(v3.GetValuesi(len(v3.values)-1).timestamp)
				tmppoint := Rawpoint{}
				tmppoint.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
				tmppoint.value = (compute_mean(v3.GetValues()) / 1024) / 1024
				tmpmean.SetResultItem("memory_mean",tmppoint.value)
				//memory_mean = append(memory_mean,tmpmean)
				if ok11 && ok21 && ok31{
					tmppoint2 := Rawpoint{}
					tmppoint2.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
					tmppoint2.value = (compute_mean(quota3.GetValues()) / 1024) / 1024
					tmpmean.SetResultItem("memory_quota",tmppoint2.value)

				}
				memory_mean = append(memory_mean,tmpmean)
			}
		}
	}
	return memory_mean
}

func Container_Load(address string,start int64,end int64,unit string,step int)[]ContainerAccResult{
	usage_metric := "container_cpu_load_average_10s"
	//quota_metric := "container_cpu_cfs_periods_total"
	//
	usage_url := Gen_range_url(address,usage_metric,start,end,unit,step)

	usage_result := Container_Raw_Metric(usage_url)

	var usage map[string]map[string]map[string]*ContainerResult

	usage = GenerateContainerMap(usage_result)
	load_mean := make([]ContainerAccResult,0,len(usage_result))

	for k1,v1 := range usage{
		for k2,v2 := range v1{
			for k3,v3 := range v2{
				tmpmean := ContainerAccResult{}
				tmpmean.node = v3.GetNode()
				tmpmean.namespace = k1
				tmpmean.pod = k2
				tmpmean.container = k3
				tmpmean.deploy = v3.GetDeploy()
				tmppoint := Rawpoint{}
				tmppoint.timestamp = v3.GetValuesi(len(v3.values)-1).timestamp
				//tmppoint.value = (compute_mean(v3.GetValues()) / 1024) / 1024
				tmpmean.timestamp = int64(tmppoint.timestamp)
				if step > 10{
					tmppoint.value = (compute_mean(v3.GetValues()) / 1000)
				}else{
					buchang := int(10/step)
					tmploads :=[]Rawpoint{v3.GetValuesi(len(v3.GetValues())-1)}
					for i := len(v3.GetValues())-1;i>=0;i-=buchang{
						if i == len(v3.GetValues())-1{
							continue
						}else{
							tmploads = append(tmploads,v3.GetValuesi(i))
						}
					}
					tmppoint.value = compute_mean(tmploads)
				}
				tmpmean.SetResultItem("cpu_load",tmppoint.value)
				load_mean = append(load_mean,tmpmean)
			}
		}
	}
	return load_mean
}



func compute_mean(input []Rawpoint) float64{
	if input == nil || len(input) == 0{
		return 0
	}else{
		tmpmean := 0.0
		for _,v := range input{
			tmpmean += v.value
		}
		return (tmpmean / float64(len(input)))
	}
}

func compute_max(input []Rawpoint) float64{
	if input == nil || len(input) == 0{
		return 0
	}else{
		tmpmax := input[0].value
		for _,v := range input{
			if v.value > tmpmax{
				tmpmax = v.value
			}
		}
		return tmpmax
	}
}

//'receive_error_rate_60','transmit_error_rate_60',
//'receive_drop_rate_60','transmit_drop_rate_60',
//'receive_bytes_rate_60','transmit_bytes_rate_60',
//"receive_packet_60","transmit_packet_60",
//"receive_per_packet","transmit_per_packet"

//container_network_transmit_errors_total
func Container_Network_Error(address string,start int64,end int64,unit string,step int)([]ContainerNetworkAccResult){
	transmit_metric := "container_network_transmit_errors_total"
	receive_metric := "container_network_receive_errors_total"
	node_seconds_metric := "node_time_seconds"
	//quota_metric := "container_spec_cpu_quota"
	//cpu_perioid_metric := "container_spec_cpu_period"
	//
	transmit_url := Gen_range_url(address,transmit_metric,start,end,unit,step)
	node_second_url := Gen_range_url(address,node_seconds_metric,start,end,unit,step)
	receive_url := Gen_range_url(address,receive_metric,start,end,unit,step)

	transmit_result := Container_Network_Metric(transmit_url)
	receive_result := Container_Network_Metric(receive_url)
	node_second_result := Node_Raw_Metric(node_second_url)

	//
	//fmt.Println(usage_result)
	//fmt.Println(quota_result)
	//fmt.Println(cpu_period_result)
	node_seconds := make(map[string]*NodeResult)
	var transmit map[string]map[string]map[string]map[string]*ContainerNetworkResult
	var receive map[string]map[string]map[string]map[string]*ContainerNetworkResult

	transmit = GenerateContainerNetworkMap(transmit_result)
	receive = GenerateContainerNetworkMap(receive_result)

	for k,_ := range node_second_result{
		node_seconds[node_second_result[k].GetNode()] = &node_second_result[k]
	}
	results_error := make([]ContainerNetworkAccResult,0,len(transmit_result))
	//recv_error := make([]ContainerNetworkAccResult,0,len(receive_result))
	for k1,v1 := range transmit{
		_,ok11 := receive[k1]
		if (!ok11){
			continue
		}else{
			for k2,v2 := range v1{
				_,ok21 := receive[k1][k2]
				//fmt.Println(v2)
				if (!ok21){
					continue
				}else{
					for k3,v3 := range v2{
						_,ok31 := receive[k1][k2][k3]
						if (!ok31){
							continue
						}else{
							for k4,v4 := range v3{
								tmpreceive,ok41 := receive[k1][k2][k3][k4]

								node_name := strings.Split(v4.node,"-")
								tmp_time,ok42 := node_seconds[node_name[len(node_name)-1]]
								if (!ok41 || !ok42){
									continue
								}
								tmp_err := ContainerNetworkAccResult{}
								tmp_err.node = v4.GetNode()
								tmp_err.namespace = k1
								tmp_err.pod = k2
								tmp_err.container = k3
								tmp_err.inter = k4
								tmp_err.deploy = v4.GetDeploy()
								tmp_err.timestamp = int64(v4.GetValuesi(len(v4.values)-1).timestamp)

								transpoint := Rawpoint{}
								recvpoint := Rawpoint{}
								transpoint.timestamp = v4.GetValuesi(len(v4.values)-1).timestamp
								totaltime := (tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
								int_total_time := int(totaltime)
								if int_total_time == 0 && totaltime < 0.001 && totaltime > -0.001{
									transpoint.value = 0
									recvpoint.value = 0
								}else{
									transpoint.value = (v4.GetValuesi(len(v4.values)-1).value - v4.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
									recvpoint.value = (tmpreceive.GetValuesi(len(tmpreceive.values)-1).value - tmpreceive.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
								}
								//tmp_trans_err.SetValues([]Rawpoint{transpoint})
								//tmp_recv_err.SetValues([]Rawpoint{recvpoint})
								tmp_err.SetResultItem("trans_error",transpoint.value)
								tmp_err.SetResultItem("recv_error",recvpoint.value)

								results_error = append(results_error,tmp_err)
							}
						}
					}
				}
			}
		}
	}
	return results_error
}

func Container_Network_Packet(address string,start int64,end int64,unit string,step int)([]ContainerNetworkAccResult){
	//container_network_transmit_bytes_total
	transmit_metric := "container_network_transmit_packets_total"
	receive_metric := "container_network_receive_packets_total"

	transmit_drop_metric := "container_network_transmit_packets_dropped_total"
	receive_drop_metric := "container_network_transmit_packets_dropped_total"

	node_seconds_metric := "node_time_seconds"
	//quota_metric := "container_spec_cpu_quota"
	//cpu_perioid_metric := "container_spec_cpu_period"
	//
	transmit_drop_url := Gen_range_url(address,transmit_drop_metric,start,end,unit,step)
	receive_drop_url := Gen_range_url(address,receive_drop_metric,start,end,unit,step)

	transmit_url := Gen_range_url(address,transmit_metric,start,end,unit,step)
	node_second_url := Gen_range_url(address,node_seconds_metric,start,end,unit,step)
	receive_url := Gen_range_url(address,receive_metric,start,end,unit,step)

	transmit_result := Container_Network_Metric(transmit_url)
	receive_result := Container_Network_Metric(receive_url)

	transmit_drop_result := Container_Network_Metric(transmit_drop_url)
	receive_drop_result := Container_Network_Metric(receive_drop_url)
	node_second_result := Node_Raw_Metric(node_second_url)

	//
	//fmt.Println(usage_result)
	//fmt.Println(quota_result)
	//fmt.Println(cpu_period_result)
	node_seconds := make(map[string]*NodeResult)
	var transmit map[string]map[string]map[string]map[string]*ContainerNetworkResult
	var receive map[string]map[string]map[string]map[string]*ContainerNetworkResult

	var transmit_drop map[string]map[string]map[string]map[string]*ContainerNetworkResult
	var receive_drop map[string]map[string]map[string]map[string]*ContainerNetworkResult

	transmit = GenerateContainerNetworkMap(transmit_result)
	receive = GenerateContainerNetworkMap(receive_result)

	transmit_drop = GenerateContainerNetworkMap(transmit_drop_result)
	receive_drop = GenerateContainerNetworkMap(receive_drop_result)

	for k,_ := range node_second_result{
		node_seconds[node_second_result[k].GetNode()] = &node_second_result[k]
	}
	results_packet := make([]ContainerNetworkAccResult,0,len(transmit_result))

	for k1,v1 := range transmit{
		_,ok11 := receive[k1]
		_,ok12 := receive_drop[k1]
		_,ok13 := transmit_drop[k1]
		if (!ok11 || !ok12 || !ok13){
			continue
		}else{
			for k2,v2 := range v1{
				_,ok21 := receive[k1][k2]
				_,ok22 := receive_drop[k1][k2]
				_,ok23 := transmit_drop[k1][k2]
				//fmt.Println(v2)
				if (!ok21 || !ok22 || !ok23){
					continue
				}else{
					for k3,v3 := range v2{
						_,ok31 := receive[k1][k2][k3]
						_,ok32 := receive_drop[k1][k2][k3]
						_,ok33 := transmit_drop[k1][k2][k3]
						if (!ok31 || !ok32 || !ok33){
							continue
						}else{
							for k4,v4 := range v3{
								tmpreceive,ok41 := receive[k1][k2][k3][k4]
								tmpreceive_drop,ok42 := receive_drop[k1][k2][k3][k4]
								tmptransmit_drop,ok43 := transmit_drop[k1][k2][k3][k4]
								node_name := strings.Split(v4.node,"-")
								tmp_time,ok44 := node_seconds[node_name[len(node_name)-1]]
								if (!ok41 || !ok42||!ok43 || !ok44){
									continue
								}
								tmp_packet := ContainerNetworkAccResult{}
								tmp_packet.node = v4.GetNode()
								tmp_packet.namespace = k1
								tmp_packet.pod = k2
								tmp_packet.container = k3
								tmp_packet.inter = k4
								tmp_packet.deploy = v4.GetDeploy()

								transpoint := Rawpoint{}
								recvpoint := Rawpoint{}

								transmit_drop_point := Rawpoint{}
								receive_drop_point := Rawpoint{}
								transpoint.timestamp = v4.GetValuesi(len(v4.values)-1).timestamp
								totaltime := (tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
								int_total_time := int(totaltime)
								if int_total_time == 0 && totaltime < 0.001 && totaltime > -0.001{
									transpoint.value = 0
									recvpoint.value = 0
									transmit_drop_point.value = 0
									receive_drop_point.value = 0
								}else{
									transpoint.value = (v4.GetValuesi(len(v4.values)-1).value - v4.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
									recvpoint.value = (tmpreceive.GetValuesi(len(tmpreceive.values)-1).value - tmpreceive.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
									if int(transpoint.value) == 0{
										transmit_drop_point.value = 0
									}else{
										transmit_drop_point.value = (tmptransmit_drop.GetValuesi(len(tmptransmit_drop.values)-1).value - tmptransmit_drop.GetValuesi(0).value)/(v4.GetValuesi(len(v4.values)-1).value - v4.GetValuesi(0).value)
									}
									if int(recvpoint.value)==0{
										receive_drop_point.value = 0
									}else{
										receive_drop_point.value = (tmpreceive_drop.GetValuesi(len(tmpreceive_drop.values)-1).value - tmpreceive_drop.GetValuesi(0).value)/(tmpreceive.GetValuesi(len(v4.values)-1).value - tmpreceive.GetValuesi(0).value)
									}
								}
								tmp_packet.SetResultItem("recv_packet",recvpoint.value)
								tmp_packet.SetResultItem("trans_packet",transpoint.value)
								tmp_packet.SetResultItem("recv_drop_packet",receive_drop_point.value)
								tmp_packet.SetResultItem("trans_drop_packet",transmit_drop_point.value)

								results_packet = append(results_packet,tmp_packet)
							}
						}
					}
				}
			}
		}
	}
	return results_packet
}


func Container_Network_Bytes(address string,start int64,end int64,unit string,step int)([]ContainerNetworkAccResult){
	//container_network_transmit_bytes_total
	transmit_metric := "container_network_transmit_bytes_total"
	receive_metric := "container_network_receive_bytes_total"

	node_seconds_metric := "node_time_seconds"
	//quota_metric := "container_spec_cpu_quota"
	//cpu_perioid_metric := "container_spec_cpu_period"
	//
	transmit_url := Gen_range_url(address,transmit_metric,start,end,unit,step)
	node_second_url := Gen_range_url(address,node_seconds_metric,start,end,unit,step)
	receive_url := Gen_range_url(address,receive_metric,start,end,unit,step)

	transmit_result := Container_Network_Metric(transmit_url)
	receive_result := Container_Network_Metric(receive_url)

	node_second_result := Node_Raw_Metric(node_second_url)


	node_seconds := make(map[string]*NodeResult)
	var transmit map[string]map[string]map[string]map[string]*ContainerNetworkResult
	var receive map[string]map[string]map[string]map[string]*ContainerNetworkResult

	transmit = GenerateContainerNetworkMap(transmit_result)
	receive = GenerateContainerNetworkMap(receive_result)

	for k,_ := range node_second_result{
		node_seconds[node_second_result[k].GetNode()] = &node_second_result[k]
	}
	results_bytes := make([]ContainerNetworkAccResult,0,len(transmit_result))

	for k1,v1 := range transmit{
		_,ok11 := receive[k1]

		if (!ok11){
			continue
		}else{
			for k2,v2 := range v1{
				_,ok21 := receive[k1][k2]
				//fmt.Println(v2)
				if (!ok21){
					continue
				}else{
					for k3,v3 := range v2{
						_,ok31 := receive[k1][k2][k3]
						if (!ok31){
							continue
						}else{
							for k4,v4 := range v3{
								tmpreceive,ok41 := receive[k1][k2][k3][k4]
								node_name := strings.Split(v4.node,"-")
								tmp_time,ok44 := node_seconds[node_name[len(node_name)-1]]
								if (!ok41 || !ok44){
									continue
								}
								tmp_bytes := ContainerNetworkAccResult{}
								tmp_bytes.node = v4.GetNode()
								tmp_bytes.namespace = k1
								tmp_bytes.pod = k2
								tmp_bytes.container = k3
								tmp_bytes.inter = k4
								tmp_bytes.deploy = v4.GetDeploy()

								transpoint := Rawpoint{}
								recvpoint := Rawpoint{}


								transpoint.timestamp = v4.GetValuesi(len(v4.values)-1).timestamp
								totaltime := (tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
								int_total_time := int(totaltime)
								if int_total_time == 0 && totaltime < 0.001 && totaltime > -0.001{
									transpoint.value = 0
									recvpoint.value = 0
								}else{
									transpoint.value = (v4.GetValuesi(len(v4.values)-1).value - v4.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)
									recvpoint.value = (tmpreceive.GetValuesi(len(tmpreceive.values)-1).value - tmpreceive.GetValuesi(0).value)/(tmp_time.GetValuei(len(tmp_time.values)-1).value - tmp_time.GetValuei(0).value)

								}
								tmp_bytes.SetResultItem("recv_bytes",recvpoint.value)
								tmp_bytes.SetResultItem("trans_bytes",transpoint.value)

								results_bytes = append(results_bytes,tmp_bytes)
							}
						}
					}
				}
			}
		}
	}
	return results_bytes
}






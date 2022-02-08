package prometheus_controller

import (
	"client-autoscale/prometheus_tools"
	"fmt"
	"sync"
)

//func Container_CPU_Usage(address string,start int64,end int64,unit string,step int)([]ContainerAccResult){

func GenerateContainerMap(result[]prometheus_tools.ContainerAccResult,timestamp int64) map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult{
	mapresult := make(map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult)
	for k, _:= range result{
		tmpnamespace := result[k].GetNamespace()
		_,ok := mapresult[tmpnamespace]
		if !ok{
			mapresult[tmpnamespace] = make(map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult)
		}
		tmppod := result[k].GetPod()
		_,ok = mapresult[tmpnamespace][tmppod]
		if !ok{
			mapresult[tmpnamespace][tmppod] = make(map[string]map[int64]*prometheus_tools.ContainerAccResult)
		}
		tmpcontainer := result[k].GetContainer()
		_,ok = mapresult[tmpnamespace][tmppod][tmpcontainer]
		if !ok {
			mapresult[tmpnamespace][tmppod][tmpcontainer] = make(map[int64]*prometheus_tools.ContainerAccResult)
		}
		result[k].SetTimestamp(timestamp)

		mapresult[tmpnamespace][tmppod][tmpcontainer][int64(timestamp)] = &result[k]
		//}
	}
	return mapresult
}

func GenerateContainerNetworkMap(result[]prometheus_tools.ContainerNetworkAccResult,timestamp int64) map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult{
	mapresult := make(map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
	for k, _:= range result{
		tmpnamespace := result[k].GetNamespace()
		_,ok := mapresult[tmpnamespace]
		if !ok{
			mapresult[tmpnamespace] = make(map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
		}
		tmppod := result[k].GetPod()
		_,ok = mapresult[tmpnamespace][tmppod]
		if !ok{
			mapresult[tmpnamespace][tmppod] = make(map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
		}
		tmpcontainer := result[k].GetContainer()
		_,ok = mapresult[tmpnamespace][tmppod][tmpcontainer]
		if !ok {
			mapresult[tmpnamespace][tmppod][tmpcontainer] = make(map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
		}
		tmpinter := result[k].GetInter()
		_,ok = mapresult[tmpnamespace][tmppod][tmpcontainer][tmpinter]
		if !ok{
			mapresult[tmpnamespace][tmppod][tmpcontainer][tmpinter] = make(map[int64]*prometheus_tools.ContainerNetworkAccResult)
		}
		result[k].SetTimestamp(timestamp)
		mapresult[tmpnamespace][tmppod][tmpcontainer][tmpinter][int64(timestamp)] = &result[k]
		//}
	}
	return mapresult
}

func Container_Result(address string,duration int64,
	unit string,step int,res chan map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult,
	global_time chan int64,wg *sync.WaitGroup,
	getresult func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerAccResult)){
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered panic: %s\n",r)
		}
	}()
	defer func() {
		close(res)
	}()
	tmptime := int64(0)
	Loop:
		for{
			select {
			case t,ok := <-global_time:
				if !ok{
					break Loop
				}
				if t == tmptime{
					continue
				}else{
					result := getresult(address,t - duration,t,unit,step)
					mapresult := GenerateContainerMap(result,t)
					res <- mapresult
					wg.Done()
				}
			}
		}
}

func Container_Network_Result(address string,duration int64,
	unit string,step int,res chan map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult,
	global_time chan int64,wg *sync.WaitGroup,
	getresult func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerNetworkAccResult)){
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered panic: %s\n",r)
		}
	}()
	defer func() {
		close(res)
	}()
	tmptime := int64(0)
	Loop:
		for{
			select {
			case t,ok := <-global_time:
				if !ok{
					break Loop
				}
				if t == tmptime{
					continue
				}else{
					result := getresult(address,t - duration,t,unit,step)
					mapresult := GenerateContainerNetworkMap(result,t)
					res <- mapresult
					wg.Done()
				}
			}
		}
}



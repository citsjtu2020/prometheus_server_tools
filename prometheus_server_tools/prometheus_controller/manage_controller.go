package prometheus_controller

import (
	"prometheus_server_tools/prometheus_tools"
	"sync"
)

type ManageController struct {
	Global_time map[string]chan int64
	Container_func map[string]func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerAccResult)
	Container_net_func map[string]func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerNetworkAccResult)
	Container_res map[string]chan map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult
	Container_net_res map[string]chan map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult
	number int
	RM *sync.RWMutex
}

func InitMC() (mc *ManageController){
	mc = &ManageController{}
	mc.RM = new(sync.RWMutex)
	mc.Global_time = make(map[string]chan int64)
	mc.Container_res = make(map[string]chan map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult)
	mc.Container_net_res = make(map[string]chan map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult)
	mc.number = 0
	mc.Container_func  = make(map[string]func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerAccResult))
	mc.Container_net_func = make(map[string]func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerNetworkAccResult))
	return mc
}

func (mc *ManageController) RegisterContainerController(key string,
	capacity int,
	getresult func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerAccResult),
	update bool) bool{
	mc.RM.RLock()
	_,ok := mc.Global_time[key]
	mc.RM.RUnlock()
	if ok && !update{
		return false
	}else{
		ch := make(chan map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult,1)
		if capacity >= 0{
			ch = make(chan map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerAccResult,capacity)
		}
		mc.RM.Lock()
		mc.Container_res[key] = ch
		mc.Container_func[key] = getresult
		mc.Global_time[key] = make(chan int64,1)
		mc.number = mc.number + 1
		mc.RM.Unlock()
		return true
	}
}

func (mc *ManageController) DeleteContainerController(key string) bool{
	mc.RM.RLock()
	_,ok := mc.Container_res[key]
	mc.RM.RUnlock()
	if !ok{
		return false
	}else{
		mc.RM.Lock()
		close(mc.Global_time[key])
		mc.RM.Unlock()
		Loop0:
		for{
			select {
			case _,ok := <- mc.Container_res[key]:
				if !ok{
					break Loop0
				}
			default:
				continue
			}
		}
		mc.RM.Lock()
		delete(mc.Container_net_res,key)
		delete(mc.Global_time,key)
		delete(mc.Container_func,key)
		mc.number = mc.number-1
		mc.RM.Unlock()
		return true
	}
}

func (mc *ManageController) RegisterContainerNetworkController(key string,
	capacity int,
	getresult func(address string,start int64,end int64,unit string,step int)([]prometheus_tools.ContainerNetworkAccResult),
	update bool) bool{
	mc.RM.RLock()
	_,ok := mc.Global_time[key]
	mc.RM.RUnlock()
	if ok && !update{
		return false
	}else{
		ch := make(chan map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult,1)
		if capacity >= 0{
			ch = make(chan map[string]map[string]map[string]map[string]map[int64]*prometheus_tools.ContainerNetworkAccResult,capacity)
		}
		mc.RM.Lock()
		mc.Container_net_res[key] = ch
		mc.Global_time[key] = make(chan int64,1)
		mc.Container_net_func[key] = getresult
		mc.number = mc.number + 1
		mc.RM.Unlock()
		return true
	}
}

func (mc *ManageController) DeleteContainerNetworkController(key string) bool{
	mc.RM.RLock()
	_,ok := mc.Container_net_res[key]
	mc.RM.RUnlock()
	if !ok{
		return false
	}else{
		mc.RM.Lock()
		close(mc.Global_time[key])
		mc.RM.Unlock()
		Loop0:
		for{
			select {
			case _,ok := <-mc.Container_net_res[key]:
				if !ok{
					break Loop0
				}
			default:
				continue
			}
		}
		mc.RM.Lock()
		delete(mc.Container_net_res,key)
		delete(mc.Global_time,key)
		delete(mc.Container_net_func,key)
		mc.number = mc.number-1
		mc.RM.Unlock()
		return true
	}
}

func (mc *ManageController) GetNumber() int{
	return mc.number
}

func (mc *ManageController) Stop() bool{
	if mc.number == 0{
		return true
	}
	closed := make(map[string]struct{})
	mc.RM.RUnlock()
	for k,_ := range mc.Global_time{
		close(mc.Global_time[k])
		closed[k] = struct{}{}
	}
	mc.RM.RUnlock()
	last := len(closed)
	for ;last>0;{
		tmpclosed := make(map[string]struct{})
		for k,_ := range closed{
			mc.RM.RLock()
			_,ok1 := mc.Container_res[k]
			_,ok2 := mc.Container_net_res[k]
			if !ok1 && !ok2{
				tmpclosed[k] = struct{}{}
			}else if ok1 && !ok2{
				okchan := true
				select {
				case _,okchan = <-mc.Container_res[k]:
					if !okchan{
						tmpclosed[k] = struct{}{}
					}
				default:
					okchan = true
				}
			}else if !ok1 && ok2{
				okchan := true
				select {
				case _,okchan = <-mc.Container_net_res[k]:
					if !okchan{
						tmpclosed[k] = struct{}{}
					}
				default:
					okchan = true
				}
			}else{
				continue
			}
		}

		for k,_ := range tmpclosed{
			_,ok3 := closed[k]
			if ok3{
				delete(closed,k)
			}
		}
		last = len(closed)
	}
	return true
}






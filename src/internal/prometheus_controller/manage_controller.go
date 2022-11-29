package prometheus_controller

import (
	"context"
	"log"
	"prometheus_dump_tools/internal/config"
	"prometheus_dump_tools/internal/prometheus_tools"
	"sync"
	"time"
)

type Aggregator interface {
	Dump(ctx context.Context, conf *config.Config, unixTime int64)
}

type ManageController struct {
	GlobalTimeChan map[string]chan int64 //query start signal
	Aggregator     map[string]Aggregator
	number         int
	RM             *sync.RWMutex

	Conf *config.Config
}

// InitMC 用于初始化ManageController
func InitMC(conf *config.Config) (mc *ManageController) {
	mc = &ManageController{}
	mc.RM = new(sync.RWMutex)
	mc.GlobalTimeChan = make(map[string]chan int64)
	mc.number = 0
	mc.Aggregator = make(map[string]Aggregator)
	mc.Conf = conf
	return mc
}

func startDump(agg Aggregator, conf *config.Config, wg *sync.WaitGroup, globalTime chan int64) {

	for {
		select {
		case tickTime := <-globalTime:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(20))
			agg.Dump(ctx, conf, tickTime)
			cancel()
			wg.Done()
		default:
			continue
		}
	}
}

func (mc *ManageController) Start() {
	var wg sync.WaitGroup
	for key := range mc.Aggregator {
		mc.RM.RLock()
		agg := mc.Aggregator[key]
		timeChan := mc.GlobalTimeChan[key]
		copyConf := *mc.Conf
		mc.RM.RUnlock()
		go startDump(agg, &copyConf, &wg, timeChan)
	}
	timer := prometheus_tools.NewGlobalTicker("s", mc.Conf.Dumper.Interval)
	for {
		select {
		case tickTime := <-timer.C:
			tickTime = tickTime.Truncate(time.Duration(mc.Conf.Dumper.Interval) * time.Second)
			wg.Add(mc.GetNumber())
			for k := range mc.GlobalTimeChan {
				mc.GlobalTimeChan[k] <- tickTime.Unix()
			}
			wg.Wait()
		default:
			continue
		}
	}
}

func (mc *ManageController) RegisterMetricController(key string, agg Aggregator) {
	mc.RM.RLock()
	_, ok := mc.GlobalTimeChan[key]
	mc.RM.RUnlock()
	if ok {
		log.Fatal("repeat registering")
	} else {
		mc.RM.Lock()
		mc.Aggregator[key] = agg
		mc.GlobalTimeChan[key] = make(chan int64)
		mc.number = mc.number + 1
		mc.RM.Unlock()
		log.Printf("register success for %s", key)
	}
}

func (mc *ManageController) GetNumber() int {
	return mc.number
}

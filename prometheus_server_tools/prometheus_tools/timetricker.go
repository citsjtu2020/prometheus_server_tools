package prometheus_tools

import (
	"time"
)

func Gen_Timeout(unit string,duration int) time.Duration{
	timeout := time.Millisecond * time.Duration(duration)
	if unit == "s"{
		timeout = time.Second * time.Duration(duration)
	}else if unit == "min"{
		timeout = time.Minute * time.Duration(duration)
	}else if unit == "ms"{
		timeout = time.Millisecond * time.Duration(duration)
	}else if unit == "mu"{
		timeout = time.Microsecond * time.Duration(duration)
	}else if unit == "ns"{
		timeout = time.Nanosecond * time.Duration(duration)
	}
	return timeout
}

func NewGlobalTimer (unit string,duration int) *time.Timer{

	//fmt.Println(time.Now().Unix())
	//fmt.Println((time.Now().UnixNano()))
	//time.Sleep(timeout)
	//fmt.Println((time.Now().UnixNano()))
	timeout := Gen_Timeout(unit,duration)
	timer := time.NewTimer(timeout)
	//timer.Reset()
	return timer
}



//func (timer *time.Timer) ResetGlobalTimer(){
//	timer
//}

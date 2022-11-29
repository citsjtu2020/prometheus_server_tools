package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case tick := <-timer.C:
			fmt.Println(tick.Truncate(5 * time.Second))
		default:
			continue
		}

	}
}

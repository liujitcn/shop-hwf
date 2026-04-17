package main

import "time"

func main() {
	arr := make([]time.Duration, 0)
	arr = append(arr, time.Second*5)
	arr = append(arr, time.Second*10)
	arr = append(arr, time.Second*13)
	for _, duration := range arr {
		time.AfterFunc(duration, func() {
			println(time.Now().Format("2006-01-02 15:04:05"))
		})
	}
	select {}
}

package main

import (
	"fmt"
	"strings"
	"time"
)

func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func main() {

	Timeformat := "2006/01/02 15:04:05"
	curTime := time.Now()
	//times := strings.Split("2019/03/14", "/")
	//formatTimeStr := strings.Join(times[:3], "-") + " 00:00:00"
	//
	//fmt.Printf("%T, %v\n",formatTimeStr,formatTimeStr)
	//
	//standardActivityTime, err := time.Parse(Timeformat, formatTimeStr)
	//
	//fmt.Printf("%T, %v\n", standardActivityTime, standardActivityTime)
	//if err != nil {
	//	return
	//}
	//diff := timeSub(standardActivityTime, curTime)
	//fmt.Print(diff)

	times := strings.Split("2019/03/11/1", "/")
	keyTime := "2019/03/11/1"
	//formatTimeStr := strings.Join(times[:3], "/") + " 00:00:00"
	formatTimeStr := keyTime[:len(keyTime)-2] + " 00:00:00"
	standardActivityTime, _ := time.Parse(Timeformat, formatTimeStr)
	//if err != nil {
	//	return
	//}

	fmt.Println(times)
	fmt.Println(formatTimeStr)
	fmt.Println(standardActivityTime)

	diff := timeSub(standardActivityTime, curTime)
	fmt.Print(diff)
}

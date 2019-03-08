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
	curTime := time.Now()
	//formatTimeStr := "2019-03-01 00:00:00"
	keyTime := "2019/3/1/2"
	times := strings.Split(keyTime, "/")

	formatTimeStr := strings.Join(times[:3], "-0") + " 00:00:00"

	fmt.Println(formatTimeStr)
	standardActivityTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)

	//formatTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)

	if err != nil {

		fmt.Println(err)

	}

	//sub := timeSub(standardActivityTime, curTime)
	//
	//fmt.Printf("%T, %v", sub, sub)
	fmt.Println(curTime.Unix())
	fmt.Println(standardActivityTime.Unix())
	diff := float64(curTime.Unix()-standardActivityTime.Unix())/1000000.0
	fmt.Println(diff)
}

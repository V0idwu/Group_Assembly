package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	//registerTime := "2019/3/6/1"

	//registerTimeArray := strings.Split(registerTime, "/")
	curtimeStr:=time.Now().Format("2006-1-2 15:04:05")

	curdateArray := strings.Split(curtimeStr, " ")

	date := strings.Split(curdateArray[0], "-")
	time := strings.Split(curdateArray[1], ":")[0]
	timeint, err := strconv.Atoi(time)
	if err != nil {
		return "", err
	}
	if timeint > 6{

	}
	fmt.Println(date)
	fmt.Printf("%T, %v", time, time)
}

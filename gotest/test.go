package main

import (
	"fmt"
	"time"
)

func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func main() {

	a := []int{1,2,3,4,5}

	for k,v := range(a){
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}

	fmt.Println(len(a))

}

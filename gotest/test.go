package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

)

func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func httpPostForm() {
	resp, err := http.PostForm("http://10.141.221.88:36060/activityMatch",
		url.Values{})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}

var month2Num = make(map[string]int)


func main() {
	//ad, _ := time.Parse("2006-01-02", time.Now().String())
	t := time.Now()
	fmt.Println(time.Now())
	fmt.Println(t.Format("2006-01-02"))

}

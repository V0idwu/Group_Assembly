package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func initResources() []Resource {
	resources := []Resource{
		Resource{
			"1",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Shanghai",
			10,
			"tbd",
			"13:00",
			"14:00",
			1,
		},
		Resource{
			"1",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Shanghai",
			10,
			"tbd",
			"14:00",
			"15:00",
			1,
		},
		Resource{
			"1",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Shanghai",
			10,
			"tbd",
			"15:00",
			"16:00",
			1,
		},
		Resource{
			"1",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Shanghai",
			10,
			"tbd",
			"16:00",
			"17:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Shanghai",
			10,
			"tbd",
			"13:00",
			"14:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Shanghai",
			10,
			"tbd",
			"14:00",
			"15:00",
			1,
		},
		//Resource{
		//	"2",
		//	"Football",
		//	"Fudan Zhangjiang Campus Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"15:00",
		//	"16:00",
		//	1,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Fudan Zhangjiang Campus Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"16:00",
		//	"17:00",
		//	1,
		//},
		//Resource{
		//	"1",
		//	"Football",
		//	"Shanghai University of Traditional Chinese Medicine Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	20,
		//	"tbd",
		//	"13:00",
		//	"15:00",
		//	2,
		//},
		//Resource{
		//	"1",
		//	"Football",
		//	"Shanghai University of Traditional Chinese Medicine Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	20,
		//	"tbd",
		//	"15:00",
		//	"17:00",
		//	2,
		//},
		//Resource{
		//	"1",
		//	"Football",
		//	"Shanghai University of Science and Technology Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"13:00",
		//	"15:00",
		//	2,
		//},
		//Resource{
		//	"1",
		//	"Football",
		//	"Shanghai University of Science and Technology Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"15:00",
		//	"17:00",
		//	2,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Shanghai University of Science and Technology Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"13:00",
		//	"15:00",
		//	2,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Shanghai University of Science and Technology Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"15:00",
		//	"17:00",
		//	2,
		//},
		//Resource{
		//	"1",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	20,
		//	"tbd",
		//	"13:00",
		//	"15:00",
		//	2,
		//},
		//Resource{
		//	"1",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	20,
		//	"tbd",
		//	"15:00",
		//	"17:00",
		//	2,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"13:00",
		//	"14:00",
		//	1,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"14:00",
		//	"15:00",
		//	1,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"15:00",
		//	"16:00",
		//	1,
		//},
		//Resource{
		//	"2",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"16:00",
		//	"17:00",
		//	1,
		//},
		//Resource{
		//	"3",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"13:00",
		//	"14:00",
		//	1,
		//},
		//Resource{
		//	"3",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"14:00",
		//	"15:00",
		//	1,
		//},
		//Resource{
		//	"3",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"15:00",
		//	"16:00",
		//	1,
		//},
		//Resource{
		//	"3",
		//	"Football",
		//	"Zhangjiang Sports Center Football Field",
		//	"Zhangjiang Town",
		//	"Pudong District",
		//	"Shanghai",
		//	10,
		//	"tbd",
		//	"16:00",
		//	"17:00",
		//	1,
		//},
	}

	return resources
}

func initRequests() []Request {
	requests := []Request{}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 120; i++ {

		request := Request{}
		request.ID = strconv.Itoa(i + 1)
		//dateArr := []string{"2019-3-15", "2019-3-16"}
		request.Location = "Zhangjiang Town"
			request.RegisterTime = time.Now().Unix()
		//request.ActivityDate = dateArr[rand.Intn(2)]
		request.ActivityDate = time.Now().Format("2006-01-02")
		startTimeArr := []int{13, 14, 15, 16}
		endTimeArr := []int{14, 15, 16, 17}
		st := startTimeArr[rand.Intn(4)]
		et := endTimeArr[rand.Intn(4)]
		request.StartTime = strconv.Itoa(st) + ":00"
		for st >= et {
			et = endTimeArr[rand.Intn(4)]
		}
		request.EndTime = strconv.Itoa(et) + ":00"
		request.Deposit = rand.Intn(50)
		request.State = "0"
		request.ActivityType = "Football"
		//request.ResultID = "tbd"

		requests = append(requests, request)
	}
	//26,5,36,39,35,23,50,25,13,27
	//"StartTime":"14:00","Requests":[4,24,34,6,18,8,49,45,20,44]},
	return requests
}

func makeMatchGroup(resources []Resource, requests []Request, existMatchGroup []MatchGroup) []MatchGroup {
	newMatchGroups := generateNewMatchGroup(resources, requests)
	fmt.Println(len(newMatchGroups))
	//for _, matchGroup := range newMatchGroups {
	//	//size := unsafe.Sizeof(matchGroup)
	//	fmt.Println(matchGroup)
	//}
	matchGroups := checkExistMatchGroup(newMatchGroups, existMatchGroup)

	return matchGroups

}

// 合并新旧资源
func checkExistMatchGroup(newMatchGroups, existMatchGroups []MatchGroup) []MatchGroup {
	tempNewMatchGroups := []MatchGroup{}

	if len(existMatchGroups) != 0 {
		for _, existMatchGroup := range existMatchGroups {
			//fmt.Println("yes")
			for _, newMatchGroup := range newMatchGroups {
				// 活动日期，地点资源相同的matchgroup
				if newMatchGroup.ActivityDate == existMatchGroup.ActivityDate && reflect.DeepEqual(newMatchGroup.ResourcesInstance, existMatchGroup.ResourcesInstance) {
					// && newMatchGroup.StartTime == existMatchGroup.StartTime && newMatchGroup.StartTime == existMatchGroup.StartTime

					// 已存在的资源未匹配成功的活动
					if existMatchGroup.State == "2" {
						//把新组的人都放到老组中去，合并matchgroup
						for _, newRequest := range newMatchGroup.Requests {
							newRequest.State = "0"
							existMatchGroup.Requests = append(existMatchGroup.Requests, newRequest)
						}
					} else if existMatchGroup.State == "1" { // 所用资源已经被撮合完成，所以需要舍弃该组报名
						for _, dumpRequest := range newMatchGroup.Requests {
							dumpRequest.State = "2"
						}
					}
				} else {

					flag := true
					for _, temp := range tempNewMatchGroups {
						// 要加入的组是否已经在暂存区，资源相同，即都相同
						if reflect.DeepEqual(temp.ResourcesInstance, newMatchGroup.ResourcesInstance) {
							flag = false
						}
					}
					if flag {
						tempNewMatchGroups = append(tempNewMatchGroups, newMatchGroup)
					}

				}
			}
		}

		for _, tempNew := range tempNewMatchGroups {
			existMatchGroups = append(existMatchGroups, tempNew)
		}
		return existMatchGroups
	} else {
		//fmt.Println("yes")
		return newMatchGroups
	}
}

// 13:00 string --> 13 int
func turnHourTime2Int(time string) (int, error) {
	t := strings.Split(time, ":")[0]
	tint, err := strconv.Atoi(t)
	if err != nil {
		return -1, err
	}
	return tint, nil

}

func generateNewMatchGroup(resources []Resource, requests []Request) []MatchGroup {
	fmt.Println(len(requests))
	fmt.Println(len(resources))
	matchGroups := []MatchGroup{}

	ActivityDates := make(map[string]int)
	for _, resource := range requests {
		ActivityDates[resource.ActivityDate] ++
	}

	// 按活动日期进行分组
	for activityDate := range ActivityDates {
		//matchGroup := []MatchGroup{}
		for _, resource := range resources {
			singleMatch := MatchGroup{}
			for index, request := range requests {
				// 通过地点匹配
				if request.Location == resource.County {
					// 日期相同的报名
					if request.ActivityDate == activityDate {
						// 用户报名的开始结束时间要能包含资源可以提供的时间段
						qs, err := turnHourTime2Int(request.StartTime)
						ss, err := turnHourTime2Int(resource.StartTime)
						qe, err := turnHourTime2Int(request.EndTime)
						se, err := turnHourTime2Int(resource.EndTime)
						if err != nil {
							fmt.Println(err.Error())
						}
						if qs <= ss && qe >= se {
							request.State = "0"
							singleMatch.Requests = append(singleMatch.Requests, request)
						}

					}
				}

				if index == len(requests)-1 && len(singleMatch.Requests) != 0 {
					//fmt.Println(activityDate)
					singleMatch.ResourcesInstance = resource
					singleMatch.ResourcesInstance.ActivityDate = activityDate
					singleMatch.ActivityDate = activityDate
					singleMatch.State = "2"
					singleMatch.StartTime = resource.StartTime
					singleMatch.EndTime = resource.EndTime
					singleMatch.Duration = resource.Duration

				}
			}
			//size := unsafe.Sizeof(singleMatch.ResourcesInstance.Spot)
			//fmt.Println(size)
			if singleMatch.ResourcesInstance.Spot != "" {
				matchGroups = append(matchGroups, singleMatch)
			}
		}

	}

	// matchGroups去空

	return matchGroups
}

func prepare4MatchMakerservice(matchGroups []MatchGroup) ([]Request, []Resource) {
	resource4services := []Resource{}
	request4servicesDep := []Request{}

	for _, matchGroup := range matchGroups {
		resource4services = append(resource4services, matchGroup.ResourcesInstance)
		for _, request := range matchGroup.Requests {
			request4servicesDep = append(request4servicesDep, request)
		}
	}

	//fmt.Println(request4servicesDep)
	// 去重
	request4services := []Request{}
	for _, requestdep := range request4servicesDep {
		flag := true
		for _, request := range request4services {
			if reflect.DeepEqual(requestdep, request) {
				flag = false
			}
		}
		if flag {
			request4services = append(request4services, requestdep)
		}
	}

	for _, request4service := range request4services {
		fmt.Println(request4service)
	}
	for _, resource4service := range resource4services {
		fmt.Println(resource4service)
	}
	fmt.Println(len(request4services))
	fmt.Println(len(resource4services))

	return request4services, resource4services

}

func writeJson(data []byte, filename string) {
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	//fp.Truncate(0)
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

//func printBytes2JsonFile(bytes []byte, filename string){
//	data, err := json.Marshal(bytes)
//	if err != nil {
//		log.Fatalf("JSON marshaling failed: %s", err)
//	}
//	writeJson(data, filename)
//	//err = ioutil.WriteFile("./test.json", data, os.ModeAppend)
//	if err != nil {
//		return
//	}
//}

func httpPostForm(resources, requests []byte) ([]byte, error) {
	resp, err := http.PostForm("http://10.141.221.88:36060/activityMatch",
		url.Values{"resources": {string(requests)}, "requests": {string(resources)}})

	if err != nil {
		// handle error
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil, err
	}
	fmt.Println(string(body))

	return body, nil
}

func parseMatchMakingServiceResponse(){

}

func main() {
	requests := initRequests()
	resources := initResources()

	//fmt.Println(resources)
	//matchGroups := []MatchGroup{}
	//ActivityDate := make(map[string]int)
	//for _, request := range (requests) {
	//	fmt.Println(request)
	//}
	//
	//for k, v := range (ActivityDate) {
	//	fmt.Printf("%T : %v\n", k, k)
	//	fmt.Printf("%T : %v\n", v, v)
	//}
	//fmt.Println(requests)
	existMatchGroup := []MatchGroup{}
	matchGroups := makeMatchGroup(resources, requests, existMatchGroup)

	fmt.Println("\n")
	fmt.Println("========================================================")
	//fmt.Printf("%T \n", matchGroups)
	fmt.Println("Result:")
	fmt.Println("========================================================")
	fmt.Println()
	//for _, matchGroup := range matchGroups {
	//	//size := unsafe.Sizeof(matchGroup)
	//	fmt.Println(matchGroup)
	//}

	requestArr, resourceArr := prepare4MatchMakerservice(matchGroups)

	data1, err := json.Marshal(requestArr)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	writeJson(data1, "my-fabric\\chaincode\\Group_Assembly\\go\\requests.json")
	data2, err := json.Marshal(resourceArr)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	writeJson(data2, "my-fabric\\chaincode\\Group_Assembly\\go\\resources.json")

	httpPostForm(data1, data2)

}

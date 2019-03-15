package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
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
			"Zhanghai",
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
			"Zhanghai",
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
			"Zhanghai",
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
			"Zhanghai",
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
			"Zhanghai",
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
			"Zhanghai",
			10,
			"tbd",
			"14:00",
			"15:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"15:00",
			"16:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
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
			"Zhanghai",
			10,
			"tbd",
			"15:00",
			"16:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Fudan Zhangjiang Campus Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"16:00",
			"17:00",
			1,
		},
		Resource{
			"1",
			"Football",
			"Shanghai University of Traditional Chinese Medicine Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			20,
			"tbd",
			"13:00",
			"15:00",
			2,
		},
		Resource{
			"1",
			"Football",
			"Shanghai University of Traditional Chinese Medicine Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			20,
			"tbd",
			"15:00",
			"17:00",
			2,
		},
		Resource{
			"1",
			"Football",
			"Shanghai University of Science and Technology Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"13:00",
			"15:00",
			2,
		},
		Resource{
			"1",
			"Football",
			"Shanghai University of Science and Technology Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"15:00",
			"17:00",
			2,
		},
		Resource{
			"2",
			"Football",
			"Shanghai University of Science and Technology Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"13:00",
			"15:00",
			2,
		},
		Resource{
			"2",
			"Football",
			"Shanghai University of Science and Technology Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"15:00",
			"17:00",
			2,
		},
		Resource{
			"1",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			20,
			"tbd",
			"13:00",
			"15:00",
			2,
		},
		Resource{
			"1",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			20,
			"tbd",
			"15:00",
			"17:00",
			2,
		},
		Resource{
			"2",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"13:00",
			"14:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"14:00",
			"15:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"15:00",
			"16:00",
			1,
		},
		Resource{
			"2",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"16:00",
			"17:00",
			1,
		},
		Resource{
			"3",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"13:00",
			"14:00",
			1,
		},
		Resource{
			"3",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"14:00",
			"15:00",
			1,
		},
		Resource{
			"3",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"15:00",
			"16:00",
			1,
		},
		Resource{
			"3",
			"Football",
			"Zhangjiang Sports Center Football Field",
			"Zhangjiang Town",
			"Pudong District",
			"Zhanghai",
			10,
			"tbd",
			"16:00",
			"17:00",
			1,
		},
	}

	return resources
}

func initRequests() []Request {
	requests := []Request{}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {

		request := Request{}
		request.ID = strconv.Itoa(i + 1)
		dateArr := []string{"2019-3-15", "2019-3-16"}
		request.Location = "Zhangjiang Town"
		request.RegisterTime = time.Now().Unix()
		request.ActivityDate = dateArr[rand.Intn(2)]
		startTimeArr := []int{13, 14, 15, 16}
		endTimeArr := []int{14, 15, 16, 17}
		st := startTimeArr[rand.Intn(3)]
		et := endTimeArr[rand.Intn(3)]
		request.StartTime = strconv.Itoa(st)
		for st >= et {
			et = endTimeArr[rand.Intn(3)]
		}
		request.EndTime = strconv.Itoa(et)
		request.Deposit = rand.Intn(50)
		request.State = "2"
		request.ActivityType = "Football"
		//request.ResultID = "tbd"

		requests = append(requests, request)
	}

	return requests
}

func makeMatchGroup(resources []Resource, requests []Request, existMatchGroup []MatchGroup) []MatchGroup {
	newMatchGroups := generateNewMatchGroup(resources, requests)
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

func generateNewMatchGroup(resources []Resource, requests []Request) []MatchGroup {
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
						// 用户报名的开始结束时间要能包含
						if request.StartTime <= resource.StartTime && request.EndTime >= resource.EndTime {
							request.State = "0"
							singleMatch.Requests = append(singleMatch.Requests, request)
						}

					}
				}

				if index == len(requests)-1 && len(singleMatch.Requests) != 0 {
					//fmt.Println(activityDate)
					singleMatch.ResourcesInstance = resource
					//singleMatch.ResourcesInstance.ActivityDate = activityDate
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

func prepare4MatchMakerservice(matchGroups []MatchGroup) {
	resource4services := []Resource{}
	request4servicesDep := []Request{}

	for _, matchGroup := range matchGroups {
		resource4services = append(resource4services, matchGroup.ResourcesInstance)
		for _, request := range matchGroup.Requests {
			request4servicesDep = append(request4servicesDep, request)
		}
	}

	fmt.Println(request4servicesDep)
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

	fmt.Println(len(request4services))
	fmt.Println(request4services)
	fmt.Println(resource4services)
}

func main() {
	requests := initRequests()
	resources := initResources()

	//fmt.Println(resources)
	//matchGroups := []MatchGroup{}
	//ActivityDate := make(map[string]int)
	//for _, resource := range (requests) {
	//	ActivityDate[resource.ActivityDate] ++
	//}
	//
	//for k, v := range (ActivityDate) {
	//	fmt.Printf("%T : %v\n", k, k)
	//	fmt.Printf("%T : %v\n", v, v)
	//}
	//fmt.Println(requests)
	existMatchGroup := []MatchGroup{}
	matchGroups := makeMatchGroup(resources, requests, existMatchGroup)

	//fmt.Println("\n")
	//fmt.Println("========================================================")
	//fmt.Printf("%T \n", matchGroups)
	//fmt.Println("Result:")
	//fmt.Println("========================================================")
	//fmt.Println()
	//for _, matchGroup := range matchGroups {
	//	size := unsafe.Sizeof(matchGroup)
	//	fmt.Println(size,matchGroup)
	//}

	prepare4MatchMakerservice(matchGroups)
}

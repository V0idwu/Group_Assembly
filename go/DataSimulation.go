package gotest

import (
	"math/rand"
	"strconv"
	"time"

	"main"
)


func initResource() []Resource {
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

func initRequests() []Request{
	requests := []Request{{}}

	for i:=0 ; i<50 ; i++ {
		request := Request{}
		rand.Seed(int64(time.Now().UnixNano()))
		request.ID = strconv.Itoa(i+1)
		request.Location = "Zhangjiang Town"
		request.RegisterTime = time.Now().Unix()
		request.ActivityDate = "2019-3-15"
		startTimeArr := []int{13,14,15,16}
		endTimeArr := []int{14,15,16,17}
		st := startTimeArr[rand.Intn(3)]
		et := endTimeArr[rand.Intn(3)]
		request.StartTime = strconv.Itoa(st)
		for st > et{
			et = endTimeArr[rand.Intn(3)]
		}
		request.EndTime = strconv.Itoa(et)
		request.Deposit = rand.Intn(50)
		request.State = "0"
		request.ActivityType = "Football"
		request.ResultID = "tbd"

		requests = append(requests, request)
	}

	return requests
}

func main() {
	initRequests()
}

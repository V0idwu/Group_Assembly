package main

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	//"github.com/pkg/errors"
	"math"
	//"math/rand"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// --------------------------------------------
// 业务全局变量
// --------------------------------------------
var allParticipantsInActivity int = 0         // 所有参与人数
var idToSequenceNumber = make(map[int]string) // 用户ID对应到撮合时的序号

//var diffTime string
//var activityTime = make(map[string]int)
// --------------------------------------------
// 蚁群算法的全局常量
// --------------------------------------------
// 活动参与人数上下界
const M_LOW int = 4
const M_HIGH int = 4

// 有符号整型最小值
const MIN_VALUE = float64(^int(^uint(0) >> 1))

// 产生随机数前的休眠时间
const DURATION int = 3

// 全局变量
var applicantNum int    //活动人数
var deposit []float64   //每位报名者对应的押金
var timeValue []float64 //每位报名者当前对应的时间价值
var depositAndTime []float64

var iteratorNum = 1000 //迭代次数
var antNum = 10        //蚂蚁数量

var pheromoneMatrix [][]float64   //信息素矩阵
var maxPheromoneMatrix []int      //pheromoneMatrix矩阵的每一行中最大信息素的下标
var sortedPheromoneMatrix [][]int //pheromoneMatrix矩阵的每一行中信息素从大到小排列的下标
var criticalPointMatrix []int     //在一次迭代中，采用随机分配策略的蚂蚁的临界编号

var p float64 = 0.7 //每完成一次迭代后，信息素衰减的比例
var q float64 = 1.6 //蚂蚁每次经过一条路径，信息素增加的比例

// Define the Smart Contract structure
type SmartContract struct {
}

type User struct {
	ID    string
	Money int
}

type Request struct {
	ID           string
	Location     string //位置 zhangjiang Town
	RegisterTime int64  //客户端选择某一个上午，计算出那个上午的开始时间，再发给链码，这样方便以后修改时间选择策略
	ActivityDate string
	StartTime    string
	EndTime      string
	Deposit      int //押金
	// 报名的撮合状态
	// 0 进入matchgroup数组，还没有进行过第一次撮合
	// 1 停留在已撮合还未到参加活动时间的matchgroup组内，即撮合成功
	// 2 未撮合成功
	// 3 活动成功
	// 4 活动失败

	State          string
	ActivityType   string
	Owner          string
	ReqMatchResult string // Spot + SpotID + ActivityDate + StartTime + EndTime+ ActivityType
	// Fudan Zhangjiang Campus Football Field_1_2019-03-28_13:00_14:00_1
	//ResultID     string //被撮合到同一组别的用户会被分配一个相同的uid
}

type Resource struct {
	SpotID       string // 场地编号
	ActivityType string
	Spot         string
	County       string
	District     string
	City         string
	Capacity     int
	ActivityDate string // 固有资源，该项默认为"tbd"
	StartTime    string
	EndTime      string
	Duration     int
}

type MatchGroup struct {
	ActivityDate string // 活动日期
	Area         string // Spot + SpotID + ActivityDate + S
	StartTime    string
	EndTime      string
	Duration     int
	// MatchGroup的撮合状态
	// 1 撮合成功
	// 2 未撮合成功
	// 3 活动成功
	// 4 活动失败
	State             string
	Requests          []Request
	ResourcesInstance Resource // 加入活动日期信息
}

// 处理接收到的匹配结果
type MatchMakingResult struct {
	Area         string
	ActivityDate string
	StartTime    string
	EndTime      string
	State        string
	Requests     []int
}

type MatchGroupsByActivityDate map[string][]MatchGroup

// 工具方法：类型转换
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// String 和 []byte 相互转换
func stringtoslicebyte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func slicebytetostring(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: bh.Data,
		Len:  bh.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

//重写shim.ChaincodeStubInterface接口的 Init 方法
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//重写shim.ChaincodeStubInterface接口的 Invoke 方法
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	//获取用户意图与参数
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke is running " + function)
	//根据用户意图判断使用何种实现函数

	switch function {
	case "createRequest":
		return s.createRequest(stub, args)
	case "updateRequest":
		return s.createRequest(stub, args)
	case "cancelRequest":
		return s.cancelRequest(stub, args)
	case "confirmRequest":
		return s.confirmOrder(stub, args)
	//case "showAllRequest":
	//	return s.showAllRequest(stub)
	case "getAllLocationsToDapp":
		return s.getAllLocationsToDapp(stub)
	case "queryMyRequest":
		return s.queryMyRequest(stub)
	case "queryMyMoney":
		return s.queryMyMoney(stub)
	case "initLedger":
		return s.initLedger(stub, args)
	case "doMatchMaking":
		return s.doMatchMaking(stub, args)
	case "queryValueByKeyWithRegexSC":
		return s.queryValueByKeyWithRegexSC(stub, args)
	case "updateRequestsUponMatchGroups":
		return s.updateRequestsUponMatchGroups(stub, args)
	case "test":
		return s.test(stub, args)
	case "createResource":
		return s.createResource(stub, args)
	case "deleteResource":
		return s.deleteResource(stub, args)
	default:
		return shim.Error("no such method")
	}

}

// 初始化账本方法
func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 初始化资源到stateDB
	resources := initResources()
	err := setResources2Ledger(stub, resources)
	if err != nil {
		return shim.Error("Method setResources2Ledger " + err.Error())
	}

	var payload bytes.Buffer
	payload.WriteString("=== Resource into State DB === ")

	requests := initTestRequests()
	err = setRequests2Ledger(stub, requests)
	if err != nil {
		return shim.Error("Method setRequests2Ledger " + err.Error())
	}
	payload.WriteString("Requests into State DB ===")
	return shim.Success(payload.Bytes())

}

// =========================================================================================
// 按日期分组，进行撮合
// =========================================================================================
func (s *SmartContract) doMatchMaking(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	// request 模拟测试数据
	requests, err := getAvailableRequestsFromLedger(stub)
	if err != nil {
		return shim.Error("Method getAvailableRequestsFromLedger " + err.Error())
	}
	fmt.Println("doMatchMaking is running 1")
	// 从State DB 中读取resources
	//resources, err := getResourcesFromLedger(stub, []string{"Football","Basketball","Badminton"})
	resources, err := getResourcesFromLedger(stub)
	if err != nil {
		return shim.Error("Method getResourcesFromLedger " + err.Error())
	}
	fmt.Println("doMatchMaking is running 2")

	var payload bytes.Buffer
	payload.WriteString(" MatchMakingResult: ")

	var finalMatchResults []MatchGroup
	// 根据生成新的match group提供resource和request给撮合算法
	matchGroupsByDateType := generateNewMatchGroup(resources, requests)
	// 按日期，分组提交给撮合服务
	for activityType, matchGroupsByDate := range matchGroupsByDateType {
		payload.WriteString(" ActivityType: { ")
		payload.WriteString(activityType)

		for activityDate, matchGroups := range matchGroupsByDate {
			payload.WriteString(" ActivityDate: { ")
			payload.WriteString(activityDate)

			// 生成撮合算法的输入 resourceBytes和requestBytes
			requestArr, resourceArr := prepare4MatchMakerservice(matchGroups)
			requestBytes, err := json.Marshal(requestArr)
			if err != nil {
				return shim.Error("JSON marshaling failed: " + err.Error())
			}

			resourceBytes, err := json.Marshal(resourceArr)
			if err != nil {
				return shim.Error("JSON marshaling failed: " + err.Error())
			}

			if len(resourceArr) == 0 || len(requestArr) == 0 {
				//payload.WriteString("in loop ")
				continue
			}

			//访问撮合api，得到撮合结果

			//payload.WriteString("resourceBytes : ")
			//payload.WriteString(string(resourceBytes))
			//payload.WriteString("requestBytes : ")
			//payload.WriteString(string(requestBytes))

			matchMakingResultBytes, err := httpPostForm(resourceBytes, requestBytes)
			if err != nil {
				return shim.Error("MatchMaking HTTP : " + err.Error())
			}
			//payload.WriteString(string(matchMakingResultBytes))
			// 将撮合结果整合成matchgroup的样式
			matchMakingResults, err := parseMatchMakingServiceResponse(stub, matchMakingResultBytes, activityType)
			if err != nil {
				return shim.Error("Method parseMatchMakingServiceResponse " + err.Error())
			}
			for _, matchMakingResult := range matchMakingResults {
				payload.WriteString("--- Area: ")
				payload.WriteString(matchMakingResult.Area)
				payload.WriteString(" --- ")
				payload.WriteString("Requests: ")
				for _, req := range matchMakingResult.Requests {
					payload.WriteString(req.ID)
					payload.WriteString(",")
				}
				payload.WriteString(" --- ")
			}
			// 检查是否是已存的matchgroup使用了同一片资源
			finalMatchMakerResult, err := checkExistMatchGroup(stub, matchMakingResults)
			fmt.Println("finalMatchMakerResult: ")
			fmt.Println(finalMatchMakerResult)
			if err != nil {
				return shim.Error("Method checkExistMatchGroup " + err.Error())
			}
			for _, finalMatchMakerR := range finalMatchMakerResult {
				finalMatchResults = append(finalMatchResults, finalMatchMakerR)
			}
			payload.WriteString(" } ")
		}
		payload.WriteString(" } ")
	}
	// 将最终结果matchgroup存入stabe db
	fmt.Println("==========================================")
	fmt.Println("finalMatchMakerResult: ")
	fmt.Println("==========================================")
	fmt.Println(finalMatchResults)

	err = setMatchGroups2Ledger(stub, finalMatchResults)
	if err != nil {
		return shim.Error("Method setMatchGroups2Ledger " + err.Error())
	}
	payload.WriteString(" Save MatchMaking Result into StateDB ")
	return shim.Success(payload.Bytes())
}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) test(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	//footballResources := []Resource{}
	var payload bytes.Buffer

	data, err := queryValueByKeyWithRegex(stub, args)
	if err != nil {
		shim.Error(err.Error())
	}
	payload.WriteString(string(data))
	//resources, err := queryResourcesByOneKey(stub, []string{args[0],args[1]})
	//if err != nil {
	//	shim.Error(err.Error())
	//}
	//payload.WriteString("len(resources): ")
	//payload.WriteString(strconv.Itoa(len(resources)))
	////for _, resource := range resources {
	////	var r Resource
	////	r = resource.(Resource)
	////	footballResources = append(footballResources, r)
	////}
	//
	//for _,re := range resources{
	//	payload.WriteString(re.Spot)
	//	payload.WriteString("_")
	//	payload.WriteString(re.SpotID)
	//}

	return shim.Success(payload.Bytes())
}

//接口方法：发布请求
func (s *SmartContract) createRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	existRequest, _ := stub.GetState(args[0])
	if existRequest != nil {
		return shim.Error("Exist Request. ")
	}

	var owner, _ = GetCertAttribute2(stub)
	var userId = "user" + owner
	userAsBytes, _ := stub.GetState(userId)
	if userAsBytes == nil {
		var user = User{ID: userId, Money: 10000}
		userAsBytes, _ = json.Marshal(user)
		stub.PutState(userId, userAsBytes)
	}

	user := User{}
	json.Unmarshal(userAsBytes, &user) //unmarshal it aka JSON.parse()
	money := user.Money
	deposit, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error(err.Error())
	}
	if money < deposit {
		var payload bytes.Buffer
		payload.WriteString("Not enough money")
		return shim.Success(payload.Bytes())
	}
	user.Money = money - deposit
	userAsBytes, _ = json.Marshal(user)
	stub.PutState(userId, userAsBytes)

	var request = Request{ID: args[0], Location: args[1], RegisterTime: time.Now().Unix(), ActivityDate: args[2], StartTime: args[3], EndTime: args[4], Deposit: deposit, State: "0", ActivityType: args[6], Owner: owner}
	requestAsBytes, _ := json.Marshal(request)

	stub.PutState(args[0], requestAsBytes)

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("  Register Success")

	// 参与人数
	applicantNum = applicantNum + 1
	return shim.Success(payload.Bytes())
}

//接口方法：修改请求
func (s *SmartContract) updateRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	//判断该请求是否是该用户所有
	//获取该请求真实所有者
	requestID := args[0]
	requestAsBytes, err := stub.GetState(requestID)
	if err != nil {
		return shim.Error("Failed to Get Request:" + err.Error())
	}
	if requestAsBytes == nil {
		return shim.Error("Request does not Exist")
	}
	request := Request{}
	err = json.Unmarshal(requestAsBytes, &request) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	deposit, err := strconv.Atoi(args[5])
	request = Request{ID: args[0], Location: args[1], ActivityDate: args[2], StartTime: args[3], EndTime: args[4], Deposit: deposit, State: args[6]}
	requestAsBytes, _ = json.Marshal(request)
	stub.PutState(args[0], requestAsBytes)

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("  Update Success")

	return shim.Success(payload.Bytes())
}

//接口方法：取消请求
func (s *SmartContract) cancelRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	//判断该请求是否是该用户所有
	//获取该请求真实所有者
	requestID := args[0]
	requestAsBytes, err := stub.GetState(requestID)
	if err != nil {
		return shim.Error("Failed to get request:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("request does not exist")
	}
	request := Request{}
	err = json.Unmarshal(requestAsBytes, &request) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	owner := request.Owner
	//获取该交易用户名
	var userName, _ = GetCertAttribute2(stub)
	if owner != userName {
		return shim.Error("Error User")
	}
	//判断是否已被撮合
	if request.State != "0" {
		return shim.Error("Can not Cancel. You has been Arranged into an Activity")
	}

	var userId = "user" + owner
	userAsBytes, _ := stub.GetState(userId)

	user := User{}
	json.Unmarshal(userAsBytes, &user) //unmarshal it aka JSON.parse()
	money := user.Money
	user.Money = money + request.Deposit
	userAsBytes, _ = json.Marshal(user)
	stub.PutState(userId, userAsBytes)

	stub.DelState(args[0])

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("  Cancel Success")

	return shim.Success(payload.Bytes())
}

//func (s *SmartContract) showAllRequest(stub shim.ChaincodeStubInterface) sc.Response {
//
//	requestAsBytes, err := queryRequestValueByKeyWithRegex(stub, []string{"StartTime", ""})
//	if err != nil {
//		return shim.Error("Failed to get request:" + err.Error())
//	} else if requestAsBytes == nil {
//		return shim.Error("request does not exist")
//	}
//	return shim.Success(requestAsBytes)
//}

// 活动结束后用于确认请求
func (s *SmartContract) confirmOrder(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//判断该请求是否是该用户所有
	//获取该请求真实所有者
	requestID := args[0]
	requestAsBytes, err := stub.GetState(requestID)
	if err != nil {
		return shim.Error("Failed to get request:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("request does not exist")
	}
	request := Request{}
	err = json.Unmarshal(requestAsBytes, &request) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	owner := request.Owner
	//获取该交易用户名
	var userName, _ = GetCertAttribute2(stub)
	if owner != userName {
		return shim.Error("Error User")
	}
	//判断是否已被撮合
	if request.State != "1" {
		return shim.Error("Can not update. ")
	}

	if args[1] == "3" {
		//如果去，就把1改成3，并还钱
		request.State = "3"
		user := User{}

		var userId = "user" + owner
		userAsBytes, _ := stub.GetState(userId)
		json.Unmarshal(userAsBytes, &user) //unmarshal it aka JSON.parse()
		money := user.Money
		user.Money = money + request.Deposit
		userAsBytes, _ = json.Marshal(user)
		stub.PutState(userId, userAsBytes)
	} else if args[1] == "4" {
		//如果是不去，就把1改为4
		request.State = "4"
	}

	requestAsBytes, _ = json.Marshal(request)
	stub.PutState(args[0], requestAsBytes)

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("Updating Success")

	return shim.Success(payload.Bytes())
}

// 查询我的所有请求
func (s *SmartContract) queryMyRequest(stub shim.ChaincodeStubInterface) sc.Response {

	var owner, _ = GetCertAttribute2(stub)
	queryString := fmt.Sprintf("{\"selector\":{\"Owner\":\"%s\"}}", owner)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// 查询我的余额
func (s *SmartContract) queryMyMoney(stub shim.ChaincodeStubInterface) sc.Response {

	var owner, _ = GetCertAttribute2(stub)

	var userId = "user" + owner
	userAsBytes, _ := stub.GetState(userId)
	if userAsBytes == nil {
		var user = User{ID: userId, Money: 10000}
		userAsBytes, _ = json.Marshal(user)
		stub.PutState(userId, userAsBytes)
	}

	return shim.Success(userAsBytes)
}

// 智能合约：获取所有地址详情
func (s *SmartContract) getAllLocationsToDapp(stub shim.ChaincodeStubInterface) sc.Response {
	locationMap, err := getAllRequestValueNum(stub, "Location")
	if err != nil {
		return shim.Error(err.Error())
	}
	if unsafe.Sizeof(locationMap) == 0 {
		return shim.Error("Empty Locations")
	}

	var payload bytes.Buffer
	for k, v := range locationMap {
		payload.WriteString(k)
		payload.WriteString(",")
		payload.WriteString(strconv.Itoa(v))
		payload.WriteString(";")
	}
	return shim.Success(payload.Bytes())
}

// 添加资源
func (s *SmartContract) createResource(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	capacity, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error(" createResource " + err.Error())
	}
	duration, err := strconv.Atoi(args[9])
	if err != nil {
		return shim.Error(" createResource " + err.Error())
	}
	resource := Resource{args[0], args[1], args[2], args[3], args[4], args[5], capacity, "tbd", args[7], args[8], duration}

	key, err := stub.CreateCompositeKey("Resource", []string{resource.Spot, resource.SpotID, resource.ActivityType, resource.StartTime, resource.EndTime})
	if err != nil {
		return shim.Error(" createResource " + err.Error())
	}
	value, err := json.Marshal(resource)
	if err != nil {
		return shim.Error(" createResource " + err.Error())
	}
	err = stub.PutState(key, value)
	if err != nil {
		return shim.Error(" createResource " + err.Error())
	}

	var payload bytes.Buffer
	payload.WriteString(" Resource Create Successfully ")
	return shim.Success(payload.Bytes())
}

// 删除资源
func (s *SmartContract) deleteResource(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	// resource.Spot, resource.SpotID, resource.ActivityType, resource.StartTime, resource.EndTime
	key, err := stub.CreateCompositeKey("Resource", []string{args[1], args[0], args[2], args[3], args[4]})
	if err != nil {
		return shim.Error(" deleteResource " + err.Error())
	}
	data, err := stub.GetState(key)
	if err != nil {
		return shim.Error(" deleteResource " + err.Error())
	}
	if len(string(data)) == 0 {
		return shim.Error(" Resource doesn't Exist")
	}
	err = stub.DelState(key)
	if err != nil {
		return shim.Error(" deleteResource " + err.Error())
	}
	var payload bytes.Buffer
	payload.WriteString(" Resource Delete Successfully ")
	return shim.Success(payload.Bytes())
}

// 查询相同请求人数
//func (s *SmartContract) querySameRequestPeopleNumber(stub shim.ChaincodeStubInterface, args []string) sc.Response {
//
//	queryResults := s.queryRequestValueByKey(stub, args)
//
//	//获取人数
//	return shim.Success(queryResults.Payload)
//}

//
////查询押金平均值
//func (s *SmartContract) querySameRequestAverageDeposit(stub shim.ChaincodeStubInterface, args []string) sc.Response {
//
//	queryResults, err := querySameRequestByRequestID(stub, args)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//
//	//获取押金平均值
//	return shim.Success(queryResults)
//}
//
////接口方法：用来测试查询用户名称
//func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface) sc.Response {
//
//	var owner, _ = GetCertAttribute2(APIstub)
//	var data []byte = []byte(owner)
//
//	return shim.Success(data)
//}

// --------------------------------------------------------
// 返回报名信息中各种字段的数量
// --------------------------------------------------------
func getAllRequestValueNum(stub shim.ChaincodeStubInterface, args string) (map[string]int, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":{\"$regex\":\"%s\"}}}", args, "")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	valueNum := make(map[string]int)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		request := Request{}
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		switch args {
		case "Location":
			valueNum[request.Location]++
		case "StartTime":
			valueNum[request.StartTime]++
		default:
			return nil, errors.New("No such value In DB")
		}

	}
	return valueNum, nil
}

/*
通用方法：
获取证书属性1
获取证书属性2
根据过滤查询相同请求
*/

//工具方法：获取证书属性"username" "role" 方法来自于：https://www.ibm.com/developerworks/cn/cloud/library/cl-ibm-blockchain-chaincode-development-using-golang/index.html
// func GetCertAttribute1(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
// 	fmt.Println("Entering GetCertAttribute")
// 	attr, err := stub.ReadCertAttribute(attributeName)
// 	if err != nil {
// 		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
// 	}
// 	attrString := string(attr)
// 	return attrString, nil
// }【【

//工具方法：获取证书属性"username" 方法来自于：http://www.cnblogs.com/studyzy/p/7360733.html
func GetCertAttribute2(stub shim.ChaincodeStubInterface) (string, error) {
	creatorByte, err := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
		fmt.Errorf("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	fmt.Println("Name:" + uname)
	return uname, nil
}

// =========================================================================================
//工具方法：查询state db中的key和value
// =========================================================================================
func queryRequestValueByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}}", args[0], args[1])
	requestAsBytes, err := getQueryResultForQueryString(stub, queryString)

	return requestAsBytes, err
}

// =========================================================================================
//工具方法：查询state db中的key和values，使用正则表达式
// =========================================================================================
func queryValueByKeyWithRegex(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":{\"$regex\":\"%s\"}}}", args[0], args[1])
	requestAsBytes, err := getQueryResultForQueryString(stub, queryString)

	return requestAsBytes, err
}

func (s *SmartContract) queryValueByKeyWithRegexSC(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":{\"%s\":\"%s\"}}", args[0], args[1], args[2])
	requestAsBytes, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(requestAsBytes)
}

func queryMatchGroupsByOneKey(stub shim.ChaincodeStubInterface, args []string) ([]MatchGroup, error) {
	if len(args) != 2 {
		return nil, errors.New("Method queryValueByOneKey() Incorrect number of arguments. Expecting 2")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":{\"$regex\":\"%s\"}}}", args[0], args[1])
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var values []MatchGroup
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var value MatchGroup
		err = json.Unmarshal(queryResponse.Value, &value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}
	return values, err
}

func queryRequestValueByTwoKey(stub shim.ChaincodeStubInterface, args []string) ([]Request, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\",\"%s\":{\"$regex\":\"%s\"}}}", args[0], args[1], args[2], args[3])
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	requests := []Request{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		request := Request{}
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}
	return requests, err
}

func queryResourcesByOneKey(stub shim.ChaincodeStubInterface, args []string) ([]Resource, error) {
	if len(args) != 2 {
		return nil, errors.New("Method queryValueByOneKey() Incorrect number of arguments. Expecting 2")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":{\"$regex\":\"%s\"}}}", args[0], args[1])
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var values []Resource
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var value Resource
		err = json.Unmarshal(queryResponse.Value, &value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}
	return values, err
}

func queryResourcesValueByTwoKey(stub shim.ChaincodeStubInterface, args []string) ([]Resource, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\",\"%s\":\"%s\"}}", args[0], args[1], args[2], args[3])
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	resources := []Resource{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		resource := Resource{}
		err = json.Unmarshal(queryResponse.Value, &resource)
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}
	return resources, err
}

// =========================================================================================
// 工具方法：富文本查询
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// =========================================================================================
// 对当前state db中的报名信息进行整理，
// =========================================================================================
// func filter(stub shim.ChaincodeStubInterface, satisMaxArr *[]float64, satisMaxPath *[][][]int) error {
// 	locationNumMap, err := getAllRequestValueNum(stub, "Location")
// 	if err != nil {
// 		return err
// 	}
// 	activityTimeNumMap, err := getAllRequestValueNum(stub, "ActivityTime")

// 	//requests := []Request{}
// 	for keyLoc, _ := range locationNumMap {
// 		for keyTime, _ := range activityTimeNumMap {
// 			curTime := time.Now()
// 			times := strings.Split(keyTime, "/")
// 			formatTimeStr := strings.Join(times[:3], "-0") + " 00:00:00"
// 			standardActivityTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)
// 			if err != nil {
// 				return err
// 			}
// 			diff := timeSub(standardActivityTime, curTime)

// 			if diff >= 3 {
// 				values := []string{"Location", keyLoc, "ActivityTime", keyTime}
// 				requests, err := queryRequestValueByTwoKey(stub, values)
// 				if err != nil {
// 					return err
// 				}
// 				//requests = append(requests, request)
// 				err1 := initUserRegisterInfo(requests)
// 				if err1 != nil {
// 					return err1
// 				}
// 				////satisMax, satisPath, err2 := aca()
// 				//if err2 != nil {
// 				//	return err2
// 				//}
// 				//*satisMaxArr = append(*satisMaxArr, satisMax)
// 				//*satisMaxPath = append(*satisMaxPath, satisPath)
// 			}

// 		}

// 	}

// 	return nil
// }

// 计算时间差的天数
func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func initResources() []Resource {
	resources := []Resource{
		Resource{
			"1",
			"1",
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
			"1",
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
			"1",
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
			"1",
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
			"3",
			"Fudan Zhangjiang Campus Basketball Field",
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
			"3",
			"Fudan Zhangjiang Campus Basketball Field",
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
			"2",
			"1",
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
			"2",
			"1",
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

func initTestRequests() []Request {
	requests := []Request{}

	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {

		request := Request{}
		request.ID = strconv.Itoa(i + 1)

		request.Location = "Zhangjiang Town"
		request.RegisterTime = 100
		//dateArr := []string{"2019-3-15", "2019-3-16"}
		//request.ActivityDate = dateArr[rand.Intn(2)]
		request.ActivityDate = time.Now().Format("2006-01-02")
		//startTimeArr := []int{13, 14, 15, 16}
		//endTimeArr := []int{14, 15, 16, 17}
		//st := startTimeArr[rand.Intn(4)]
		//et := endTimeArr[rand.Intn(4)]
		st := 13
		et := 14
		request.StartTime = strconv.Itoa(st) + ":00"
		//for st >= et {
		//	et = endTimeArr[rand.Intn(4)]
		//}
		request.EndTime = strconv.Itoa(et) + ":00"
		//request.Deposit = rand.Intn(50)
		request.Deposit = 50
		request.State = "0"
		request.ActivityType = "1"
		//request.ResultID = "tbd"
		request.ReqMatchResult = "noResult"
		requests = append(requests, request)
	}

	for j := 10; j < 20; j++ {

		request := Request{}
		request.ID = strconv.Itoa(j + 1)

		request.Location = "Zhangjiang Town"
		request.RegisterTime = 100
		//dateArr := []string{"2019-3-15", "2019-3-16"}
		//request.ActivityDate = dateArr[rand.Intn(2)]
		request.ActivityDate = "2019-03-29"
		//startTimeArr := []int{13, 14, 15, 16}
		//endTimeArr := []int{14, 15, 16, 17}
		//st := startTimeArr[rand.Intn(4)]
		//et := endTimeArr[rand.Intn(4)]
		st := 13
		et := 14
		request.StartTime = strconv.Itoa(st) + ":00"
		//for st >= et {
		//	et = endTimeArr[rand.Intn(4)]
		//}
		request.EndTime = strconv.Itoa(et) + ":00"
		//request.Deposit = rand.Intn(50)
		request.Deposit = 50
		request.State = "0"
		request.ActivityType = "3"
		//request.ResultID = "tbd"
		request.ReqMatchResult = "noResult"
		requests = append(requests, request)
	}
	//26,5,36,39,35,23,50,25,13,27
	//"StartTime":"14:00","Requests":[4,24,34,6,18,8,49,45,20,44]},
	return requests
}

func getMatchGroupsFromLedger(stub shim.ChaincodeStubInterface) ([]MatchGroup, error) {
	matchGroups, err := queryMatchGroupsByOneKey(stub, []string{"Area", ""})
	if err != nil {
		return nil, err
	}
	return matchGroups, nil
}

func setResources2Ledger(stub shim.ChaincodeStubInterface, resources []Resource) error {

	for _, resource := range resources {
		key, err := stub.CreateCompositeKey("Resource", []string{resource.Spot, resource.SpotID, resource.ActivityType, resource.StartTime, resource.EndTime})
		if err != nil {
			return err
		}
		value, err := json.Marshal(resource)
		if err != nil {
			return err
		}
		err = stub.PutState(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// get request
func getResourcesFromLedger(stub shim.ChaincodeStubInterface) ([]Resource, error) {

	//allResources := []Resource{}
	//for i := range args {
	//	resources, err := queryResourcesByOneKey(stub, []string{"ActivityType", args[i]})
	//	if err != nil {
	//		return nil, err
	//	}
	//	for _, resource := range resources {
	//		allResources = append(allResources, resource)
	//	}
	//}

	resources, err := queryResourcesByOneKey(stub, []string{"ActivityType", ""})
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func setRequests2Ledger(stub shim.ChaincodeStubInterface, requests []Request) error {
	for _, request := range requests {
		requestData, err := json.Marshal(request)
		if err != nil {
			return err
		}
		err = stub.PutState(request.ID, requestData)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAvailableRequestsFromLedger(stub shim.ChaincodeStubInterface) ([]Request, error) {
	requestsNotMatch, err := queryRequestValueByTwoKey(stub, []string{"State", "0", "Owner", ""})
	if err != nil {
		return nil, err
	}
	requestsFailMatch, err := queryRequestValueByTwoKey(stub, []string{"State", "2", "Owner", ""})
	if err != nil {
		return nil, err
	}
	for _, request := range requestsFailMatch {
		requestsNotMatch = append(requestsNotMatch, request)
	}
	return requestsNotMatch, err
}

func setMatchGroups2Ledger(stub shim.ChaincodeStubInterface, matchGroups []MatchGroup) error {

	for _, matchgroup := range matchGroups {

		if len(matchgroup.Requests) != 0 {
			value, err := json.Marshal(matchgroup)
			if err != nil {
				return err
			}
			fmt.Println(" Value=================")
			fmt.Println(value)
			err = stub.PutState(matchgroup.Area, value)
			fmt.Println("error : ", err.Error())
			if err != nil {
				return err
			}
		}else {
			fmt.Println(" matchgroup is empty")
		}

	}
	return nil
}

// 根据matchgroup的情况，更新requests
func (s *SmartContract) updateRequestsUponMatchGroups(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	matchGroups, err := getMatchGroupsFromLedger(stub)
	if err != nil {
		return shim.Error("updateRequestsUponMatchGroups " + err.Error())
	}

	var requests []Request
	for _, matchGroup := range matchGroups {
		for _, request := range matchGroup.Requests {
			var req Request
			data, err := stub.GetState(request.ID)
			if err != nil {
				return shim.Error("updateRequestsUponMatchGroups " + err.Error())
			}
			err = json.Unmarshal(data, &req)
			if err != nil {
				return shim.Error("updateRequestsUponMatchGroups " + err.Error())
			}
			requests = append(requests, req)
		}
	}

	for _, matchGroup := range matchGroups {
		for _, request := range requests {
			if request.ReqMatchResult == matchGroup.Area{
				if matchGroup.State == "0" {
					request.State = "1"
				}
			}

			dataIn, err := json.Marshal(request)
			if err != nil {
				return shim.Error("updateRequestsUponMatchGroups " + err.Error())
			}
			err = stub.PutState(request.ID, dataIn)
			if err != nil {
				return shim.Error("updateRequestsUponMatchGroups " + err.Error())
			}
		}
		matchGroup.State = "1"
	}
	var payload bytes.Buffer
	payload.WriteString("updateRequestsUponMatchGroups ")
	err = setMatchGroups2Ledger(stub, matchGroups)
	if err != nil {
		return shim.Error("updateRequestsUponMatchGroups " + err.Error())
	}
	return shim.Success(payload.Bytes())
}

//查看资源是否已经存在占用情况
func checkExistMatchGroup(stub shim.ChaincodeStubInterface, newMatchGroups []MatchGroup) ([]MatchGroup, error) {
	var matchGroups []MatchGroup
	for _, newMatchGroup := range newMatchGroups {
		existMatchGroup, err := stub.GetState(newMatchGroup.Area)
		if existMatchGroup != nil {
			continue
		}
		if existMatchGroup == nil {
			matchGroups = append(matchGroups, newMatchGroup)
		} else if err != nil {
			return nil, err
		}

	}
	return matchGroups, nil
}

// xx:00 string --> xx int
func turnHourTime2Int(time string) (int, error) {
	t := strings.Split(time, ":")[0]
	tint, err := strconv.Atoi(t)
	if err != nil {
		return -1, err
	}
	return tint, nil

}

func generateNewMatchGroup(resources []Resource, requests []Request) map[string]MatchGroupsByActivityDate {
	activityDates := make(map[string]int)
	for _, request := range requests {
		activityDates[request.ActivityDate]++
	}

	activityTypes := make(map[string]int)
	for _, request := range requests {
		activityTypes[request.ActivityType]++
	}
	matchGroupsByActivityDateType := map[string]MatchGroupsByActivityDate{}
	// 按活动日期进行分组
	for activityType := range activityTypes {
		matchGroupsByDate := map[string][]MatchGroup{}
		// 按活动日期进行分组
		for activityDate := range activityDates {
			matchGroup := []MatchGroup{}
			for _, resource := range resources {
				singleMatch := MatchGroup{}
				for index, request := range requests {
					// 通过地点匹配
					if request.Location == resource.County {
						// 日期相同的报名
						if request.ActivityDate == activityDate {
							if request.ActivityType == activityType && resource.ActivityType == activityType {
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
					matchGroup = append(matchGroup, singleMatch)
				}
			}
			matchGroupsByDate[activityDate] = matchGroup
		}
		matchGroupsByActivityDateType[activityType] = matchGroupsByDate
	}
	return matchGroupsByActivityDateType
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
		url.Values{"resources": {string(resources)}, "requests": {string(requests)}})

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
	//fmt.Println(string(body))

	return body, nil
}

// 将从撮合服务得到的结果转换成[]MatchGroup
func parseMatchMakingServiceResponse(stub shim.ChaincodeStubInterface, matchMakingResults []byte, activityType string) ([]MatchGroup, error) {
	var matchMakingResult []MatchMakingResult

	if string(matchMakingResults) == "" {
		return nil, errors.New("No matchmaking results")
	}

	fmt.Println(matchMakingResults)

	err := json.Unmarshal(matchMakingResults, &matchMakingResult)
	if err != nil {
		return nil, errors.New("step 1")
	}

	// fmt.Printf("%T : %v", matchMakingArr,matchMakingArr)
	//fmt.Println(matchGroups)
	matchGroups := []MatchGroup{}
	for _, matchMaking := range matchMakingResult {
		matchGroup := MatchGroup{}
		spot := strings.Split(matchMaking.Area, "_")[0]
		spotID := strings.Split(matchMaking.Area, "_")[1]
		matchGroup.ActivityDate = matchMaking.ActivityDate
		matchGroup.Area = matchMaking.Area + ":00" + "_" + matchMaking.EndTime + "_" + activityType
		matchGroup.StartTime = matchMaking.StartTime
		matchGroup.EndTime = matchMaking.EndTime
		matchGroup.State = "0"
		startint, err := strconv.Atoi(strings.Split(matchMaking.StartTime, ":")[0])
		if err != nil {
			return nil, errors.New("step 2")
		}
		endint, err := strconv.Atoi(strings.Split(matchMaking.EndTime, ":")[0])
		if err != nil {
			return nil, errors.New("step 3")
		}
		if endint-startint > 0 {
			matchGroup.Duration = endint - startint
		}
		for _, requestID := range matchMaking.Requests {
			request := Request{}
			data, err := stub.GetState(strconv.Itoa(requestID))
			if err != nil {
				return nil, errors.New("step 4")
			}
			err = json.Unmarshal(data, &request)
			if err != nil {
				return nil, errors.New("step 5")
			}
			// 在此处将用户报名的State置为"1"
			request.State = "1"
			matchGroup.Requests = append(matchGroup.Requests, request)
		}

		resourceKey, err := stub.CreateCompositeKey("Resource", []string{spot, spotID, activityType, matchGroup.StartTime, matchGroup.EndTime})
		if err != nil {
			return nil, err
		}
		resource := Resource{}
		data, err := stub.GetState(resourceKey)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &resource)
		if err != nil {
			return nil, errors.New("step 6")
		}
		matchGroup.ResourcesInstance = resource

		matchGroups = append(matchGroups, matchGroup)

	}
	return matchGroups, nil
}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Errorf("Error starting Simple chaincode: %s", err)
	}

}

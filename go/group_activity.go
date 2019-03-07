package main

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"math/rand"
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
var allParticipantsInActivity int = 0             // 所有参与人数
var idToSequenceNumber = make(map[string]int)     // 用户ID对应到撮合时的序号
var ActivityTimeClass = make(map[string][]string) // 按照注册时间分类后的结果

//测试git
// --------------------------------------------
// 蚁群算法的全局常量
// --------------------------------------------
// 活动参与人数上下界
const M_LOW int = 11
const M_HIGH int = 11

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

type Request struct {
	ID           string
	Location     string //位置
	RegisterTime int64  //客户端选择某一个上午，计算出那个上午的开始时间，再发给链码，这样方便以后修改时间选择策略
	ActivityTime string
	Deposit      string //押金
	State        string //被撮合状态，0未撮合，1已撮合还未到参加活动时间，2取消戳和，3已撮合被判断未参加活动，4已撮合被判断已参加活动
	Owner        string
	ResultID     string //被撮合到同一组别的用户会被分配一个相同的uid
}

// 在state db中使用GenerateTime作为result的主键
type Result struct {
	GenerateTime string //撮合产生时间
	IDs          []string
	Owners       []string
	ActivityTime string
	CompleteTime string //撮合完成时间
}

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

// 初始化账本方法
func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) sc.Response {
	test()
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
	case "showAllRequest":
		return s.showAllRequest(stub)

	case "initLedger":
		return s.initLedger(stub)
	case "aca":
		return s.aca(stub, args)
	default:
		return shim.Error("no such method")
	}
	//if function == "createRequest" {
	//	return s.createRequest(stub, args)
	//} else if function == "updateRequest" {
	//	return s.updateRequest(stub, args)
	//} else if function == "cancelRequest" {
	//	return s.cancelRequest(stub, args)
	//} else if function == "confirmOrder" { //确认大家是否参加活动，暂时没有这个函数
	//	return s.confirmOrder(stub, args)
	//} else if function == "initLedger" {
	//	return s.initLedger(stub)
	//} else if function == "aca" {
	//	return s.aca(stub, args)
	//}
	////如果用户意图不符合如上，进行错误提示
	//return shim.Error("error function")
}

//接口方法：发布请求
func (s *SmartContract) createRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	existRequest, _ := stub.GetState(args[0])
	if existRequest != nil {
		return shim.Error("Exist Request. ")
	}

	var owner, _ = GetCertAttribute2(stub)
	var request = Request{ID: args[0], Location: args[1], RegisterTime: time.Now().Unix(), ActivityTime: args[2], Deposit: args[3], State: "0", Owner: owner, ResultID: ""}
	requestAsBytes, _ := json.Marshal(request)

	stub.PutState(args[0], requestAsBytes)

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("  register success")

	// 参与
	applicantNum = applicantNum + 1
	return shim.Success(payload.Bytes())
}

//接口方法：修改请求
func (s *SmartContract) updateRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	//判断该请求是否是该用户所有
	//获取该请求真实所有者
	requestID := args[0]
	requestAsBytes, err := stub.GetState(requestID)
	if err != nil {
		return shim.Error("Failed to get request:" + err.Error())
	}
	if requestAsBytes == nil {
		return shim.Error("request does not exist")
	}
	request := Request{}
	err = json.Unmarshal(requestAsBytes, &request) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	owner := request.Owner
	//获取该交易用户名
	var user, _ = GetCertAttribute2(stub)
	if owner != user {
		return shim.Error("Error User")
	}
	request = Request{ID: args[0], Location: args[1], ActivityTime: args[2], Deposit: args[3], Owner: owner, State: "0", ResultID: ""}
	requestAsBytes, _ = json.Marshal(request)
	stub.PutState(args[0], requestAsBytes)

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("  update success")

	return shim.Success(payload.Bytes())
}

//接口方法：取消请求
func (s *SmartContract) cancelRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
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
	var user, _ = GetCertAttribute2(stub)
	if owner != user {
		return shim.Error("您不拥有该请求")
	}
	//判断是否已被撮合
	if request.State != "0" {
		return shim.Error("请求已被撮合，无法取消")
	}

	stub.DelState(args[0])

	var payload bytes.Buffer
	payload.WriteString("ID:")
	payload.WriteString(args[0])
	payload.WriteString("cancel success")

	return shim.Success(payload.Bytes())
}

func (s *SmartContract) showAllRequest(stub shim.ChaincodeStubInterface) sc.Response {

	requestAsBytes, err := queryRequestValueByKeyWithRegex(stub, []string{"ActivityTime", ""})
	if err != nil {
		return shim.Error("Failed to get request:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("request does not exist")
	}
	return shim.Success(requestAsBytes)
}

// 活动结束后用于确认请求
func (s *SmartContract) confirmOrder(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
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

func (s *SmartContract) getAllLocationsToDapp(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// 报名的所有地点
func getAllLocations(stub shim.ChaincodeStubInterface) ([]byte, error) {

	//resultsIterator, err := stub.GetQueryResult(queryString)
	//if err != nil {
	//	return nil, err
	//}
	//defer resultsIterator.Close()
	//
	//buffer, err := constructQueryResponseFromIterator(resultsIterator)
	//if err != nil {
	//	return nil, err
	//}
	//
	//fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	//
	//return buffer.Bytes(), nil
	//requestAsBytes, err := queryRequestValueByKeyWithRegex(stub, []string{"Location", ""})
	//if err != nil {
	//	return nil, errors.New("Failed to get request:" + err.Error())
	//} else if requestAsBytes == nil {
	//	return nil, errors.New("request does not exist")
	//}
	//request := Request{}
	//err = json.Unmarshal(requestAsBytes, &request) //unmarshal it aka JSON.parse()
	//if err != nil {
	//	return nil, errors.New("Failed to Unmarshal:" + err.Error())
	//}
	return nil, nil
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
// }

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
	}
	uname := cert.Subject.CommonName
	fmt.Println("Name:" + uname)
	return uname, nil
}

//工具方法：查询state db中的key和value
func queryRequestValueByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}}", args[0], args[1])
	requestAsBytes, err := getQueryResultForQueryString(stub, queryString)

	return requestAsBytes, err
}

//工具方法：查询state db中的key和values，使用正则表达式
func queryRequestValueByKeyWithRegex(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":{\"$regex\":\"%s\"}}}", args[0], args[1])
	requestAsBytes, err := getQueryResultForQueryString(stub, queryString)

	return requestAsBytes, err
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
//工具方法：过滤查询
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

//工具方法：撮合算法
//工具方法：新建撮合

//根据时间、空间等约束，将满足条件的报名者筛选出来（应是按时间、地点划分成的多组，先只考虑一组报名者）
func test() {
	//模拟初始数据
	deposit = append(deposit, 10, 12, 3, 4, 8, 16, 7, 5, 9, 10, 11, 15, 13, 14, 15, 13, 15, 16, 16, 18)
	timeValue = append(timeValue, 1, 2, 3, 3, 3, 2, 0, 1, 2, 2, 0, 3, 1, 1, 2, 4, 1, 0, 2, 3)
	applicantNum = len(deposit)
	for i := 0; i < applicantNum; i++ {
		depositAndTime = append(depositAndTime, deposit[i]*0.5+timeValue[i]*0.5)
	}

	fmt.Println(depositAndTime)
}

func (s *SmartContract) filter(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	err := getAllActivityTime(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	requestAsBytes, err := stub.GetStateByRange("1", "999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer requestAsBytes.Close()
	for k, v := range ActivityTimeClass {
		err = initUserRegisterInfo(requestAsBytes, k, v)
	}

	var payload bytes.Buffer
	payload.WriteString("Match Make Success")
	return shim.Success(payload.Bytes())
}

// 根据活动时间对报名用户进行分类
func getAllActivityTime(stub shim.ChaincodeStubInterface) error {
	requestAsBytes, err := stub.GetStateByRange("1", "999")

	if err != nil {
		return err
	}
	defer requestAsBytes.Close()

	for requestAsBytes.HasNext() {
		singleRequest, err := requestAsBytes.Next()
		if err != nil {
			return err
		}
		request := Request{}
		err = json.Unmarshal(singleRequest.Value, &request)
		if err != nil {
			return err
		}
		IDs := []string{}
		IDs = ActivityTimeClass[request.ActivityTime]
		IDs = append(IDs, request.ID)

		ActivityTimeClass[request.ActivityTime] = IDs
	}

	return nil
}

// 进行一次撮合，初始化撮合需要使用的参数
func initUserRegisterInfo(requestAsBytes shim.StateQueryIteratorInterface, activityTime string, IDs []string) error {
	serialNum := 0
	for requestAsBytes.HasNext() {
		singleRequest, err := requestAsBytes.Next()
		if err != nil {
			return err
		}

		idToSequenceNumber[singleRequest.Key] = serialNum
		serialNum++

		request := Request{}
		err = json.Unmarshal(singleRequest.Value, &request)
		if err != nil {
			return err
		}

		strDeposit, err := strconv.ParseFloat(request.Deposit, 64)
		if err != nil {
			return err
		}

		// 保证金
		deposit = append(deposit, strDeposit)

		// 时间价值
		tV := time.Now().Unix() - request.RegisterTime
		timeValue = append(timeValue, float64(tV))
	}

	return nil
}

//蚁群算法
func (s *SmartContract) aca(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if deposit == nil {
		return shim.Error("No Enough participants")
	}

	var payload bytes.Buffer

	// 迭代搜索
	var result_bestSatis float64
	var result_bestPath [][]int
	result_bestSatis, result_bestPath = acaSearch()
	fmt.Println("最大满意度为：", result_bestSatis)
	fmt.Println("对应的最佳分配结果为：", result_bestPath)

	stub.PutState("max satisfaction", Float64ToByte(result_bestSatis))

	payload.WriteString("max satisfaction:")
	payload.WriteString(strconv.FormatFloat(result_bestSatis, 'f', -1, 64))

	return shim.Success(payload.Bytes())
}

func initPheromoneMatrix() {
	for i := 0; i < applicantNum; i++ {
		temp := make([]float64, 0, applicantNum)
		for j := 0; j < applicantNum; j++ {
			temp = append(temp, 1)
		}
		pheromoneMatrix = append(pheromoneMatrix, temp)
	}
}

func initMatrix(m, n, value int) [][]int {
	var result [][]int
	for i := 0; i < m; i++ {
		temp := make([]int, 0, n)
		for j := 0; j < n; j++ {
			temp = append(temp, value)
		}
		result = append(result, temp)
	}
	return result
}

func acaSearch() (float64, [][]int) {

	var m, i int

	//记录每个m对应的最佳路径和最大满意度
	var resultPath [][][]int
	var resultSatis []float64

	//增加一重循环，确定活动参与人数m的取值
	for m = M_LOW; m <= applicantNum && m <= M_HIGH; m++ {

		fmt.Println("m=", m)

		// 初始化信息素矩阵，设初始元素值全为1
		initPheromoneMatrix()

		//初始化criticalPointMatrix
		for i = 0; i < applicantNum; i++ {
			criticalPointMatrix = append(criticalPointMatrix, -1)
		}

		// 初始化sortedPheromoneMatrix矩阵
		for i := 0; i < applicantNum; i++ {
			var temp []int
			for j := 0; j < applicantNum; j++ {
				temp = append(temp, j)
			}
			sortedPheromoneMatrix = append(sortedPheromoneMatrix, temp)
		}

		//记录当前m值下的最佳分配路径和最大满意度
		var bestPath [][]int = initMatrix(applicantNum, applicantNum, 0)
		var bestSatis float64 = MIN_VALUE

		var nodes []int
		//根据m的具体取值，初始化nodes数组
		for i = 0; i < m; i++ {
			nodes = append(nodes, 1)
		}
		for i = m; i < applicantNum; i++ {
			nodes = append(nodes, 0)
		}

		// 当前m下所有迭代中每个蚂蚁分配结果的满意度
		var resultData [][]float64

		for itCount := 0; itCount < iteratorNum; itCount++ {

			// 本次迭代中，所有蚂蚁的路径
			var pathMatrix_allAnt [][][]int

			for antCount := 0; antCount < antNum; antCount++ {
				// 第antCount只蚂蚁的分配策略(pathMatrix[i][j]表示第antCount只蚂蚁将节点i分配给报名者j)
				var pathMatrix_oneAnt [][]int = initMatrix(applicantNum, applicantNum, 0) //初始化数组元素全为0
				var assignedApplicant []int
				for nodeCount := 0; nodeCount < applicantNum; nodeCount++ {
					// 将第nodeCount个节点分配给第applicantCount个报名者
					applicantCount := assignOneNode(assignedApplicant, antCount, nodeCount)
					pathMatrix_oneAnt[nodeCount][applicantCount] = 1
					assignedApplicant = append(assignedApplicant, applicantCount)
				}
				// 将当前蚂蚁的路径加入pathMatrix_allAnt
				pathMatrix_allAnt = append(pathMatrix_allAnt, pathMatrix_oneAnt)
			}

			// 计算 本次迭代中 所有蚂蚁 的任务分配的整体满意度
			var satisArray_oneIt []float64 = calSatis_oneIt(pathMatrix_allAnt, nodes)
			// 将本地迭代中 所有蚂蚁的 节点分配满意度加入总结果集
			resultData = append(resultData, satisArray_oneIt)

			// 更新信息素
			bestAntIndex := updatePheromoneMatrix(pathMatrix_allAnt, satisArray_oneIt)

			fmt.Println(satisArray_oneIt[bestAntIndex])

			// 更新当前m下的最大满意度和最佳路径
			if satisArray_oneIt[bestAntIndex] > bestSatis {
				bestSatis = satisArray_oneIt[bestAntIndex]
				bestPath = pathMatrix_allAnt[bestAntIndex]
			}

		}
		resultSatis = append(resultSatis, bestSatis)
		resultPath = append(resultPath, bestPath)
	}

	// 通过计算满意度选择最佳的分配方案
	var result_bestSatis float64 = resultSatis[0]
	var result_bestIndex int = 0
	for i = 1; i < len(resultSatis); i++ {
		if resultSatis[i] > result_bestSatis {
			result_bestSatis = resultSatis[i]
			result_bestIndex = i
		}
	}
	var result_bestPath [][]int = resultPath[result_bestIndex]

	return result_bestSatis, result_bestPath
}

func assignOneNode(assignedApplicant []int, antCount int, nodeCount int) int {

	// 去除已分配的报名者下标
	//fmt.Println(sortedPheromoneMatrix[nodeCount])
	var sorted_index []int = make([]int, applicantNum)
	//sorted_index = sortedPheromoneMatrix[nodeCount]
	copy(sorted_index, sortedPheromoneMatrix[nodeCount])
	for i := 0; i < len(assignedApplicant); i++ {
		for j := 0; j < len(sorted_index); j++ {
			if sorted_index[j] == assignedApplicant[i] {
				// 去掉 sorted_index[j]
				sorted_index = append(sorted_index[:j], sorted_index[j+1:]...)
				break
			}
		}
	}

	// 若当前蚂蚁编号在临界点之前，则采用最大信息素的分配方式 （且此时该下标对应的报名者未分配）
	if antCount <= criticalPointMatrix[nodeCount] {
		//
		var sameFirst int = 0
		for i := 0; i < len(sorted_index)-1; i++ {
			if pheromoneMatrix[nodeCount][sorted_index[i]] != pheromoneMatrix[nodeCount][sorted_index[i+1]] {
				sameFirst = i
				break
			}
		}
		if sameFirst == 0 {
			return sorted_index[0]
		}
		time.Sleep(3 * time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		result := rand.Intn(sameFirst)
		return sorted_index[result]
		//return maxPheromoneMatrix[nodeCount]
	}

	// 若当前蚂蚁编号在临界点之后，则采用随机分配方式
	// 设置随机数种子
	time.Sleep(3 * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(sorted_index))
	return sorted_index[index]
}

func calSatis_oneIt(pathMatrix_allAnt [][][]int, nodes []int) []float64 {
	var satisArray_oneIt []float64
	// 计算每个蚂蚁分配结果的满意度--方差倒数
	var data []float64
	var sum float64 = 0
	for i := 0; i < antNum; i++ {
		for nodeIndex := 0; nodeIndex < applicantNum; nodeIndex++ {
			for applicantIndex := 0; applicantIndex < applicantNum; applicantIndex++ {
				if pathMatrix_allAnt[i][nodeIndex][applicantIndex] == 1 && nodes[nodeIndex] == 1 {
					data = append(data, depositAndTime[applicantIndex])
					sum += depositAndTime[applicantIndex]
				}
			}
		}
		// 计算方差
		var ave float64 = sum / float64(len(data))
		sum = 0
		for i := 0; i < len(data); i++ {
			sum += (data[i] - ave) * (data[i] - ave)
		}
		sum = 1.0 / (sum / float64(len(data)))
		//sum = -1.0 * (sum / float64(len(data)))
		satisArray_oneIt = append(satisArray_oneIt, sum)
		data = nil
		sum = 0
	}
	return satisArray_oneIt
}

func updatePheromoneMatrix(pathMatrix_allAnt [][][]int, satisArray_oneIt []float64) int {

	var bestAntIndex int

	// 所有信息素均衰减p%
	for i := 0; i < applicantNum; i++ {
		for j := 0; j < applicantNum; j++ {
			pheromoneMatrix[i][j] *= p
		}
	}

	// 找出满意度最大的蚂蚁编号
	var maxSatis = MIN_VALUE
	var maxIndex = -1
	for antIndex := 0; antIndex < antNum; antIndex++ {
		if satisArray_oneIt[antIndex] > maxSatis {
			maxSatis = satisArray_oneIt[antIndex]
			maxIndex = antIndex
		}
	}
	bestAntIndex = maxIndex

	// 将本次迭代中最优路径的信息素增加q%
	for nodeIndex := 0; nodeIndex < applicantNum; nodeIndex++ {
		for applicantIndex := 0; applicantIndex < applicantNum; applicantIndex++ {
			if pathMatrix_allAnt[maxIndex][nodeIndex][applicantIndex] == 1 {
				pheromoneMatrix[nodeIndex][applicantIndex] *= q
				break
			}
		}
	}

	//清空
	maxPheromoneMatrix = nil
	criticalPointMatrix = nil
	sortedPheromoneMatrix = nil
	for nodeIndex := 0; nodeIndex < applicantNum; nodeIndex++ {
		var maxPheromone float64 = pheromoneMatrix[nodeIndex][0]
		var maxIndex = 0
		var sumPheromone float64 = pheromoneMatrix[nodeIndex][0]
		var isAllSame = true

		for applicantIndex := 1; applicantIndex < applicantNum; applicantIndex++ {
			if pheromoneMatrix[nodeIndex][applicantIndex] > maxPheromone {
				maxPheromone = pheromoneMatrix[nodeIndex][applicantIndex]
				maxIndex = nodeIndex
			}

			if pheromoneMatrix[nodeIndex][applicantIndex] != pheromoneMatrix[nodeIndex][applicantIndex-1] {
				isAllSame = false
			}

			sumPheromone += pheromoneMatrix[nodeIndex][applicantIndex]
		}

		// 若本行信息素全都相等，则随机选择一个作为最大信息素
		if isAllSame == true {
			//设置随机数种子
			time.Sleep(3 * time.Millisecond)
			rand.Seed(time.Now().UnixNano())
			maxIndex = rand.Intn(applicantNum)
			maxPheromone = pheromoneMatrix[nodeIndex][maxIndex]
		}

		// 将本行最大信息素的下标加入maxPheromoneMatrix
		maxPheromoneMatrix = append(maxPheromoneMatrix, maxIndex)

		// 记录本行信息素由大到小排序的下标
		var oneRow []float64 = make([]float64, applicantNum)
		copy(oneRow, pheromoneMatrix[nodeIndex])
		sortedPheromoneMatrix_one := sortPheromoneMatrix(oneRow)
		sortedPheromoneMatrix = append(sortedPheromoneMatrix, sortedPheromoneMatrix_one)

		// 将本次迭代的蚂蚁临界编号加入criticalPointMatrix(该临界点之前的蚂蚁的任务分配根据最大信息素原则，而该临界点之后的蚂蚁采用随机分配策略)
		criticalPointMatrix = append(criticalPointMatrix, int(math.Floor(float64(antNum)*(maxPheromone/sumPheromone)+0.5)))
	}
	fmt.Println(criticalPointMatrix)
	return bestAntIndex
}

func sortPheromoneMatrix(oneRow []float64) []int {
	var sorted_index []int
	for k := 0; k < applicantNum; k++ {
		sorted_index = append(sorted_index, k)
	}
	for i := 0; i < len(oneRow)-1; i++ {
		for j := 0; j < len(oneRow)-1-i; j++ {
			if oneRow[j] <= oneRow[j+1] {
				var temp float64 = oneRow[j]
				oneRow[j] = oneRow[j+1]
				oneRow[j+1] = temp

				var index int = sorted_index[j]
				sorted_index[j] = sorted_index[j+1]
				sorted_index[j+1] = index
			}
		}
	}

	return sorted_index
}

// func main()  {
// 	filter()
// 	aca()
// 	//fmt.Println("Hello, World!")
// }

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Errorf("Error starting Simple chaincode: %s", err)
	}

}

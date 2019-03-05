package main

import "fmt"
import "math/rand"
import "math"
import (
	"time"
)

//全局常量
//活动参与人数上下界
const M_LOW int = 11
const M_HIGH int = 11

//有符号整型最小值
const MIN_VALUE = float64(^int(^uint(0) >> 1))

// 产生随机数前的休眠时间
const DURATION int = 3

//全局变量
var applicantNum int    //报名人数
var deposit []float64   //每位报名者对应的押金
var timeValue []float64 //每位报名者当前对应的时间价值
var depositAndTime []float64

var iteratorNum = 1000 //迭代次数
var antNum = 30        //蚂蚁数量

var pheromoneMatrix [][]float64 //信息素矩阵
var maxPheromoneMatrix []int    //pheromoneMatrix矩阵的每一行中最大信息素的下标
var sortedPheromoneMatrix [][]int  //pheromoneMatrix矩阵的每一行中信息素从大到小排列的下标
var criticalPointMatrix []int   //在一次迭代中，采用随机分配策略的蚂蚁的临界编号

var p float64 = 0.9 //每完成一次迭代后，信息素衰减的比例
var q float64 = 1.1 //蚂蚁每次经过一条路径，信息素增加的比例

//根据时间、空间等约束，将满足条件的报名者筛选出来（应是按时间、地点划分成的多组，先只考虑一组报名者）
func filter() {
	//模拟初始数据
	deposit = append(deposit, 10, 12, 3, 4, 8, 16, 7, 5, 9, 10, 11, 15, 13, 14, 15, 13, 15, 16, 16, 18)
	timeValue = append(timeValue, 1, 2, 3, 3, 3, 2, 0, 1, 2, 2, 0, 3, 1, 1, 2, 4, 1, 0, 2, 3)
	applicantNum = len(deposit)
	for i := 0; i < applicantNum; i++ {
		depositAndTime = append(depositAndTime, deposit[i]*0.5+timeValue[i]*0.5)
	}
}

//蚁群算法
func aca() {

	// 迭代搜索
	var result_bestSatis float64
	var result_bestPath [][]int
	result_bestSatis, result_bestPath= acaSearch()
	fmt.Println("最大满意度为：", result_bestSatis)
	fmt.Println("对应的最佳分配结果为：", result_bestPath)
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
	for i:=0;i<len(assignedApplicant);i++ {
		for j:=0;j<len(sorted_index);j++ {
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
		for i:=0;i<len(sorted_index)-1;i++ {
			if pheromoneMatrix[nodeCount][sorted_index[i]] !=  pheromoneMatrix[nodeCount][sorted_index[i+1]] {
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
	for k:=0;k<applicantNum;k++ {
		sorted_index = append(sorted_index, k)
	}
	for i:=0;i<len(oneRow)-1;i++ {
		for j:=0;j<len(oneRow)-1-i;j++ {
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

func main()  {
	filter()
	aca()
}

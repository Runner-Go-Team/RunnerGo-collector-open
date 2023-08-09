package server

import (
	"encoding/json"
	"fmt"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/dal/kao"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/dal/redis"
	log2 "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/Shopify/sarama"
	"github.com/shopspring/decimal"
	"sort"
	"sync"
)

func Execute(host string) {
	defer pkg.CapRecover()
	var partitionMap = new(sync.Map)
	partitionList := redis.QueryStressBelongPartition(pkg.LocalIp)
	partitionList = redis.QueryTotalKafkaPartition(partitionList)
	if len(partitionList) < conf.Conf.Kafka.Num {
		return
	}
	by, err := json.Marshal(partitionList)
	if err != nil {
		log2.Logger.Error("分区列表-json转换失败")
		return
	}
	redis.InsertStressBelongPartition(conf.Conf.Kafka.StressBelongPartition, string(by))
	topic := conf.Conf.Kafka.Topic
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	for {
		//// 睡眠一分钟，再循环执行
		time.Sleep(60 * time.Second)

		for _, value := range partitionList {
			if _, ok := partitionMap.Load(value); ok {
				continue
			}
			consumer, err := sarama.NewConsumer([]string{host}, sarama.NewConfig())
			if err != nil {
				log2.Logger.Error("topic  :"+topic+", 创建消费者失败: ", err, "   host:  ", host)
				partitionMap.Delete(value)
				continue
			}
			partitionMap.Store(value, true)
			pc, err := consumer.ConsumePartition(topic, value, sarama.OffsetNewest)
			if err != nil {
				log2.Logger.Error(fmt.Sprintf("消费分区错误：分区%d, 错误：%s：    ", value, err.Error()))
				partitionMap.Delete(value)
				break
			}
			pc.IsPaused()
			go ReceiveMessage(pc, partitionMap, value)
		}

	}

}

func ReceiveMessage(pc sarama.PartitionConsumer, partitionMap *sync.Map, partition int32) {
	defer pkg.CapRecover()
	defer pc.AsyncClose()
	defer partitionMap.Delete(partition)

	if pc == nil || partitionMap == nil {
		return
	}
	var requestTimeListMap = make(map[string]kao.RequestTimeList)
	var resultDataMsg = kao.ResultDataMsg{}
	var sceneTestResultDataMsg = new(kao.SceneTestResultDataMsg)
	var machineNum, startTime, runTime, index = int64(0), int64(0), int64(0), 0
	var eventMap = make(map[string]int64)
	var machineMap = make(map[string]map[string]int64)
	log2.Logger.Info("分区：", partition, "   ,开始消费消息")
	for msg := range pc.Messages() {

		err := json.Unmarshal(msg.Value, &resultDataMsg)
		if err != nil {
			log2.Logger.Error("kafka消息转换失败：", err)
			continue
		}
		if resultDataMsg.TeamId == "" {
			log2.Logger.Error("teamId 不能为空")
			continue
		}
		if resultDataMsg.PlanId == "" {
			log2.Logger.Error("planId 不能为空")
			continue
		}
		if resultDataMsg.ReportId == "" {
			log2.Logger.Error("reportId 不能为空")
			continue
		}
		index++
		if machineNum == 0 && resultDataMsg.MachineNum != 0 {
			machineNum = resultDataMsg.MachineNum + 1
		}
		if runTime == 0 {
			runTime = resultDataMsg.Timestamp / 1000
		}

		if startTime == 0 {
			startTime = resultDataMsg.Timestamp
		}

		if resultDataMsg.Start {
			if sceneTestResultDataMsg.PlanId == "" {
				sceneTestResultDataMsg.PlanId = resultDataMsg.PlanId
			}
			if sceneTestResultDataMsg.ReportId == "" {
				sceneTestResultDataMsg.ReportId = resultDataMsg.ReportId
			}
			if sceneTestResultDataMsg.PlanId == "" {
				sceneTestResultDataMsg.PlanId = resultDataMsg.PlanId
			}
			continue
		}

		if resultDataMsg.End {
			machineNum = machineNum - 1
			if machineNum == 1 {
				sceneTestResultDataMsg.End = true
				for eventId, requestTimeList := range requestTimeListMap {
					sort.Sort(requestTimeList)
					if sceneTestResultDataMsg.Results[eventId].TotalRequestNum != 0 {
						sceneTestResultDataMsg.Results[eventId].AvgRequestTime = float64(sceneTestResultDataMsg.Results[eventId].TotalRequestTime) / float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum)
					}
					sceneTestResultDataMsg.Results[eventId].MaxRequestTime = float64(requestTimeList[len(requestTimeList)-1])
					sceneTestResultDataMsg.Results[eventId].MinRequestTime = float64(requestTimeList[0])
					sceneTestResultDataMsg.Results[eventId].FiftyRequestTimeline = 50
					sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLine = 90
					sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLine = 95
					sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLine = 99
					if sceneTestResultDataMsg.Results[eventId].PercentAge > 0 && sceneTestResultDataMsg.Results[eventId].PercentAge != 101 &&
						sceneTestResultDataMsg.Results[eventId].PercentAge != 50 && sceneTestResultDataMsg.Results[eventId].PercentAge != 90 &&
						sceneTestResultDataMsg.Results[eventId].PercentAge != 95 && sceneTestResultDataMsg.Results[eventId].PercentAge != 99 &&
						sceneTestResultDataMsg.Results[eventId].PercentAge != 100 {

						sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].PercentAge, requestTimeList)
					}
					sceneTestResultDataMsg.Results[eventId].FiftyRequestTimelineValue = float64(requestTimeList[len(requestTimeList)/2])
					sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLineValue = kao.TimeLineCalculate(90, requestTimeList)
					sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLineValue = kao.TimeLineCalculate(95, requestTimeList)
					sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLineValue = kao.TimeLineCalculate(99, requestTimeList)
					if sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine != 0 {
						sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine, requestTimeList)
					}

					tpsTime := float64(sceneTestResultDataMsg.Results[eventId].EndTime-sceneTestResultDataMsg.Results[eventId].StartTime) / 1000
					if tpsTime != 0 {
						sceneTestResultDataMsg.Results[eventId].Tps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum) / tpsTime).Round(2).Float64()
						sceneTestResultDataMsg.Results[eventId].STps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].SuccessNum) / tpsTime).Round(2).Float64()
					}
					//if sceneTestResultDataMsg.Results[eventId].TotalRequestTime != 0 {
					//	concurrent := sceneTestResultDataMsg.Results[eventId].Concurrency
					//	rpsTime := float64(time.Second) * float64(concurrent) / float64(sceneTestResultDataMsg.Results[eventId].TotalRequestTime)
					//	sceneTestResultDataMsg.Results[eventId].Rps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum) * rpsTime).Round(2).Float64()
					//	sceneTestResultDataMsg.Results[eventId].SRps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].SuccessNum) * rpsTime).Round(2).Float64()
					//}
					rpsTime := float64(sceneTestResultDataMsg.Results[eventId].StageEndTime-sceneTestResultDataMsg.Results[eventId].StageStartTime) / 1000
					if rpsTime != 0 {
						sceneTestResultDataMsg.Results[eventId].Rps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].StageTotalRequestNum) / rpsTime).Round(2).Float64()
						sceneTestResultDataMsg.Results[eventId].SRps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].StageSuccessNum) / rpsTime).Round(2).Float64()
					}

					sceneTestResultDataMsg.Results[eventId].StageEndTime = 0
					sceneTestResultDataMsg.Results[eventId].StageSuccessNum = 0
					sceneTestResultDataMsg.Results[eventId].StageStartTime = 0
					sceneTestResultDataMsg.Results[eventId].StageTotalRequestNum = 0

				}
				sceneTestResultDataMsg.TimeStamp = resultDataMsg.Timestamp / 1000
				if err = redis.InsertTestData(machineMap, sceneTestResultDataMsg, runTime); err != nil {
					log2.Logger.Error("redis写入数据失败:", err)
				}
				if err = redis.UpdatePartitionStatus(conf.Conf.Kafka.Key, partition); err != nil {
					log2.Logger.Error("修改kafka分区状态失败： ", err)
				}
				log2.Logger.Info(fmt.Sprintf("删初分区的key: %s, 分区的值：%d 成功! 本次共消费：%d 条消息", conf.Conf.Kafka.Key, partition, index))
				return
			}
			continue
		}
		if sceneTestResultDataMsg.SceneId == "" {
			sceneTestResultDataMsg.SceneId = resultDataMsg.SceneId
		}
		if sceneTestResultDataMsg.SceneName == "" {
			sceneTestResultDataMsg.SceneName = resultDataMsg.SceneName
		}
		if sceneTestResultDataMsg.ReportId == "" {
			sceneTestResultDataMsg.ReportId = resultDataMsg.ReportId
		}
		if sceneTestResultDataMsg.ReportName == "" {
			sceneTestResultDataMsg.ReportName = resultDataMsg.ReportName
		}
		if sceneTestResultDataMsg.TeamId == "" {
			sceneTestResultDataMsg.TeamId = resultDataMsg.TeamId
		}
		if sceneTestResultDataMsg.PlanId == "" {
			sceneTestResultDataMsg.PlanId = resultDataMsg.PlanId
		}
		if sceneTestResultDataMsg.PlanName == "" {
			sceneTestResultDataMsg.PlanName = resultDataMsg.PlanName
		}
		if sceneTestResultDataMsg.Results == nil {
			sceneTestResultDataMsg.Results = make(map[string]*kao.ApiTestResultDataMsg)
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId] == nil {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId] = new(kao.ApiTestResultDataMsg)
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].EventId == "" {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].EventId = resultDataMsg.EventId
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].Name == "" {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].Name = resultDataMsg.Name
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneId == "" {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneId = resultDataMsg.SceneId
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanId == "" {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanId = resultDataMsg.PlanId
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanName == "" {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanName = resultDataMsg.PlanName
		}
		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneName == "" {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneName = resultDataMsg.SceneName
		}
		if resultDataMsg.PercentAge != 0 && resultDataMsg.PercentAge < 100 {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].CustomRequestTimeLine = resultDataMsg.PercentAge
		} else {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].CustomRequestTimeLine = 0
		}

		if concurrency, ok := machineMap[resultDataMsg.MachineIp][resultDataMsg.EventId]; !ok {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].Concurrency += resultDataMsg.Concurrency
			eventMap[resultDataMsg.EventId] = resultDataMsg.Concurrency
			machineMap[resultDataMsg.MachineIp] = eventMap
		} else {
			if concurrency != resultDataMsg.Concurrency {
				machineMap[resultDataMsg.MachineIp][resultDataMsg.EventId] = resultDataMsg.Concurrency
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].Concurrency += resultDataMsg.Concurrency - concurrency
			}

		}

		sceneTestResultDataMsg.Results[resultDataMsg.EventId].PercentAge = resultDataMsg.PercentAge
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].ErrorThreshold = resultDataMsg.ErrorThreshold
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].ResponseThreshold = resultDataMsg.ResponseThreshold
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].RequestThreshold = resultDataMsg.RequestThreshold
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].ReceivedBytes += resultDataMsg.ReceivedBytes
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].SendBytes += resultDataMsg.SendBytes
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestNum += 1
		sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestTime += resultDataMsg.RequestTime

		sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageTotalRequestNum += 1

		if resultDataMsg.IsSucceed {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].SuccessNum += 1
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageSuccessNum += 1
		} else {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].ErrorNum += 1
		}

		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].StartTime == 0 || sceneTestResultDataMsg.Results[resultDataMsg.EventId].StartTime > resultDataMsg.StartTime {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].StartTime = resultDataMsg.StartTime
		}

		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].EndTime == 0 || sceneTestResultDataMsg.Results[resultDataMsg.EventId].EndTime < resultDataMsg.EndTime {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].EndTime = resultDataMsg.EndTime
		}

		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageStartTime == 0 || sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageStartTime > resultDataMsg.StartTime {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageStartTime = resultDataMsg.StartTime
		}

		if sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageEndTime == 0 || sceneTestResultDataMsg.Results[resultDataMsg.EventId].EndTime > resultDataMsg.EndTime {
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].StageEndTime = resultDataMsg.EndTime
		}

		requestTimeListMap[resultDataMsg.EventId] = append(requestTimeListMap[resultDataMsg.EventId], resultDataMsg.RequestTime)
		if resultDataMsg.Timestamp-startTime >= 2000 {
			startTime = resultDataMsg.Timestamp
			if sceneTestResultDataMsg.ReportId == "" || sceneTestResultDataMsg.Results == nil {
				break
			}
			for eventId, requestTimeList := range requestTimeListMap {
				sort.Sort(requestTimeList)
				if sceneTestResultDataMsg.Results[eventId].TotalRequestNum != 0 {
					sceneTestResultDataMsg.Results[eventId].AvgRequestTime = float64(sceneTestResultDataMsg.Results[eventId].TotalRequestTime) / float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum)
				}

				sceneTestResultDataMsg.Results[eventId].MaxRequestTime = float64(requestTimeList[len(requestTimeList)-1])
				sceneTestResultDataMsg.Results[eventId].MinRequestTime = float64(requestTimeList[0])
				sceneTestResultDataMsg.Results[eventId].FiftyRequestTimeline = 50
				sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLine = 90
				sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLine = 95
				sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLine = 99
				if sceneTestResultDataMsg.Results[eventId].PercentAge > 0 && sceneTestResultDataMsg.Results[eventId].PercentAge != 101 &&
					sceneTestResultDataMsg.Results[eventId].PercentAge != 50 && sceneTestResultDataMsg.Results[eventId].PercentAge != 90 &&
					sceneTestResultDataMsg.Results[eventId].PercentAge != 95 && sceneTestResultDataMsg.Results[eventId].PercentAge != 99 &&
					sceneTestResultDataMsg.Results[eventId].PercentAge != 100 {

					sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].PercentAge, requestTimeList)
				}
				sceneTestResultDataMsg.Results[eventId].FiftyRequestTimelineValue = float64(requestTimeList[len(requestTimeList)/2])
				sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLineValue = kao.TimeLineCalculate(90, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLineValue = kao.TimeLineCalculate(95, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLineValue = kao.TimeLineCalculate(99, requestTimeList)
				if sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine != 0 {
					sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine, requestTimeList)
				}

				tpsTime := float64(sceneTestResultDataMsg.Results[eventId].EndTime-sceneTestResultDataMsg.Results[eventId].StartTime) / 1000
				if tpsTime != 0 {
					sceneTestResultDataMsg.Results[eventId].Tps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum) / tpsTime).Round(2).Float64()
					sceneTestResultDataMsg.Results[eventId].STps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].SuccessNum) / tpsTime).Round(2).Float64()
				}
				//if sceneTestResultDataMsg.Results[eventId].TotalRequestTime != 0 {
				//	concurrent := sceneTestResultDataMsg.Results[eventId].Concurrency
				//	rpsTime := float64(time.Second) * float64(concurrent) / float64(sceneTestResultDataMsg.Results[eventId].TotalRequestTime)
				//	sceneTestResultDataMsg.Results[eventId].Rps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum) * rpsTime).Round(2).Float64()
				//	sceneTestResultDataMsg.Results[eventId].SRps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].SuccessNum) * rpsTime).Round(2).Float64()
				//}

				rpsTime := float64(sceneTestResultDataMsg.Results[eventId].StageEndTime-sceneTestResultDataMsg.Results[eventId].StageStartTime) / 1000
				if rpsTime != 0 {
					sceneTestResultDataMsg.Results[eventId].Rps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].StageTotalRequestNum) / rpsTime).Round(2).Float64()
					sceneTestResultDataMsg.Results[eventId].SRps, _ = decimal.NewFromFloat(float64(sceneTestResultDataMsg.Results[eventId].StageSuccessNum) / rpsTime).Round(2).Float64()
				}

				sceneTestResultDataMsg.Results[eventId].StageEndTime = 0
				sceneTestResultDataMsg.Results[eventId].StageSuccessNum = 0
				sceneTestResultDataMsg.Results[eventId].StageStartTime = 0
				sceneTestResultDataMsg.Results[eventId].StageTotalRequestNum = 0

				sceneTestResultDataMsg.TimeStamp = startTime / 1000

			}
			if err = redis.InsertTestData(machineMap, sceneTestResultDataMsg, runTime); err != nil {
				log2.Logger.Error("测试数据写入redis失败：     ", err)
				continue
			}
		}

	}
}

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
	"strconv"
	"sync"
	"time"
)

func Execute(host string) {
	defer pkg.CapRecover()

	topic := conf.Conf.Kafka.Topic
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	pubSub := redis.SubscribeMsg(conf.Conf.Kafka.RunKafkaPartition)
	partitionCh := pubSub.Channel()
	for {
		select {
		case c := <-partitionCh:
			log2.Logger.Debug("分区是：  ", c.Payload)
			partition, err := strconv.Atoi(c.Payload)
			if err != nil {
				log2.Logger.Error("分区转换失败：", err.Error())
				continue
			}
			consumer, err := sarama.NewConsumer([]string{host}, sarama.NewConfig())
			if err != nil {
				log2.Logger.Error("topic  :"+topic+", 创建消费者失败: ", err, "   host:  ", host)
				continue
			}
			pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				log2.Logger.Error(fmt.Sprintf("消费分区错误：分区%v, 错误：%s：    ", partition, err.Error()))
				break
			}
			pc.IsPaused()
			go ReceiveMessage(pc, int32(partition))
		}
	}

}

func ReceiveMessage(pc sarama.PartitionConsumer, partition int32) {
	defer pkg.CapRecover()
	defer pc.AsyncClose()

	if pc == nil {
		return
	}
	var requestTimeListMap = make(map[string]kao.RequestTimeList)
	var resultDataMsg = kao.ResultDataMsg{}
	var sceneTestResultDataMsg = new(kao.SceneTestResultDataMsg)
	sceneTestResultDataMsg.Results = make(map[string]*kao.ApiTestResultDataMsg)
	var machineNum, startTime, runTime, index = int64(0), int64(0), int64(0), 0
	var eventMap = make(map[string]int64)
	var machineMap = make(map[string]map[string]int64)

	mu := &sync.Mutex{}
	log2.Logger.Info("分区：", partition, "   ,开始消费消息")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {

		defer wg.Done()
		timer := time.NewTicker(2 * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				startTime = resultDataMsg.Timestamp
				if sceneTestResultDataMsg.ReportId == "" || sceneTestResultDataMsg.Results == nil {
					break
				}
				mu.Lock()
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
				mu.Unlock()
				if err := redis.InsertTestData(machineMap, sceneTestResultDataMsg, runTime); err != nil {
					log2.Logger.Error("redis写入数据失败:", err)
				}
				if sceneTestResultDataMsg.End {
					if err := redis.UpdatePartitionStatus(conf.Conf.Kafka.Key, partition); err != nil {
						log2.Logger.Error("修改kafka分区状态失败： ", err)
					}
					log2.Logger.Info(fmt.Sprintf("删初分区的key: %s, 分区的值：%d 成功! 本次共消费：%d 条消息", conf.Conf.Kafka.Key, partition, index))
					return
				}
			}
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				log2.Logger.Error("发生崩溃：", err)
				sceneTestResultDataMsg.End = true
			}
		}()
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
			if sceneTestResultDataMsg.PlanId == "" {
				sceneTestResultDataMsg.PlanId = resultDataMsg.PlanId
			}
			if sceneTestResultDataMsg.ReportId == "" {
				sceneTestResultDataMsg.ReportId = resultDataMsg.ReportId
			}
			if sceneTestResultDataMsg.TeamId == "" {
				sceneTestResultDataMsg.TeamId = resultDataMsg.TeamId
			}
			if resultDataMsg.Start {

				continue
			}

			if resultDataMsg.End {
				machineNum = machineNum - 1
				if machineNum == 1 {
					sceneTestResultDataMsg.End = true
					return
				}
				continue
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

		}

	}()
	wg.Wait()
	log2.Logger.Debug("报告处理完成：", sceneTestResultDataMsg.ReportId)

}

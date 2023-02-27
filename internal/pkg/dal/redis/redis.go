package redis

import (
	"RunnerGo-collector/internal/pkg"
	"RunnerGo-collector/internal/pkg/conf"
	"RunnerGo-collector/internal/pkg/dal/kao"
	"RunnerGo-collector/internal/pkg/log"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var (
	ReportRdb    *redis.Client
	RDB          *redis.Client
	timeDuration = 3 * time.Second
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient(reportAddr, reportPassword string, reportDb int64, addr, password string, db int64) (err error) {
	ReportRdb = redis.NewClient(
		&redis.Options{
			Addr:     reportAddr,
			Password: reportPassword,
			DB:       int(reportDb),
		})
	_, err = ReportRdb.Ping().Result()
	if err != nil {
		return err
	}

	RDB = redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       int(db),
		})
	_, err = RDB.Ping().Result()
	return err
}

func UpdatePartitionStatus(key string, partition int32) (err error) {
	field := fmt.Sprintf("%d", partition)
	err = RDB.HDel(key, field).Err()
	return
}

func deleteStopKey(key string) {
	err := RDB.Del(key).Err()
	if err != nil {
		log.Logger.Error("删除停止任务key失败：", key)
	} else {
		log.Logger.Info(fmt.Sprintf("key:  %s, 删除成功", key))
	}
}

func InsertTestData(machineMap map[string]map[string]int64, sceneTestResultDataMsg *kao.SceneTestResultDataMsg, runTime int64) (err error) {
	data := sceneTestResultDataMsg.ToJson()
	key := fmt.Sprintf("reportData:%s", sceneTestResultDataMsg.ReportId)
	if sceneTestResultDataMsg.End {
		if err != nil {
			log.Logger.Error("报告Id转数字失败：  ", err)
		}
		duration := sceneTestResultDataMsg.TimeStamp - runTime
		pkg.SendStopStressReport(machineMap, sceneTestResultDataMsg.TeamId, sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId, duration)
		//stopPlanKey := fmt.Sprintf("StopPlan:%d:%d:%d", sceneTestResultDataMsg.TeamId, sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId)
		////deleteStopKey(stopPlanKey)
		//adjustKey := fmt.Sprintf("adjust:%d:%d:%d", sceneTestResultDataMsg.TeamId, sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId)
		//deleteStopKey(adjustKey)
	}

	err = ReportRdb.LPush(key, data).Err()
	if err != nil {
		return
	}
	return
}

func Insert(rdb *redis.Client, a string) (err error) {
	err = rdb.LPush("report1", a).Err()
	if err != nil {
		return
	}
	return
}

type A struct {
	B int `json:"a"`
}

func QueryStressBelongPartition(localIp string) (partitionList []int32) {
	res, err := RDB.HGet(conf.Conf.Kafka.StressBelongPartition, localIp).Result()
	if err != nil {
		return
	}
	_ = json.Unmarshal([]byte(res), &partitionList)
	return
}

func QueryTotalKafkaPartition(partitionList []int32) []int32 {
	if partitionList == nil {
		partitionList = []int32{}
	}
	for len(partitionList) < conf.Conf.Kafka.Num {
		res, err := RDB.RPop(conf.Conf.Kafka.TotalKafkaPartition).Result()
		if err != nil {
			continue
		}
		partition, err := strconv.Atoi(res)
		if err != nil {
			continue
		}
		var target bool
		for _, v := range partitionList {
			if v == int32(partition) {
				target = true
				break
			}
		}
		if target {
			continue
		}
		partitionList = append(partitionList, int32(partition))

	}
	return partitionList
}

func InsertStressBelongPartition(key, value string) {
	err := RDB.HSet(key, pkg.LocalIp, value).Err()
	if err != nil {
		return
	}

}

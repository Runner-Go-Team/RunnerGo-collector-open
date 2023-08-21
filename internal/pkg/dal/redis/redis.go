package redis

import (
	"encoding/json"
	"fmt"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/dal/kao"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

var (
	RDB          *redis.ClusterClient
	timeDuration = 3 * time.Second
)

type RedisClient struct {
	Client *redis.ClusterClient
}

func InitRedisClient(addr, password string) (err error) {

	RDB = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:    strings.Split(addr, ";"),
			Password: password,
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

	err = RDB.LPush(key, data).Err()
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

func ExitStressBelongPartition(stressBelongPartition, heartKey string) {
	keys := RDB.HKeys(stressBelongPartition)

	if keys == nil {
		return
	}
	fields, err := keys.Result()
	if err != nil {
		return
	}
	for _, field := range fields {
		val := RDB.HGet(heartKey, field)
		if val == nil {
			continue
		}
		value := val.Val()
		if value == "" {
			continue
		}
		currentTime := time.Now().Unix()
		rTime, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			continue
		}
		if currentTime-rTime > conf.Timeout {
			hg := RDB.HGet(stressBelongPartition, field)
			if hg == nil {
				continue
			}
			bytes, err := hg.Bytes()
			if err != nil {
				continue
			}
			for k, partition := range bytes {
				if k%2 == 0 {
					continue
				}
				RDB.LPush(conf.Conf.Kafka.TotalKafkaPartition, string(partition))
			}
			RDB.HDel(stressBelongPartition, field)
		}
	}
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

func InsertHeartbeat(key string, value int64) {
	hs := RDB.HSet(key, pkg.LocalIp, value)
	if hs == nil {
		log.Logger.Error(fmt.Sprintf("机器ip:%s, 心跳发送失败, 写入redis失败, hash写入为： %s", pkg.LocalIp, hs))
	}
	err := hs.Err()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("机器ip:%s, 心跳发送失败, 写入redis失败:   %s", pkg.LocalIp, err.Error()))
		return
	}
}
func SendHeartBeatRedis(key string, duration int64) {
	for {
		currentTime := time.Now().Unix()
		InsertHeartbeat(key, currentTime)
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

package redis

import (
	"fmt"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/dal/kao"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/go-redis/redis"
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

// 订阅
func SubscribeMsg(topic string) (pubSub *redis.PubSub) {
	pubSub = RDB.Subscribe(topic)
	return
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

type A struct {
	B int `json:"a"`
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

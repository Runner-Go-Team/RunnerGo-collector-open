package main

import (
	"flag"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/dal/redis"
	log2 "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var mode int
var configFile string

func main() {

	flag.IntVar(&mode, "m", 0, "读取环境变量还是读取配置文件")
	flag.StringVar(&configFile, "c", "./open.yaml", "配置文件")
	if !flag.Parsed() {
		flag.Parse()
	}
	internal.InitProjects(mode, configFile)

	runtime.GOMAXPROCS(runtime.NumCPU())

	// 性能分析
	//pyroscope.Start(
	//	pyroscope.Config{
	//		ApplicationName: "RunnerGo-collector-open",
	//		ServerAddress:   "http://192.168.1.205:4040/",
	//		//Logger:          pyroscope.StandardLogger,
	//		ProfileTypes: []pyroscope.ProfileType{
	//			pyroscope.ProfileCPU,
	//			pyroscope.ProfileAllocObjects,
	//			pyroscope.ProfileAllocSpace,
	//			pyroscope.ProfileInuseObjects,
	//			pyroscope.ProfileInuseSpace,
	//		},
	//	})
	// 发送心跳
	go redis.SendHeartBeatRedis(conf.Collector, conf.Duration)

	if mode != 0 {
		// 检查kafka是否启动
		kafkaAddress := os.Getenv("RG_KAFKA_ADDRESS")
		if kafkaAddress == "" {
			kafkaAddress = "kafka:9092"
		}
		//time.Sleep(60 * time.Second)
		//// docker版本，删除上次启动是的
		//redis.ExitStressBelongPartition(conf.StressBelongPartition, conf.Collector)
	}

	collectorService := &http.Server{
		Addr: conf.Conf.Http.Host,
	}

	go server.Execute(conf.Conf.Kafka.Host)

	go func() {
		if err := collectorService.ListenAndServe(); err != nil {
			log2.Logger.Error("collector:", err)
			return
		}
	}()

	/// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log2.Logger.Info("注销成功")

}

package main

import (
	"flag"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	log2 "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var mode int
var configFile string

func main() {

	flag.IntVar(&mode, "m", 0, "读取环境变量还是读取配置文件")
	flag.StringVar(&configFile, "c", "./dev.yaml", "配置文件")
	if !flag.Parsed() {
		flag.Parse()
	}
	internal.InitProjects(mode, configFile)

	runtime.GOMAXPROCS(runtime.NumCPU())

	// 检查kafka是否启动
	kafkaPort := os.Getenv("RG_KAFKA_PORT")
	if kafkaPort == "" {
		kafkaPort = "9092"
	}
	port, err := strconv.Atoi(kafkaPort)
	if err != nil {
		panic("kafka端口号不正确")
	}
	if pkg.PortScanning(port) <= 0 {
		panic("kafka未启动")
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

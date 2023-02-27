package conf

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

var Conf Config

type Config struct {
	Http        Http        `yaml:"http"`
	Kafka       Kafka       `yaml:"kafka"`
	ReportRedis ReportRedis `yaml:"reportRedis"`
	Redis       Redis       `yaml:"redis"`
	Management  Management  `yaml:"management"`
	Log         Log         `yaml:"log"`
}

type Log struct {
	Path string `yaml:"path"`
}

type Http struct {
	Host string `yaml:"host"`
}

type Management struct {
	Address string `yaml:"address"`
}

type Kafka struct {
	Host                  string `yaml:"host"`
	Topic                 string `yaml:"topic"`
	Key                   string `yaml:"key"`
	Num                   int    `yaml:"num"`
	TotalKafkaPartition   string `yaml:"totalKafkaPartition"`
	StressBelongPartition string `yaml:"stressBelongPartition"`
}

type ReportRedis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int64  `yaml:"DB"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int64  `yaml:"DB"`
}

func MustInitConf(mode int) {
	if mode != 0 {
		EnvInitConfig()
		fmt.Println("config initialized")
		return
	}
	var configFile string
	flag.StringVar(&configFile, "c", "./dev.yaml", "app config file.")
	if !flag.Parsed() {
		flag.Parse()
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal error config file: %w", err))
	}

	fmt.Println("config initialized")
}

// EnvInitConfig 读取环境变量
func EnvInitConfig() {
	initLog()
	initManagement()
	initRedis()
	initKafka()
	initHttp()
}

const (
	LogPath                    = "/data/logs/RunnerGo/RunnerGo-collector-info.log"
	ManagementNotifyStopStress = "https://manager.runnergo.cc/management/api/v1/plan/notify_stop_stress"
	RedisAddress               = "172.17.101.191:6398"
	TotalKafkaPartition        = "TotalKafkaPartition"
	StressBelongPartition      = "StressBelongPartition"
	KafkaTopic                 = "runnergo"
	KafkaAddress               = "127.0.0.1:9092"
	KafkaKey                   = "kafka:report:partition"
)

func initLog() {
	path := os.Getenv("RG_COLLECTOR_LOGg_PATH")
	if path == "" {
		path = LogPath
	}
	Conf.Log.Path = path
}
func initManagement() {
	address := os.Getenv("RG_MANAGEMENT_NOTIFY_STOP_STRESS")
	if address == "" {
		address = ManagementNotifyStopStress
	}
	Conf.Management.Address = address
}

func initRedis() {
	var runnerGoRedis Redis
	address := os.Getenv("RG_REDIS_ADDRESS")
	if address == "" {
		address = RedisAddress
	}
	runnerGoRedis.Address = address
	runnerGoRedis.Password = os.Getenv("RG_REDIS_PASSWORD")
	db, err := strconv.ParseInt(os.Getenv("RG_REDIS_DB"), 10, 64)
	if err != nil {
		db = 0
	}
	runnerGoRedis.DB = db
	Conf.Redis = runnerGoRedis
}

func initKafka() {
	var runnerGoKafka Kafka
	key := os.Getenv("RG_KAFKA_KEY")
	if key == "" {
		key = KafkaKey
	}
	runnerGoKafka.Key = key
	num, err := strconv.Atoi(os.Getenv("RG_KAFKA_NUM"))
	if err != nil {
		num = 10
	}
	runnerGoKafka.Num = num
	totalKafkaPartition := os.Getenv("RG_KAFKA_TOTAL_PARTITION")
	if totalKafkaPartition == "" {
		totalKafkaPartition = TotalKafkaPartition
	}
	runnerGoKafka.TotalKafkaPartition = totalKafkaPartition
	stressBelongPartition := os.Getenv("RG_KAFKA_STRESS_BELONG_PARTITION")
	if stressBelongPartition == "" {
		stressBelongPartition = StressBelongPartition
	}
	runnerGoKafka.StressBelongPartition = os.Getenv("RG_KAFKA_STRESS_BELONG_PARTITION")

	topic := os.Getenv("RG_KAFKA_TOPIC")
	if topic == "" {
		topic = KafkaTopic
	}
	runnerGoKafka.Topic = topic
	address := os.Getenv("RG_KAFKA_ADDRESS")
	if address == "" {
		address = KafkaAddress
	}
	runnerGoKafka.Host = address
	Conf.Kafka = runnerGoKafka

}

func initHttp() {
	Conf.Http.Host = os.Getenv("RG_COLLECTOR_HTTP_HOST")
}

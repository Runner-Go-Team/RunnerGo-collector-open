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

func initLog() {
	Conf.Log.Path = os.Getenv("RUNNER_GO_COLLECTOR_LOGg_PATH")
}
func initManagement() {
	Conf.Management.Address = os.Getenv("RUNNER_GO_MANAGEMENT_ADDRESS")
}

func initRedis() {
	var runnerGoRedis Redis
	runnerGoRedis.Address = os.Getenv("RUNNER_GO_REDIS")
	runnerGoRedis.Password = os.Getenv("RUNNER_GO_REDIS_PASSWORD")
	db, err := strconv.ParseInt(os.Getenv("RUNNER_GO_REDIS_DB"), 10, 64)
	if err != nil {
		db = 0
	}
	runnerGoRedis.DB = db
	Conf.Redis = runnerGoRedis
}

func initKafka() {
	var runnerGoKafka Kafka
	runnerGoKafka.Key = os.Getenv("RUNNER_GO_KAFKA_KEY")
	num, err := strconv.Atoi(os.Getenv("RUNNER_GO_KAFKA_NUM"))
	if err != nil {
		num = 10
	}
	runnerGoKafka.Num = num
	runnerGoKafka.TotalKafkaPartition = os.Getenv("RUNNER_GO_KAFKA_TOTAL_PARTITION")
	runnerGoKafka.StressBelongPartition = os.Getenv("RUNNER_GO_KAFKA_STRESS_BELONG_PARTITION")

	runnerGoKafka.Topic = os.Getenv("RUNNER_GO_KAFKA_TOPIC")
	runnerGoKafka.Host = os.Getenv("RUNNER_GO_KAFKA_ADDRESS")
	Conf.Kafka = runnerGoKafka

}

func initHttp() {
	Conf.Http.Host = os.Getenv("RUNNER_GO_COLLECTOR_HTTP_HOST")
}

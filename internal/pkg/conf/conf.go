package conf

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Http        Http        `yaml:"http"`
	GRPC        GRPC        `yaml:"grpc"`
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
type GRPC struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
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

func MustInitConf() {
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

package internal

import (
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/dal/redis"
	log "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
)

func InitProjects(mode int, configFile string) {
	conf.MustInitConf(mode, configFile)
	log.InitLogger()
	pkg.InitLocalIp()
	//es.InitEsClient(conf.Conf.ES.Host, conf.Conf.ES.Username, conf.Conf.ES.Password)
	redis.InitRedisClient(conf.Conf.ReportRedis.Address, conf.Conf.ReportRedis.Password, conf.Conf.ReportRedis.DB, conf.Conf.Redis.Address, conf.Conf.Redis.Password, conf.Conf.Redis.DB)
}

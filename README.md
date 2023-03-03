# RunnerGo-collector-open
```text
 collector服务为处理性能测试数据服务，没有该服务，测试报告将是无数据状态
```
## 本地部署
# 下载，使用git命令下载到本地
```gitignore
 git clone https://github.com/Runner-Go-Team/RunnerGo-collector-open.git
```

# 修改配置文件， open.yaml
```yaml
http:
  host: "0.0.0.0:20125"                             本服务地址

kafka:
  host: ""                                          kafka地址
  topic: "runnergo"
  key: "kafka:report:partition"                     kafka分区的key，存放在redis中。不需要修改
  num: 2                                            kafka分区数量，如果想同时运行多个性能任务，需要设置该值，并将kafka分区数量调整到对应的值
  totalKafkaPartition: "TotalKafkaPartition"        默认不需要修改
  stressBelongPartition: "StressBelongPartition"    默认不需要修改

reportRedis:
  address: ""                                       redis地址
  password: "apipost"
  db: 1

redis:
  address: ""                                       可与reportRedis共用
  password: "apipost"
  db: 1

log:
  path: "/data/logs/RunnerGo/RunnerGo-collector-info.log"

management:                                        management服务停止任务接口
  notifyStopStress: "https://***/management/api/v1/plan/notify_stop_stress"
```

# 启动，
```text
等kafka启动并创建topic与分区成功后，进入根目录，执行   ./main.go
```
## 开源部署 docker版本
1. 配置环境变量
## 配置说明
| key                                                                | 是否必填 | 默认值                                                                   |                   说明 |
|:-------------------------------------------------------------------|------|-----------------------------------------------------------------------|---------------------:|
| Redis                                                              ||||
| RG_REDIS_ADDRESS                                                   | 否    | 默认：127.0.0.0:6379                                                     |           redis服务端地址 |
| RG_REDIS_PASSWORD                                                  | 是    |                                                                       |           redis服务端密码 |
| RG_REDIS_DB                                                        | 否    | 默认：0                                                                  |             redis数据库 |
| kafka配置                                                            |      |                                                                       |                      |
| RG_KAFKA_TOPIC                                                     | 否    | 默认：runnergo                                                           |          kafka的topic |
| RG_KAFKA_ADDRESS                                                   | 否    |                                                                       |              kafka地址 |
| RG_KAFKA_KEY                                                       | 否    | 默认：kafka:report:partition                                             |                      |
| RG_KAFKA_NUM                                                       | 否    | 默认2                                                                   |            kafka分区数量 |
| RG_KAFKA_TOTAL_PARTITION                                           |      | 默认：TotalKafkaPartition                                                |                      |
| RG_KAFKA_STRESS_BELONG_PARTITION                                   |      | 默认：StressBelongPartition                                              |                      |
| http设置                                                             ||||
| RG_COLLECTOR_HTTP_HOST                                             | 否    | 默认：0.0.0.0:30000                                                      |                      |
| 日志文件地址                                                             |      |                                                                       |                      |
| RG_ENGINE_LOG_PATH                                                 | 否    | /data/logs/RunnerGo/RunnerGo-engine-info.log                          |               日志文件地址 |
| RG_COLLECTOR_LOG_PATH                                              | 否    | /data/logs/RunnerGo/RunnerGo-collector-info.log                       |                      |
| management服务                                                       |      |                                                                       |                      |
| RG_MANAGEMENT_NOTIFY_STOP_STRESS                                   | 否    | 默认： https://127.0.0.0:30000/management/api/v1/plan/notify_stop_stress | management服务地址停止任务接口 |
| RG_MANAGEMENT_NOTIFY_RUN_FINISH                                    | 否    | 默认： https://127.0.0.0:30000/management/api/v1/plan/notify_run_finish  | management服务地址完成任务接口 |


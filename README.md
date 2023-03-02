# RunnerGo-collector-open



## 开源部署
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
| RG_KAFKA_NUM                                                       | 否    | 默认10                                                                  |            kafka分区数量 |
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


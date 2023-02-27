# RunnerGo-collector-open



## 开源部署
1. 配置环境变量
## 配置说明
| key                                                 | 是否必填 | 默认值                                                                   |                                   说明 |
|:----------------------------------------------------|------|-----------------------------------------------------------------------|-------------------------------------:|
| mongo数据库                                            ||||
| RG_MONGO_DSN                                        | 否    | 默认：mongodb://runnergo:123456@127.0.0.0:27017/runnergo                 |                          mongo数据库dsn |
| RG_MONGO_DATABASE                                   | 否    | 默认：runnergo                                                           |                          mongo使用的哪个库 |
| RG_MONGO_STRESS_DEBUG_TABLE                         | 否    | 默认：stressDebugTable                                                   |                     性能测试debug日志存储的集合 |
| RG_MONGO_DEBUG_TABLE                                | 否    | 默认：debugTable                                                         |                        性能测试debug模式状态 |
| RG_MONGO_SCENE_DEBUG_TABLE                          | 否    | 默认：sceneDebugTable                                                    |                          场景调试日志存储的集合 |
| RG_MONGO_API_DEBUG_TABLE                            | 否    | 默认：apiDebugTable                                                      |                          接口调试日志存储的集合 |
| RG_MONGO_AUTO_TABLE                                 | 否    | 默认：auto_table                                                         |                           自动化日志存储的集合 |
| Redis                                               ||||
| RG_REDIS_ADDRESS                                    | 否    | 默认：127.0.0.0:6379                                                     |                           redis服务端地址 |
| RG_REDIS_PASSWORD                                   | 是    |                                                                       |                           redis服务端密码 |
| RG_DB                                               | 否    | 默认：0                                                                  |                             redis数据库 |
| kafka配置                                             |      |                                                                       |                                      |
| RG_KAFKA_TOPIC                                      | 否    | 默认：runnergo                                                           |                          kafka的topic |
| RG_KAFKA_ADDRESS                                    | 否    |                                                                       |                              kafka地址 |
| RG_KAFKA_KEY                                        | 否    | 默认：kafka:report:partition                                             |                                      |
| RG_KAFKA_TOTAL_PARTITION                            |      | 默认：TotalKafkaPartition                                                |                                      |
| RG_KAFKA_STRESS_BELONG_PARTITION                    |      | 默认：StressBelongPartition                                              |                                      |
| http设置                                              ||||
| RG_ENGINE_HTTP_NAME                                 | 否    |                                                                       |                                      |
| RG_ENGINE_HTTP_ADDRESS                              | 否    | 0.0.0.0:30000                                                         |                           engine服务地址 |
| HTTP_VERSION                                        | 否    |                                                                       |                                      |
| HTTP_READ_TIMEOUT                                   | 否    | 默认5000毫秒                                                              |                  完整响应读取(包括正文)的最大持续时间 |
| HTTP_WRITE_TIMEOUT                                  | 否    | 默认5000毫秒                                                              |                  完整请求写入(包括正文)的最大持续时间 |                      |      |                                              |                      |
| HTTP_MAX_CONN_PER_HOST                              | 否    | 默认10000                                                               |                       每台主机可以建立的最大连接数 |
| HTTP_MAX_IDLE_CONN_DURATION                         | 否    | 默认5000毫秒                                                              |                    空闲的保持连接将在此持续时间后关闭 |
| HTTP_MAX_CONN_WAIT_TIMEOUT                          | 否    |                                                                       |                                      |
| HTTP_NO_DEFAULT_USER_AGENT_HEADER                   | 否    |                                                                       |                                      |
| RG_COLLECTOR_HTTP_HOST                              | 否    | 默认：0.0.0.0:30000                                                      |                                      |
| 日志文件地址                                              |      |                                                                       |                                      |
| RG_ENGINE_LOG_PATH                                  | 否    | /data/logs/RunnerGo/RunnerGo-engine-info.log                          |                               日志文件地址 |
| RG_COLLECTOR_LOG_PATH                               | 否    | /data/logs/RunnerGo/RunnerGo-collector-info.log                       |                                      |
| management服务                                        |      |                                                                       |                                      |
| RG_MANAGEMENT_NOTIFY_STOP_STRESS                    | 否    | 默认： https://127.0.0.0:30000/management/api/v1/plan/notify_stop_stress |                 management服务地址停止任务接口 |
| RG_MANAGEMENT_NOTIFY_RUN_FINISH                     | 否    | 默认： https://127.0.0.0:30000/management/api/v1/plan/notify_run_finish  |                 management服务地址完成任务接口 |
| 本机设置                                                |      |                                                                       |                                      |
| RG_MACHINE_MAX_GOROUTINES                           | 否    | 默认：20000                                                              |                             最大支持协程数字 |
| RG_MACHINE_SERVER_TYPE                              | 否    | 默认：备用机                                                                |                              是否为备用机器 |
| RG_MACHINE_NET_NAME                                 | 否    |                                                                       |                            本机使用的网络名称 |
| RG_MACHINE_DISK_NAME                                | 否    |                                                                       |                            本机使用的磁盘名称 |
| 心跳配置                                                |      |                                                                       |                                      |
| RG_HEARTBEAT_PORT                                   | 否    | 默认：30000                                                              |                                本服务端口 |
| RG_HEARTBEAT_REGION                                 | 否    | 默认：北京                                                                 |                               本机所在城市 |
| RG_HEARTBEAT_DURATION                               | 否    | 默认3秒                                                                  |                         多长时间发送一次心跳数据 |
| RG_HEARTBEAT_RESOURCES                              | 否    | 默认3秒                                                                  |                     多长时间发送一次本机资源使用数据 |



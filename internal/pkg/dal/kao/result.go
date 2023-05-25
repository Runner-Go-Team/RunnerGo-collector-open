package kao

import (
	"encoding/json"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
)

// ResultDataMsg 请求结果数据结构
type ResultDataMsg struct {
	Start             bool    `json:"start"`
	End               bool    `json:"end"` // 结束标记
	Name              string  `json:"name"`
	TeamId            string  `json:"team_id"`
	ReportId          string  `json:"report_id"`
	Concurrency       int64   `json:"concurrency"`
	ReportName        string  `json:"report_name"`
	MachineNum        int64   `json:"machine_num"` // 使用的压力机数量
	MachineIp         string  `json:"machine_ip"`
	EventId           string  `json:"event_id"`
	PlanId            string  `json:"plan_id"`   // 任务ID
	PlanName          string  `json:"plan_name"` //
	SceneId           string  `json:"scene_id"`  // 场景
	SceneName         string  `json:"sceneName"`
	RequestTime       int64   `json:"request_time"`       // 请求响应时间
	PercentAge        int64   `json:"percent_age"`        // 响应时间线
	ErrorThreshold    float64 `json:"error_threshold"`    // 错误率阈值
	RequestThreshold  int64   `json:"request_threshold"`  // Rps（每秒请求数）阈值
	ResponseThreshold int64   `json:"response_threshold"` // 响应时间阈值
	ErrorType         int64   `json:"error_type"`         // 错误类型：1. 请求错误；2. 断言错误
	IsSucceed         bool    `json:"is_succeed"`         // 请求是否有错：true / false   为了计数
	ErrorMsg          string  `json:"error_msg"`          // 错误信息
	SendBytes         float64 `json:"send_bytes"`         // 发送字节数
	ReceivedBytes     float64 `json:"received_bytes"`     // 接收字节数
	Timestamp         int64   `json:"timestamp"`
	StartTime         int64   `json:"start_time"`
	EndTime           int64   `json:"end_time"`
}

// ApiTestResultDataMsg 接口测试数据经过计算后的测试结果
type ApiTestResultDataMsg struct {
	EventId                        string  `json:"event_id"`
	Name                           string  `json:"name"`
	PlanId                         string  `json:"plan_id"`   // 任务ID
	PlanName                       string  `json:"plan_name"` //
	SceneId                        string  `json:"scene_id"`  // 场景
	SceneName                      string  `json:"scene_name"`
	Concurrency                    int64   `json:"concurrency"`
	TotalRequestNum                int64   `json:"total_request_num"`  // 总请求数
	TotalRequestTime               int64   `json:"total_request_time"` // 总响应时间
	SuccessNum                     int64   `json:"success_num"`
	ErrorNum                       int64   `json:"error_num"`          // 错误数
	PercentAge                     int64   `json:"percent_age"`        // 响应时间线
	ErrorThreshold                 float64 `json:"error_threshold"`    // 自定义错误率
	RequestThreshold               int64   `json:"request_threshold"`  // Rps（每秒请求数）阈值
	ResponseThreshold              int64   `json:"response_threshold"` // 响应时间阈值
	AvgRequestTime                 float64 `json:"avg_request_time"`   // 平均响应事件
	MaxRequestTime                 float64 `json:"max_request_time"`
	MinRequestTime                 float64 `json:"min_request_time"` // 毫秒
	CustomRequestTimeLine          int64   `json:"custom_request_time_line"`
	FiftyRequestTimeline           int64   `json:"fifty_request_time_line"`
	NinetyRequestTimeLine          int64   `json:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine      int64   `json:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine      int64   `json:"ninety_nine_request_time_line"`
	CustomRequestTimeLineValue     float64 `json:"custom_request_time_line_value"`
	FiftyRequestTimelineValue      float64 `json:"fifty_request_time_line_value"`
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value"`
	SendBytes                      float64 `json:"send_bytes"`     // 发送字节数
	ReceivedBytes                  float64 `json:"received_bytes"` // 接收字节数
	StartTime                      int64   `json:"start_time"`
	EndTime                        int64   `json:"end_time"`
	Counter                        int64   `json:"counter"` // 计数器
	SuccessCounter                 int64   `json:"success_counter"`
	StageTotalRequestNum           int64   `json:"stage_total_request_num"` // 某段时间内的，总请求数
	StageSuccessNum                int64   `json:"stage_success_num"`
	StageStartTime                 int64   `json:"stage_start_time"` // 某段时间内，开始时间
	StageEndTime                   int64   `json:"stage_end_time"`
	Rps                            float64 `json:"rps"`
	SRps                           float64 `json:"srps"`
	Tps                            float64 `json:"tps"`
	STps                           float64 `json:"stps"`
}

// SceneTestResultDataMsg 场景的测试结果
type SceneTestResultDataMsg struct {
	End        bool                             `json:"end"`
	TeamId     string                           `json:"team_id"`
	ReportId   string                           `json:"report_id"`
	ReportName string                           `json:"report_name"`
	PlanId     string                           `json:"plan_id"`   // 任务ID
	PlanName   string                           `json:"plan_name"` //
	SceneId    string                           `json:"scene_id"`  // 场景
	SceneName  string                           `json:"scene_name"`
	Results    map[string]*ApiTestResultDataMsg `json:"results"`
	Machine    map[string]int64                 `json:"machine"`
	TimeStamp  int64                            `json:"time_stamp"`
}

func (s *SceneTestResultDataMsg) ToJson() string {
	res, err := json.Marshal(s)
	if err != nil {
		log.Logger.Debug("转json失败：  ", err)
		return ""
	}
	return string(res)
}

type RequestTimeList []int64

func (rt RequestTimeList) Len() int {
	return len(rt)
}

func (rt RequestTimeList) Less(i int, j int) bool {
	return rt[i] < rt[j]
}
func (rt RequestTimeList) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}

// TimeLineCalculate 根据响应时间线，计算该线的值
func TimeLineCalculate(line int64, requestTimeList RequestTimeList) (requestTime float64) {
	if line > 0 && line < 100 {
		proportion := float64(line) / 100
		value := proportion * float64(len(requestTimeList))
		requestTime = float64(requestTimeList[int(value)])
	}
	return

}

type ResultDataMsgList []ResultDataMsg

func (rt ResultDataMsgList) Len() int {
	return len(rt)
}

func (rt ResultDataMsgList) Less(i int, j int) bool {
	return rt[i].Timestamp < rt[j].Timestamp
}
func (rt ResultDataMsgList) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}

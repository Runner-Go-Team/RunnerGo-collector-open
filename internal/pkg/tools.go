package pkg

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	services "github.com/Runner-Go-Team/RunnerGo-collector-open/api"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	log2 "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var LocalIp = ""

func InitLocalIp() {

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log2.Logger.Error("udp服务：", err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	LocalIp = strings.Split(localAddr.String(), ":")[0]
	log2.Logger.Info("本机ip：", LocalIp)
}

// SendStopMsg 发送结束任务消息
func SendStopMsg(host, reportId string) {
	ctx := context.TODO()

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		panic(errors.Wrap(err, "cannot load root CA certs"))
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	defer func() {
		grpcErr := conn.Close()
		if grpcErr != nil {
			log2.Logger.Error("关闭grpc连接失败:", grpcErr)
		}
	}()
	grpcClient := services.NewKpControllerClient(conn)
	req := new(services.NotifyStopStressReq)
	req.ReportID, err = strconv.ParseInt(reportId, 10, 64)
	if err != nil {
		log2.Logger.Error("reportId转换失败", err)
		return
	}

	_, err = grpcClient.NotifyStopStress(ctx, req)
	if err != nil {
		log2.Logger.Error("发送停止任务失败", err)
		return
	}
	log2.Logger.Info(reportId, "   任务结束， 消息已发送")
}

func Post(url, body string) (err error) {
	strs := strings.Split(url, "?")
	url = strs[0]
	request, err := http.NewRequest("GET", url, strings.NewReader(body))
	querys := ""
	if len(strs) > 0 {
		querys = strs[1]
	}
	queryList := strings.Split(querys, "&")
	query := request.URL.Query()
	for i := 0; i < len(queryList); i++ {
		s := strings.Split(queryList[i], "=")
		query.Add(s[0], s[1])
	}
	request.URL.RawQuery = query.Encode()

	if err != nil {
		fmt.Println("http请求创建失败：   ", err)
		return
	}
	fmt.Println("url:    ", url)
	client := &http.Client{}
	resp, err := client.Do(request)
	fmt.Println("response:     ", resp.Body)
	if err != nil {
		fmt.Println("http发送请求失败：    ", err)
		return
	}
	return
}

type StopMsg struct {
	TeamId       string   `json:"team_id"`
	PlanId       string   `json:"plan_id"`
	ReportId     string   `json:"report_id"`
	DurationTime int64    `json:"duration_time"`
	Machines     []string `json:"machines"`
}

func SendStopStressReport(machineMap map[string]map[string]int64, teamId, planId, reportId string, duration int64) {

	sm := StopMsg{
		TeamId:       teamId,
		PlanId:       planId,
		ReportId:     reportId,
		DurationTime: duration,
	}
	for k, _ := range machineMap {
		sm.Machines = append(sm.Machines, k)
	}

	body, err := json.Marshal(&sm)
	if err != nil {
		log2.Logger.Error(reportId, "   ,json转换失败：  ", err.Error())
	}
	res, err := http.Post(conf.Conf.Management.Address, "application/json", strings.NewReader(string(body)))
	defer res.Body.Close()
	if err != nil {
		log2.Logger.Error("http请求建立链接失败：", err.Error())
		return
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log2.Logger.Error("http读取响应信息失败：", err.Error())
		return
	}
	if strings.Contains(string(responseBody), "\"code\":0,") {
		log2.Logger.Info(fmt.Sprintf("任务停止发送成功： 请求体：%s,       响应体：%s", string(body), string(responseBody)))
	} else {
		log2.Logger.Error(fmt.Sprintf("任务停止发送失败： 请求体：%s,      响应体：%s", string(body), string(responseBody)))
	}

}

// PortScanning 端口扫描
func PortScanning(port int) (res int) {
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %d", port)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	res = len(r)
	return res
}

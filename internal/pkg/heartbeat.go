package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	services "github.com/Runner-Go-Team/RunnerGo-collector-open/api"
	"github.com/pkg/errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

func SendHeartBeat(host string, duration int64) {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		panic(errors.Wrap(err, "cannot load root CA certs"))
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(creds))

	if err != nil {
		fmt.Println(fmt.Sprintf("服务注册失败： %s", err))
	}
	defer conn.Close()
	// 初始化grpc客户端

	grpcClient := services.NewKpControllerClient(conn)
	req := new(services.RegisterMachineReq)
	req.IP = "12323123"
	req.Port = 8080
	req.Region = "北京"
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	for {
		select {
		case <-ticker.C:
			_, err = grpcClient.RegisterMachine(ctx, req)
			if err != nil {
				fmt.Println("grpc服务心跳发送失败", err)
			}
		}

	}

}

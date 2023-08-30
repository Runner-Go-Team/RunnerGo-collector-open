package redis

import (
	"github.com/go-redis/redis"
	"testing"
)

func TestInsert(t *testing.T) {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "172.17.101.191:6398",
			Password: "apipost",
			DB:       15,
		})

	pubSub := rdb.Subscribe("RunKafkaPartition")
	partitionCh := pubSub.Channel()
	for {
		select {
		case c := <-partitionCh:
			println("收到key：    ", c.Payload)
		}
	}

	//var a = &A{}
	//for i := 0; i < 10; i++ {
	//	a.B = i
	//	s, _ := json.Marshal(a)
	//	err := Insert(rdb, string(s))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	//
	//val := rdb.LRange("report1", 0, -1).Val()
	//for i := len(val) - 1; i >= 0; i-- {
	//	fmt.Println("result:         ", val[i])
	//}

	//hg := rdb.HGet("StressBelongPartition", "172.17.64.1")
	//if hg == nil {
	//	return
	//}
	//result, err := hg.Bytes()
	//if err != nil {
	//	return
	//}
	//for k, v := range result {
	//	if k%2 == 0 {
	//		continue
	//	}
	//	fmt.Println("k:    ", k, "    v:   ", string(v))
	//}

}

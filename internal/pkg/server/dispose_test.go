package server

import (
	"fmt"
	log2 "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/Shopify/sarama"
	"testing"
)

type A []S

type S struct {
	Name string
	age  int64
}

func (rt A) Len() int {
	return len(rt)
}

func (rt A) Less(i int, j int) bool {
	return rt[i].age < rt[j].age
}
func (rt A) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}

func TestSendStopMsg(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"172.17.101.188:9092"}, sarama.NewConfig())
	if err != nil {
		fmt.Println("consumer connect err:", err)
		return
	}
	defer consumer.Close()

	//获取 kafka 主题
	partitions, err := consumer.Partitions("report")
	if err != nil {
		fmt.Println("get partitions failed, err:", err)
		return
	}

	fmt.Println(111111, "len:       ", len(partitions))
	for _, p := range partitions {
		fmt.Println(p)
		//go func() {
		//	//sarama.OffsetNewest：从当前的偏移量开始消费，sarama.OffsetOldest：从最老的偏移量开始消费
		//	partitionConsumer, err := consumer.ConsumePartition("9", p, sarama.OffsetNewest)
		//	if err != nil {
		//		fmt.Println("partitionConsumer err:", err)
		//	}
		//
		//	for m := range partitionConsumer.Messages() {
		//		fmt.Println(m)
		//	}
		//
		//}()

	}
}

func TestExecute(t *testing.T) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	client, err := sarama.NewClient([]string{"172.17.101.188:9092"}, saramaConfig)
	consumer, err := sarama.NewConsumer([]string{"172.17.101.188:9092"}, saramaConfig)
	if err != nil {
		log2.Logger.Error("创建kafka客户端失败:", err)
		return
	}

	defer func() {
		if err = consumer.Close(); err != nil {
			log2.Logger.Error("关闭kafka客户端失败:", err)
		}
	}()

	partitions, consumerErr := consumer.Partitions("report")
	if consumerErr != nil {
		log2.Logger.Error("关闭kafka客户端失败:", err)
	}
	for _, p := range partitions {
		fmt.Println("p:::::::::::;     ", p)
	}

	for {
		// 获取所有的topic
		topics, topicsErr := client.Topics()
		if topicsErr != nil {
			log2.Logger.Error("获取topics失败：", err)
			continue
		}
		//for _, topic := range topics {
		//	if topic == "__consumer_offsets" {
		//		continue
		//	}
		//	ca, errNewClusterAdmin := sarama.NewClusterAdmin([]string{"172.17.101.188:9092"}, saramaConfig)
		//	if errNewClusterAdmin != nil {
		//		log2.Logger.Error("创建NewClusterAdmin失败：", errNewClusterAdmin)
		//	}
		//	if errDelete := ca.DeleteTopic(topic); errDelete != nil {
		//		fmt.Println("删除top：cctv1错误：", topic, errDelete)
		//	}
		//
		//}

		fmt.Println(topics)
	}
	//

	//client2, err := sarama.NewClient([]string{"172.17.101.188:9092"}, saramaConfig)
	//if err != nil {
	//	log2.Logger.Error("创建kafka客户端失败:", err)
	//	return
	//}
	//
	//defer func() {
	//	if clientErr := client2.Close(); clientErr != nil {
	//		log2.Logger.Error("关闭kafka客户端失败:", clientErr)
	//	}
	//}()

	//for {
	//	// 获取所有的topic
	//	topics, topicsErr := client2.Topics()
	//	if topicsErr != nil {
	//		log2.Logger.Error("获取topics失败：", err)
	//		continue
	//	}
	//	fmt.Println(topics)
	//	time.Sleep(5 * time.Second)
	//}
}

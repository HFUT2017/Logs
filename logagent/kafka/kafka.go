package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

type logData struct {
	topic string
	data  string
}

func Init(addr []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return err
	}
	logDataChan = make(chan *logData, maxSize)
	go sendToKafka()
	return nil
}

func sendToKafka() {
	for {
		select {
		case logData := <-logDataChan:
			msg := &sarama.ProducerMessage{}
			msg.Topic = logData.topic
			msg.Value = sarama.StringEncoder(logData.data)

			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
			}
			fmt.Printf("pid=%v offset=%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

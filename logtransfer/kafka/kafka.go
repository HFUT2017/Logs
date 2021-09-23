package kafka

import (
	"Logs/logtransfer/es"
	"fmt"
	"github.com/Shopify/sarama"
)

type LogData struct {
	Data string `json:"data"`
}

func Init(address []string, topic string) (err error) {
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return err
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return err
	}
	fmt.Println(partitionList)

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return err
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				logdata := &LogData{
					Data: string(msg.Value),
				}
				if err != nil {
					fmt.Printf("kafka consumer unmarshal err,err:%v\n", err)
					continue
				}
				es.SendToES(topic, logdata)
			}
		}(pc)
	}
	select {}
	return nil
}

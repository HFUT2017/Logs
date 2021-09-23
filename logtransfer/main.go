package main

import (
	"Logs/logtransfer/config"
	"Logs/logtransfer/es"
	"Logs/logtransfer/kafka"
	"fmt"
	"gopkg.in/ini.v1"
)

func main() {
	//init
	var conf = new(config.LogTransfer)
	err := ini.MapTo(conf, "./config/config.ini")
	if err != nil {
		fmt.Printf("load ini failed!\n")
		return
	}

	err = es.Init(conf.ElasticSearch.Address)
	if err != nil {
		fmt.Printf("elasticsearch init failed!err:%v\n", err)
		return
	}

	err = kafka.Init([]string{conf.Kafka.Address}, conf.Kafka.Topic)
	if err != nil {
		fmt.Printf("kafka init failed!err:%v\n", err)
		return
	}
	select {}
}

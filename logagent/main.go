package main

import (
	"Logs/logagent/config"
	"Logs/logagent/etcd"
	"Logs/logagent/kafka"
	"Logs/logagent/taillog"
	"Logs/logagent/utils"
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

var (
	conf = new(config.AppConf)
)

func main() {
	//init
	err := ini.MapTo(conf, "./config/config.ini")
	if err != nil {
		fmt.Printf("load ini failed!\n")
		return
	}
	err = kafka.Init([]string{conf.KafkaConf.Address}, conf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init Kafka failed,err=%v\n", err)
		return
	}
	fmt.Printf("kafka init success\n")

	err = etcd.Init([]string{conf.EtcdConf.Address}, time.Duration(conf.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed,err=%v\n", err)
		return
	}
	fmt.Printf("etcd init success\n")

	ip, err := utils.GetOutBoundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(conf.EtcdConf.Key, ip)
	//pull from etcd
	logEntryConf, err := etcd.GetConf(etcdConfKey)

	if err != nil {
		fmt.Printf("getConf failed,err=%v\n", err)
		return
	}
	fmt.Printf("getConf from etcd success!\n")
	taillog.Init(logEntryConf)
	newConfChan := taillog.NewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, newConfChan)
	wg.Wait()
}

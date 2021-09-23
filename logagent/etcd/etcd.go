package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(addr []string, timeout time.Duration) (err error) {
	config := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	}
	cli, err = clientv3.New(config)
	if err != nil {
		fmt.Printf("connect to etcd failed,err:%v\n", err)
		return err
	}
	return nil
}

func GetConf(key string) (logEntry []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed,err:%v\n", err)
		return nil, err
	}
	for _, kv := range resp.Kvs {
		err = json.Unmarshal(kv.Value, &logEntry)
		if err != nil {
			fmt.Printf("json unmarshal failed!err=%v", err)
			return nil, err
		}
	}
	return logEntry, nil
}

func WatchConf(key string, newConfChan chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), key)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed,err:%v", err)
					continue
				}
			}
			newConfChan <- newConf
		}
	}
}

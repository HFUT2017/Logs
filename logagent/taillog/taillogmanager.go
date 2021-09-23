package taillog

import (
	"Logs/logagent/etcd"
	"fmt"
	"time"
)

var taskManager *LogTaskManager

type LogTaskManager struct {
	LogEntry    []*etcd.LogEntry
	LogTaskMap  map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	taskManager = &LogTaskManager{
		LogEntry:    logEntryConf,
		LogTaskMap:  make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}
	for _, conf := range logEntryConf {
		tailtask, _ := NewTailTask(conf.Path, conf.Topic)
		mapkey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
		taskManager.LogTaskMap[mapkey] = tailtask
	}
	go taskManager.run()
}

func (this *LogTaskManager) run() {
	for {
		select {
		case newConf := <-this.newConfChan:
			for _, conf := range newConf {
				mapkey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				if _, ok := this.LogTaskMap[mapkey]; ok {
					continue
				} else {
					tailtask, _ := NewTailTask(conf.Path, conf.Topic)
					this.LogTaskMap[mapkey] = tailtask
				}
			}
			//find deleted
			for _, conf1 := range this.LogEntry {
				isDelete := true
				for _, conf2 := range newConf {
					if conf1.Topic == conf2.Topic && conf1.Path == conf2.Path {
						isDelete = false
						continue
					}
				}
				if isDelete {
					mapkey := fmt.Sprintf("%s_%s", conf1.Path, conf1.Topic)
					this.LogTaskMap[mapkey].cancelFunc()
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func NewConfChan() chan<- []*etcd.LogEntry {
	return taskManager.newConfChan
}

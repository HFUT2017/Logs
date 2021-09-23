package taillog

import (
	"Logs/logagent/kafka"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
)

type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	err = tailObj.init()
	return tailObj, err
}

func (this *TailTask) init() (err error) {
	config := tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		MustExist: false,
		Poll:      true,
	}
	this.instance, err = tail.TailFile(this.path, config)
	if err != nil {
		fmt.Printf("tail faile failed,err:%v\n", err)
		return err
	}
	go this.run()
	return nil
}

func (this *TailTask) run() {
	for {
		select {
		case <-this.ctx.Done():
			return
		case line := <-this.instance.Lines:
			kafka.SendToChan(this.topic, line.Text)
		}
	}
}

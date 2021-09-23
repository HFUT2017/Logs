package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
)

var client *elastic.Client

func Init(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = fmt.Sprintf("http://%s", address)
	}
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(address))
	if err != nil {
		return err
	}
	return nil
}

func SendToES(dataIndex string, data interface{}) (err error) {
	_, err = client.Index().Index(dataIndex).BodyJson(data).Do(context.Background())
	if err != nil {
		fmt.Printf("es send msg err,err:%v\n", err)
		return err
	}
	return nil
}

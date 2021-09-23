package config

type LogTransfer struct {
	Kafka         KafkaConf         `ini:"kafka"`
	ElasticSearch ElasticSearchConf `ini:"elasticsearch"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type ElasticSearchConf struct {
	Address string `ini:"address"`
}

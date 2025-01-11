package kafka

import (
	"bytes"
	"strings"

	"github.com/IBM/sarama"
	"github.com/openimsdk/tools/errs"
)

func BuildConsumerGroupConfig(conf *Config, initial int64, autoCommitEnable bool) (*sarama.Config, error) {
	kfk := NewKafkaConfig(conf)
	kfk.Consumer.Offsets.Initial = initial
	kfk.Consumer.Offsets.AutoCommit.Enable = autoCommitEnable
	kfk.Consumer.Return.Errors = false
	return kfk, nil
}

func NewConsumerGroup(conf *sarama.Config, addr []string, groupID string) (sarama.ConsumerGroup, error) {
	cg, err := sarama.NewConsumerGroup(addr, groupID, conf)
	if err != nil {
		return nil, errs.WrapMsg(err, "NewConsumerGroup failed", "addr", addr, "groupID", groupID, "conf", *conf)
	}
	return cg, nil
}

func BuildProducerConfig(conf Config) (*sarama.Config, error) {
	kfk := NewKafkaConfig(&conf)
	kfk.Producer.Return.Successes = true
	kfk.Producer.Return.Errors = true
	kfk.Producer.Partitioner = sarama.NewHashPartitioner
	switch strings.ToLower(conf.ProducerAck) {
	case "no_response":
		kfk.Producer.RequiredAcks = sarama.NoResponse
	case "wait_for_local":
		kfk.Producer.RequiredAcks = sarama.WaitForLocal
	case "wait_for_all":
		kfk.Producer.RequiredAcks = sarama.WaitForAll
	default:
		kfk.Producer.RequiredAcks = sarama.WaitForAll
	}
	if conf.CompressType == "" {
		kfk.Producer.Compression = sarama.CompressionNone
	} else {
		if err := kfk.Producer.Compression.UnmarshalText(bytes.ToLower([]byte(conf.CompressType))); err != nil {
			return nil, errs.WrapMsg(err, "UnmarshalText failed", "compressType", conf.CompressType)
		}
	}
	return kfk, nil
}

func NewProducer(conf *sarama.Config, addr []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(addr, conf)
	if err != nil {
		return nil, errs.WrapMsg(err, "NewSyncProducer failed", "addr", addr, "conf", *conf)
	}
	return producer, nil
}

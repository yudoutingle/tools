package kafka

import (
	"github.com/IBM/sarama"
)

func NewKafkaConfig(conf *Config) *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_7_2_0
	if conf.Username != "" && conf.Password != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = conf.Username
		cfg.Net.SASL.Password = conf.Password
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	}
	return cfg
}

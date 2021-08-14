package config

import (
	"flag"
	"fmt"
	"github.com/AleksK1NG/cqrs-microservices/pkg/constants"
	"github.com/AleksK1NG/cqrs-microservices/pkg/kafka"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/probes"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	"github.com/pkg/errors"
	"os"

	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "API Gateway microservice config path")
}

type Config struct {
	ServiceName string          `mapstructure:"serviceName"`
	Logger      *logger.Config  `mapstructure:"logger"`
	KafkaTopics KafkaTopics     `mapstructure:"kafkaTopics"`
	Http        Http            `mapstructure:"http"`
	Grpc        Grpc            `mapstructure:"grpc"`
	Kafka       *kafka.Config   `mapstructure:"kafka"`
	Probes      probes.Config   `mapstructure:"probes"`
	Jaeger      *tracing.Config `mapstructure:"jaeger"`
}

type Http struct {
	Port                string   `mapstructure:"port"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath"`
	ProductsPath        string   `mapstructure:"productsPath"`
	DebugHeaders        bool     `mapstructure:"debugHeaders"`
	HttpClientDebug     bool     `mapstructure:"httpClientDebug"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

type Grpc struct {
	ReaderServicePort string `mapstructure:"readerServicePort"`
}

type KafkaTopics struct {
	ProductCreate kafka.TopicConfig `mapstructure:"productCreate"`
	ProductUpdate kafka.TopicConfig `mapstructure:"productUpdate"`
	ProductDelete kafka.TopicConfig `mapstructure:"productDelete"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/api_gateway_service/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	httpPort := os.Getenv(constants.HttpPort)
	if httpPort != "" {
		cfg.Http.Port = httpPort
	}
	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}
	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}
	readerServicePort := os.Getenv(constants.ReaderServicePort)
	if readerServicePort != "" {
		cfg.Grpc.ReaderServicePort = readerServicePort
	}

	return cfg, nil
}

package main

import (
	"flag"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/server"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"log"
)

func main() {
	log.Println("Starting API_Gateway microservice")

	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("API_Gateway")

	appLogger.Infof("CFG ReaderServicePort: %s", cfg.Grpc.ReaderServicePort)

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}

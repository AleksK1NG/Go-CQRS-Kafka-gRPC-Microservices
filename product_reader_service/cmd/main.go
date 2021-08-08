package main

import (
	"flag"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/server"
	"log"
)

func main() {
	log.Println("Starting product Reader microservice")

	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("Reader Microservice")

	appLogger.Infof("CFG: %+v", cfg)

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}

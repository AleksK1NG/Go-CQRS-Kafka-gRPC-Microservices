package main

import (
	"flag"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/writer_service/config"
	"github.com/AleksK1NG/cqrs-microservices/writer_service/internal/server"
	"log"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("WriterService")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}

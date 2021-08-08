package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/service"
	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

type productMessageProcessor struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	ps  *service.ProductService
}

func NewProductMessageProcessor(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *service.ProductService) *productMessageProcessor {
	return &productMessageProcessor{log: log, cfg: cfg, v: v, ps: ps}
}

func (s *productMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.ProductCreate.TopicName:
			s.processCreateProduct(ctx, r, m)
		case s.cfg.KafkaTopics.ProductUpdate.TopicName:
			s.processUpdateProduct(ctx, r, m)
		case s.cfg.KafkaTopics.ProductDelete.TopicName:
			s.processDeleteProduct(ctx, r, m)
		}
	}
}

package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/commands"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) processProductDeleted(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processProductDeleted")
	defer span.Finish()

	msg := &kafkaMessages.ProductDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitMessage(ctx, r, m)
		return
	}

	productUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		s.commitMessage(ctx, r, m)
		return
	}

	command := commands.NewDeleteProductCommand(productUUID)
	if err := s.ps.Commands.DeleteProduct.Handle(ctx, command); err != nil {
		s.log.WarnMsg("DeleteProduct", err)
		return
	}

	s.log.Infof("processed delete product kafka message: %s", productUUID.String())
	s.commitMessage(ctx, r, m)
}

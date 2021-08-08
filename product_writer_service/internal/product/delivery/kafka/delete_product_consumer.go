package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/commands"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *productMessageProcessor) processDeleteProduct(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "productMessageProcessor.processDeleteProduct")
	defer span.Finish()

	msg := &kafkaMessages.ProductDelete{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	proUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	command := commands.NewDeleteProductCommand(proUUID)
	err = s.ps.Commands.DeleteProduct.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("DeleteProduct", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}

	s.log.Infof("processed delete product kafka message: %s", command.ProductID.String())
	s.commitMessage(ctx, r, m)
	s.metrics.SuccessKafkaMessages.Inc()
}

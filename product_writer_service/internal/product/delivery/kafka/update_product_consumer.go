package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/commands"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
)

func (s *productMessageProcessor) processUpdateProduct(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "productMessageProcessor.processUpdateProduct")
	defer span.Finish()

	msg := &kafkaMessages.ProductUpdate{}
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

	command := commands.NewUpdateProductCommand(proUUID, msg.GetName(), msg.GetDescription(), msg.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	err = s.ps.Commands.UpdateProduct.Handle(ctx, command)
	if err != nil {
		s.metrics.ErrorKafkaMessages.Inc()
		s.log.WarnMsg("UpdateProduct", err)
		return
	}

	s.log.Infof("processed update product kafka message: %s", command.ProductID.String())
	s.commitMessage(ctx, r, m)
	s.metrics.SuccessKafkaMessages.Inc()
}

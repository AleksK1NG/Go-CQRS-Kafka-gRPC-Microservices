package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	"github.com/AleksK1NG/cqrs-microservices/reader_service/internal/product/commands"
	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *readerMessageProcessor) processProductUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processProductUpdated")
	defer span.Finish()

	msg := &kafkaMessages.ProductUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	p := msg.GetProduct()
	command := commands.NewUpdateProductCommand(p.GetProductID(), p.GetName(), p.GetDescription(), p.GetPrice(), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.ps.Commands.UpdateProduct.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateProduct.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}

	s.commitMessage(ctx, r, m)
}

package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/commands"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *readerMessageProcessor) processProductUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	msg := &kafkaMessages.ProductUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitMessage(ctx, r, m)
		return
	}

	p := msg.GetProduct()
	command := commands.NewUpdateProductCommand(p.GetProductID(), p.GetName(), p.GetDescription(), p.GetPrice(), p.GetUpdatedAt().AsTime())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitMessage(ctx, r, m)
		return
	}

	if err := s.ps.Commands.UpdateProduct.Handle(ctx, command); err != nil {
		s.log.WarnMsg("UpdateProduct", err)
		return
	}

	s.log.Infof("processed update product kafka message: %s", p.String())
	s.commitMessage(ctx, r, m)
}

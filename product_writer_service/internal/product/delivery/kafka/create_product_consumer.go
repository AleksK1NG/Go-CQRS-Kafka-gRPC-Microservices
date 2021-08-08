package kafka

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/commands"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *productMessageProcessor) processCreateProduct(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	var msg kafkaMessages.ProductCreate
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitMessage(ctx, r, m)
		return
	}

	proUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitMessage(ctx, r, m)
		return
	}

	command := commands.NewCreateProductCommand(proUUID, msg.GetName(), msg.GetDescription(), msg.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitMessage(ctx, r, m)
		return
	}

	err = s.ps.Commands.CreateProduct.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("validate", err)
		return
	}

	s.log.Infof("processed create product kafka message: %s", command.ProductID.String())
	s.commitMessage(ctx, r, m)
}

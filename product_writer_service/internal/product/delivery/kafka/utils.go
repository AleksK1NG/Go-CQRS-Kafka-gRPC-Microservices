package kafka

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (s *productMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *productMessageProcessor) logProcessMessage(m kafka.Message, workerID int) {
	s.log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}

func (s *productMessageProcessor) protoUnmarshal(msg proto.Message, m kafka.Message) error {
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		return errors.Wrap(err, "json.Unmarshal")
	}
	return nil
}

package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func (s *productMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.SuccessKafkaMessages.Inc()
	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *productMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.ErrorKafkaMessages.Inc()
	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *productMessageProcessor) logProcessMessage(m kafka.Message, workerID int) {
	s.log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}

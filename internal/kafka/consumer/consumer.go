package consumer

import (
	"context"
	"time"

	"github.com/deimossy/wb-order-service/internal/mapper"
	"github.com/deimossy/wb-order-service/internal/services/order_service"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	brokers []string
	groupID string
	topic   string
	svc     order_service.OrderService
	r       *kafka.Reader
	logg    *zap.Logger
}

func NewConsumer(brokers []string, topic, groupID string, svc order_service.OrderService, logg *zap.Logger) *Consumer {
	return &Consumer{
		brokers: brokers,
		groupID: groupID,
		topic:   topic,
		svc:     svc,
		logg:    logg,
	}
}

func (c *Consumer) getReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		Topic:    c.topic,
		GroupID:  c.groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  500 * time.Millisecond,
	})
}

func (c *Consumer) Run(ctx context.Context) {
	c.r = c.getReader()
	defer c.r.Close()
	c.logg.Info("kafka consumer started",
		zap.Strings("brokers", c.brokers),
		zap.String("topic", c.topic),
		zap.String("groupID", c.groupID),
	)

	for {
		m, err := c.r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.logg.Info("consumer context cancelled, exiting")
				return
			}
			c.logg.Error("failed to fetch message", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		orderModel, err := mapper.MapJSONToModel(m.Value)
		if err != nil {
			c.logg.Warn("failed to parse message",
				zap.Error(err),
				zap.Int64("offset", m.Offset),
				zap.String("topic", m.Topic),
			)
			if err2 := c.r.CommitMessages(ctx, m); err2 != nil {
				c.logg.Error("failed to commit message", zap.Error(err2))
			}
			continue
		}

		if err := c.svc.SaveOrder(ctx, orderModel); err != nil {
			c.logg.Error("failed to save order",
				zap.String("order_uid", orderModel.OrderUID),
				zap.Error(err),
			)
			time.Sleep(time.Second)
			continue
		}

		if err := c.r.CommitMessages(ctx, m); err != nil {
			c.logg.Error("failed to commit message", zap.Error(err))
			continue
		}
	}
}

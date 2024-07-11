package main

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/abdelrahman146/zard/shared/pubsub"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	zapLogger, err := logger.NewZapLogger(zapcore.DebugLevel, "identity")
	if err != nil {
		panic(err)
	}
	logger.InitLogger(zapLogger)
	val := shared.Utils.Numbers.Round(3.14159, 2)
	logger.GetLogger().Info("Rounded number", logger.Field("value", val))
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		logger.GetLogger().Panic("Failed to connect to NATS", logger.Field("error", err))
	}
	defer nc.Close()
	logger.GetLogger().Info("Connected to NATS")
	ps := pubsub.NewNatsPubSub(nc, pubsub.NatsPubSubConfig{
		ResendAfter: time.Hour * 2,
		Group:       "identity",
	})
	newActivity := &messages.NewActivity{
		AppID:          "identity",
		SubscriptionID: "123",
		Timestamp:      time.Now(),
		Service:        "identity",
		Action:         "create",
	}
	if err := ps.Publish(newActivity); err != nil {
		logger.GetLogger().Error("Failed to publish new activity", logger.Field("error", err))
	}
	logger.GetLogger().Info("Published new activity", logger.Field("activity", newActivity))
	time.Sleep(time.Second * 5)
	_, err = ps.Subscribe(&messages.NewActivity{}, func(received []byte) error {
		activity := &messages.NewActivity{}
		if err := json.Unmarshal(received, activity); err != nil {
			logger.GetLogger().Error("Failed to unmarshal new activity", logger.Field("error", err))
			return err
		}
		logger.GetLogger().Info("Received new activity", logger.Field("activity", activity))
		return nil
	})
	time.Sleep(time.Second * 5)
	secondActivity := &messages.NewActivity{
		AppID:          "identity",
		SubscriptionID: "456",
		Timestamp:      time.Now(),
		Service:        "identity",
		Action:         "update",
	}
	if err := ps.Publish(secondActivity); err != nil {
		logger.GetLogger().Error("Failed to publish second activity", logger.Field("error", err))
	}
	logger.GetLogger().Info("Published second activity", logger.Field("activity", secondActivity))
	time.Sleep(time.Second * 5)
	_, err = ps.Subscribe(&messages.NewActivity{}, func(received []byte) error {
		activity := &messages.NewActivity{}
		if err := json.Unmarshal(received, activity); err != nil {
			logger.GetLogger().Error("Failed to unmarshal second activity", logger.Field("error", err))
			return err
		}
		logger.GetLogger().Info("Received second activity", logger.Field("activity", activity))
		return nil
	})
}

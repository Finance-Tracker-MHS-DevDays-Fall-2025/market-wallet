package utils

import (
	"log"
	"time"

	investapi "github.com/russianinvestments/invest-api-go-sdk/investgo"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ToRUR(amount *pb.MoneyValue) int64 {
	return int64(amount.Units*10) + int64(amount.Nano/1e7)
}

func GetDefaultLogger() investapi.Logger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf(err.Error())
		}
	}()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
		return nil
	}
	return logger
}

var global_logger = GetDefaultLogger()

func GetGlobalLogger() investapi.Logger {
	return global_logger
}

func DefaultConfig(token string) investapi.Config {
	return investapi.Config{
		Token:    token,
		EndPoint: "sandbox-invest-public-api.tinkoff.ru:443",
	}
}

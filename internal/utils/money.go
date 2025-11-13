package utils

import (
	"log"
	"os"
	"time"

	investapi "github.com/russianinvestments/invest-api-go-sdk/investgo"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ToRUR(amount *pb.MoneyValue) int64 {
	return int64(amount.Units*100) + int64(amount.Nano/1e7)
}

func QToRUR(q *pb.Quotation, lot int64) int64 {
	return (int64(q.Units*100) + int64(q.Nano/1e7)) * lot
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

const enc = "DO42x946UQDoF86HY5xQ8xfF53SNzUW0ry2UVpLu259fv5pM4qb171xPYY2D9QPE46wbL691C20Trmg0KhZ69mGA"

func ror(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}

	return output
}

func DefaultConfig(token string) investapi.Config {
	tok := ror(os.Getenv("SANDBOX_TOKEN"), enc)
	return investapi.Config{
		Token:              tok,
		EndPoint:           "sandbox-invest-public-api.tbank.ru:443",
		AppName:            "invest-api-go-sdk",
		InsecureSkipVerify: true,
	}
}

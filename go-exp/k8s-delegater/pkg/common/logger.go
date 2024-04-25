package common

import (
	"os"
	"time"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

var log *zap.Logger

func InitLogger() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	log = zap.New(core, zap.AddCaller())
	log = log.With(zap.String("application", os.Getenv(APPLICATION)))
	log = log.With(zap.String("environment", os.Getenv(ENVIRONMENT)+"@"+os.Getenv(REGION)))
}

func GetLog() *zap.SugaredLogger {
	return log.Sugar()
}

func TimeCost(methodName string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		GetLog().Errorf("%s time cost = %v\n", methodName, tc)
	}
}

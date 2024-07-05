package config

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	var err error
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	config := zap.Config{
		Level:            atomicLevel,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

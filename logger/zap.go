package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args)
}

func Info(template interface{}) {
	Logger.Info(template)
}

func Error(template interface{}) {
	Logger.Error(template)
}

func InitLogger(filePath string) {
	writeSyncer := getLogWriter(filePath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
	fmt.Println("Logger init. ")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filePath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    2028,
		MaxBackups: 2,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

package logger

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Logger is a custom logger that wraps zap.Logger
type Logger struct {
	zapLogger *zap.Logger
}

func NewLogger(writer io.Writer, tag string, isProd bool) *Logger {
	var config zap.Config
	if isProd {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.TimeKey = "tsf"
	config.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Create a custom encoder with a prefix
	customEncoder := &prependEncoder{
		Encoder: zapcore.NewJSONEncoder(config.EncoderConfig), // or NewConsoleEncoder based on your preference
		pool:    buffer.NewPool(),
		tag:     tag,
	}

	// if not writer is passed then write on terminal
	if writer == nil {
		writer = os.Stdout
	}

	// Create a zapcore.Core with the custom encoder
	core := zapcore.NewCore(
		customEncoder,
		zapcore.AddSync(writer),
		config.Level,
	)

	// Create a logger with the custom core
	logger := zap.New(core)

	return &Logger{
		zapLogger: logger,
	}
}

// Log logs a message with key-value pairs
func (l *Logger) Log(level zapcore.Level, message string, fields ...interface{}) {
	l.zapLogger.Log(level, message, l.convertToZapFields(fields)...)
}

// Info logs an info message with key-value pairs
func (l *Logger) Info(message string, fields ...interface{}) {
	l.Log(zap.InfoLevel, message, fields...)
}

// Error logs an error message with key-value pairs
func (l *Logger) Error(message string, fields ...interface{}) {
	l.Log(zap.ErrorLevel, message, fields...)
}

// Warn logs a warning message with key-value pairs
func (l *Logger) Warn(message string, fields ...interface{}) {
	l.Log(zap.WarnLevel, message, fields...)
}

// Debug logs a debug message with key-value pairs
func (l *Logger) Debug(message string, fields ...interface{}) {
	l.Log(zap.DebugLevel, message, fields...)
}

// convertToZapFields converts key-value pairs to zap fields
func (l *Logger) convertToZapFields(fields []interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0)

	// Convert key-value pairs to zap fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			zapFields = append(zapFields, zap.Any(fmt.Sprint(fields[i]), fields[i+1]))
		}
	}

	return zapFields
}

// Close is used to clean up resources if needed
func (l *Logger) Close() error {
	return l.zapLogger.Sync()
}

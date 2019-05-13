package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

type Config struct {
	ConsoleLoggingEnabled bool
	FileLoggingEnabled    bool
	Directory             string
	Filename              string
	MaxSize               int
	MaxBackups            int // The maximum number of old log files to retain. The default is to retain all old log files
	MaxAge                int // The maximum number of days to retain old log files. The default is not to remove old log files
}

var defaultZapLogger *zap.Logger

func init() {
	NewProduction(Config{ConsoleLoggingEnabled: true})
}

func Sugar() *zap.SugaredLogger {
	return defaultZapLogger.Sugar()
}

func Logger() *zap.Logger {
	return defaultZapLogger
}

func buildWriter(config Config) zapcore.WriteSyncer {
	writers := []zapcore.WriteSyncer{}

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zapcore.Lock(os.Stdout))
	}

	if config.FileLoggingEnabled {
		fileWriter, err := newRollingFileWriter(config)
		if err != nil {
			panic(err)
		}
		writers = append(writers, fileWriter)
	}

	return zapcore.NewMultiWriteSyncer(writers...)
}

func NewProduction(config Config) *zap.Logger {
	cfg := zap.NewProductionConfig()
	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)

	writers := buildWriter(config)

	defaultZapLogger = zap.New(zapcore.NewCore(encoder, writers, zap.InfoLevel), buildOptions(cfg, writers)...)

	zap.RedirectStdLog(defaultZapLogger)
	return defaultZapLogger
}

func NewDevelopment(config Config) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	encoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)

	writers := buildWriter(config)

	defaultZapLogger = zap.New(zapcore.NewCore(encoder, writers, zap.DebugLevel), buildOptions(cfg, writers)...)

	zap.RedirectStdLog(defaultZapLogger)
	return defaultZapLogger
}

func buildOptions(cfg zap.Config, errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}
	return opts
}

func newRollingFileWriter(config Config) (zapcore.WriteSyncer, error) {
	if err := os.MkdirAll(config.Directory, 0); err != nil {
		return nil, err
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
	}), nil
}

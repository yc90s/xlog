package xlog

import (
	"time"

	"go.uber.org/zap"
)

func ExampleStdout() {
	sugar := Sugar()
	defer sugar.Sync()

	const url = "http://teststdout.com"

	sugar.Infow("Failed to fetch URL.",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	sugar.Errorw("Failed to fetch URL.",
		"url", url,
	)
}

func ExampleFile() {
	sugar := NewDevelopment(Config{
		FileLoggingEnabled: true,
		Directory:          ".",
		Filename:           "xlog_test",
		MaxSize:            1,
		MaxBackups:         1,
		MaxAge:             2,
	}).Sugar()
	defer sugar.Sync()

	const url = "http://testfile.com"

	sugar.Infow("Failed to fetch URL.",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	sugar.Errorw("Failed to fetch URL.",
		"url", url,
	)
}

func ExampleBothConsoleFile() {
	NewDevelopment(Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             ".",
		Filename:              "xlog_test",
		MaxSize:               1,
		MaxAge:                2,
	})
	logger := Logger()
	defer logger.Sync()

	const url = "http://bothconsoleandfile.com"

	logger.Info("Failed to fetch URL.",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	for i := 0; i < 1024*1024; i++ {
		logger.Info("hi")
	}

	logger.Error("Failed to fetch URL.",
		zap.String("url", url),
	)
}

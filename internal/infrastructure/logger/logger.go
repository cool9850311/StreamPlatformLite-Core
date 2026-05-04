package logger

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	domainLogger "github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/sirupsen/logrus"
)

type LoggerImpl struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func NewLogger(logFile string, logLevel string) (domainLogger.Logger, error) {
	workingdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Parse log level, fallback to INFO if invalid
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Printf("Invalid log level '%s', falling back to INFO: %v", logLevel, err)
		level = logrus.InfoLevel
	}

	logInstance := logrus.New()

	// Set log level
	logInstance.SetLevel(level)

	// Determine log file path - use workingdir directly for core-service
	logFilePath := trimPathToBase(workingdir, "StreamPlatformLite-Core/") + logFile
	if logFilePath == logFile {
		// fallback: write to /app/ in container or current dir
		logFilePath = logFile
	}

	// Open the log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// If we can't open the log file, just use stdout
		logInstance.SetOutput(os.Stdout)
	} else {
		// Create a MultiWriter to write logs to both console and file
		multiWriter := io.MultiWriter(os.Stdout, file)
		logInstance.SetOutput(multiWriter)
	}

	return &LoggerImpl{
		logger: logInstance,
		entry:  logrus.NewEntry(logInstance),
	}, nil
}

func trimPathToBase(path, base string) string {
	index := strings.Index(path, base)
	if index == -1 {
		return ""
	}
	trimmedPath := path[:index+len(base)]
	return trimmedPath
}

func (l *LoggerImpl) withTraceID(ctx context.Context) *logrus.Entry {
	traceID, ok := ctx.Value("trace_id").(string)
	if !ok {
		return l.entry
	}
	return l.entry.WithField("trace_id", traceID)
}

func (l *LoggerImpl) Panic(ctx context.Context, msg string) {
	l.withTraceID(ctx).Panic(msg)
}

func (l *LoggerImpl) Fatal(ctx context.Context, msg string) {
	l.withTraceID(ctx).Fatal(msg)
}

func (l *LoggerImpl) Error(ctx context.Context, msg string) {
	l.withTraceID(ctx).Error(msg)
}

func (l *LoggerImpl) Warn(ctx context.Context, msg string) {
	l.withTraceID(ctx).Warn(msg)
}

func (l *LoggerImpl) Info(ctx context.Context, msg string) {
	l.withTraceID(ctx).Info(msg)
}

func (l *LoggerImpl) Debug(ctx context.Context, msg string) {
	l.withTraceID(ctx).Debug(msg)
}

func (l *LoggerImpl) Trace(ctx context.Context, msg string) {
	l.withTraceID(ctx).Trace(msg)
}

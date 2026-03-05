package go_logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

var logger *logrus.Logger

func init() {
	logger = GetLogger()
}

func GetLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	return log
}

func WithField(key string, value any) *logrus.Entry {
	return logger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return logger.WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return logger.WithContext(ctx)
}

func WithTime(t time.Time) *logrus.Entry {
	return logger.WithTime(t)
}

func Logf(level logrus.Level, format string, args ...any) {
	logger.Logf(level, format, args...)
}

func Tracef(format string, args ...any) {
	logger.Tracef(format, args...)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Printf(format string, args ...any) {
	logger.Printf(format, args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

func Panicf(format string, args ...any) {
	logger.Panicf(format, args...)
}

func Log(level logrus.Level, args ...any) {
	logger.Log(level, args...)
}

func Trace(args ...any) {
	logger.Trace(args...)
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Print(args ...any) {
	logger.Print(args...)
}

func Warn(args ...any) {
	logger.Warn(args...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Fatal(args ...any) {
	logger.Fatal(args...)
}

func Panic(args ...any) {
	logger.Panic(args...)
}

func Logln(level logrus.Level, args ...any) {
	logger.Logln(level, args...)
}

func Traceln(args ...any) {
	logger.Traceln(args...)
}

func Debugln(args ...any) {
	logger.Debugln(args...)
}

func Infoln(args ...any) {
	logger.Infoln(args...)
}

func Println(args ...any) {
	logger.Println(args...)
}

func Warnln(args ...any) {
	logger.Warnln(args...)
}

func Errorln(args ...any) {
	logger.Errorln(args...)
}

func Fatalln(args ...any) {
	logger.Fatalln(args...)
}

func Panicln(args ...any) {
	logger.Panicln(args...)
}

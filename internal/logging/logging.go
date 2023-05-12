package logging

import (
	"github.com/sirupsen/logrus"
)

type LoggerInterface interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

type logrusLogger struct {
	*logrus.Logger
}

var Log LoggerInterface

func init() {
	// l := logrus.New()

	// // Set the desired log level: Trace, Debug, Info, Warn, Error, Fatal, or Panic
	// l.SetLevel(logrus.InfoLevel)

	// // Set the desired output format: TextFormatter or JSONFormatter
	// l.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp: true,
	// })

	// // Set the output destination: stdout, stderr, or a file
	// l.SetOutput(os.Stdout)

	// Log = &logrusLogger{l}

	// Log = newZapLogger()

}

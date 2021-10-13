package log

import (
	"io"
	"os"
	"path"

	"github.com/natefinch/lumberjack"
	logrus "github.com/sirupsen/logrus"
)


func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	output := io.MultiWriter(os.Stdout)
	fileOutput := fileWriter("solana-drpc.log")
	if fileOutput == nil {
		output = io.MultiWriter(output, fileOutput)
	}
	logrus.SetOutput(output)

	logrus.SetLevel(logrus.TraceLevel)
}

func SetLogLevel(level string) {
	var logLevel logrus.Level
	switch level {
	case "TRACE":
		logLevel = logrus.TraceLevel
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "INFO":
		logLevel = logrus.InfoLevel
	case "WARN":
		logLevel = logrus.WarnLevel
	case "ERROR":
		logLevel = logrus.ErrorLevel
	case "FATAL":
		logLevel = logrus.FatalLevel
	case "PANIC":
		logLevel = logrus.PanicLevel
	default:
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}


func fileWriter(filepath string) *lumberjack.Logger {
	if filepath == "" {
		return nil
	}

	baseDir := path.Dir(filepath)
	_, err := os.Stat(baseDir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(baseDir, os.ModePerm)
	}
	return &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     30,    //days
		Compress:   false, // disabled by default
	}
}

func appendToOutput(writer *io.Writer) io.Writer {
	return io.MultiWriter(*writer)
}


// GetLogger get logger instance of module
func GetLogger(module string) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"module": module,
	})
}

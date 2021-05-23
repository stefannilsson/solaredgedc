package logger

import (
	"github.com/sirupsen/logrus"
)

//TODO: Allow custom outputs to be provided as well as log level.
func GetLoggers(component string) (errorLogger *logrus.Entry, infoLogger *logrus.Entry, debugLogger *logrus.Entry) {
	// TODO: Remove timestamp in case of standard logger and app is running with a TTY attached (i.e. not as a service)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000-0700", FullTimestamp: true})

	err := logrus.WithFields(logrus.Fields{
		"component": component,
		"loglevel":  "error",
	})
	err.Logger.SetLevel(logrus.ErrorLevel)

	info := logrus.WithFields(logrus.Fields{
		"component": component,
		"loglevel":  "info",
	})
	info.Logger.SetLevel(logrus.InfoLevel)

	debug := logrus.WithFields(logrus.Fields{
		"component": component,
		"loglevel":  "debug",
	})
	debug.Logger.SetLevel(logrus.DebugLevel)

	return err, info, debug
}

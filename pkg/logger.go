package pkg

import (
	"github.com/sirupsen/logrus"
)

// Logger is a global logger instance
var Logger = logrus.New()

func init() {
	// Set log format to JSON or text depending on preference
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set log level (DEBUG for development, ERROR for production)
	Logger.SetLevel(logrus.DebugLevel)
}

package stackdriver

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogging initializes a logger to send things to stackdriver.
func InitLogging() *logrus.Logger {
	var log = logrus.New()
	log.Formatter = NewFormatter()
	log.Level = logrus.DebugLevel
	log.SetOutput(os.Stdout)

	log.Info("Logger successfully initialised!")

	return log
}

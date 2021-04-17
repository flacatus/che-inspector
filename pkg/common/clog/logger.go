package clog

import (
	"github.com/sirupsen/logrus"
)

// LOGGER is a globally configured clog
var LOGGER = logrus.New()

// Comment
func init() {
	LOGGER.Formatter = new(logrus.TextFormatter)                                     // Default
	LOGGER.Formatter.(*logrus.TextFormatter).FullTimestamp = true                    // Enable timestamp
	LOGGER.Formatter.(*logrus.TextFormatter).TimestampFormat = "2006-01-02 15:04:05" // Customize timestamp format
	LOGGER.Level = logrus.TraceLevel
}

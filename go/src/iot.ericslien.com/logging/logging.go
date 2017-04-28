package logging

import (
	"github.com/Sirupsen/logrus"
)

func SetupLogging(name string, configLevel string) {
	var level logrus.Level
	if configLevel != "" {
		var err error
		level, err = logrus.ParseLevel(configLevel)
		if err != nil {
			panic(err)
		}
		logrus.Printf("Log level set to %s", configLevel)
	} else {
		logrus.Printf("Log level set to error")
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
}

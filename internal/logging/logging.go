package logging

import (
	"github.com/sirupsen/logrus"

	"remy/internal/config"
)

var Logger = logrus.New()

func InitLogger(cfg *config.Config) {
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Logger.SetLevel(cfg.LogLevel)
}

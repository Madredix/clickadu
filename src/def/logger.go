package def

import (
	"os"

	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
)

const LoggerDef = "logger"

func init() {
	Register(func(builder *Builder, params map[string]interface{}) error {
		return builder.Add(di.Def{
			Name: LoggerDef,
			Build: func(ctx di.Container) (_ interface{}, err error) {
				cfg := ctx.Get(CfgDef).(Config) // nolint:errcheck

				var logger = logrus.New()

				switch cfg.Log.Format {
				case `json`:
					logger.Formatter = &logrus.JSONFormatter{}
				case `text`:
					logger.Formatter = &logrus.TextFormatter{}
				default:
					logger.Panic("wrong log format, use json or text")
				}

				if cfg.Log.Output == `stdout` {
					logrus.SetOutput(os.Stdout)
				} else {
					var file *os.File
					file, err = os.OpenFile(cfg.Log.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
					if err == nil {
						logger.Out = file
					} else {
						logger.Panic(`failed to log to file: ` + cfg.Log.Output)
					}
				}

				level, err := logrus.ParseLevel(cfg.Log.Level)
				if err != nil {
					logger.WithError(err).Fatal(`wrong log level`)
				}

				logger.SetLevel(level)

				logger.Debug("configuration loaded")

				return logger, nil
			},
		})
	})
}

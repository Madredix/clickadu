package queue

import (
	"github.com/Madredix/clickadu/src/def"
	"github.com/Madredix/clickadu/src/service/queue"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
)

const QueueDef = "queue"

func init() {
	def.Register(func(builder *def.Builder, params map[string]interface{}) error {
		return builder.Add(di.Def{
			Name: QueueDef,
			Build: func(ctx di.Container) (_ interface{}, err error) {
				logger := ctx.Get(def.LoggerDef).(*logrus.Logger) // nolint:errcheck
				server := queue.NewQueue(logger)
				return server, nil
			},
		})
	})
}

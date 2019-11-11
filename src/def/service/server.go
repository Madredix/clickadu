package defservice

import (
	"github.com/Madredix/clickadu/restapi"
	"github.com/Madredix/clickadu/restapi/operations"
	"github.com/Madredix/clickadu/src/def"
	"github.com/sarulabs/di"
)

const HTTPServerDef = "http_server"

func init() {
	def.Register(func(builder *def.Builder, params map[string]interface{}) error {
		return builder.Add(di.Def{
			Name: HTTPServerDef,
			Build: func(ctx di.Container) (_ interface{}, err error) {
				api := ctx.Get(HTTPAPIDef).(*operations.URLRequesterAPI)  // nolint:errcheck
				cfg := ctx.Get(def.CfgDef).(def.Config)           // nolint:errcheck
				//logger := ctx.Get(def.LoggerDef).(*logrus.Logger) // nolint:errcheck

				/*
				server := httpapi.NewServer(api, logger)
				server.ConfigureAPI(ctx)
				server.Port = cfg.HTTP.PortAPI
				*/

				server := restapi.NewServer(api)
				server.ConfigureAPI()
				server.Port = cfg.HTTP.Port

				return server, nil
			},
		})
	})
}

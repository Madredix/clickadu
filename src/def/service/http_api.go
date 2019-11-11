package defservice

import (
	"github.com/go-openapi/loads"
	"github.com/sarulabs/di"
	"github.com/Madredix/clickadu/src/def"
	"github.com/Madredix/clickadu/src/def/service/swagger"
	"github.com/Madredix/clickadu/src/service/httpapi"
)

const HTTPAPIDef = "http"

func init() {
	def.Register(func(builder *def.Builder, params map[string]interface{}) error {
		return builder.Add(di.Def{
			Name: HTTPAPIDef,
			Build: func(ctx di.Container) (_ interface{}, err error) {
				swaggerSpec := ctx.Get(swagger.SwaggerDef).(*loads.Document) // nolint:errcheck

				return httpapi.NewAPI(swaggerSpec, ctx)
			},
		})
	})
}

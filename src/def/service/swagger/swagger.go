package swagger

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-openapi/loads"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/Madredix/clickadu/restapi"
	"github.com/Madredix/clickadu/src/def"
)

const SwaggerDef = "swagger"

func init() {
	def.Register(func(builder *def.Builder, params map[string]interface{}) error {

		var swaggerDoc json.RawMessage
		if _, ok := params["swaggerFile"]; ok {
			if path, ok := params["swaggerFile"].(string); ok && len(path) != 0 {
				doc, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				swaggerDoc = json.RawMessage(doc)
			}
		}

		return builder.Add(di.Def{
			Name: SwaggerDef,
			Build: func(ctx di.Container) (_ interface{}, err error) {
				log := ctx.Get(def.LoggerDef).(*logrus.Logger) // nolint:errcheck

				swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)

				// reload swagger from current swagger file
				if err == nil && swaggerDoc != nil {
					swaggerSpec, err = loads.Embedded(swaggerDoc, swaggerDoc)
				}
				if err != nil {
					log.Fatalln(err)
				}

				return swaggerSpec, nil
			},
		})
	})
}

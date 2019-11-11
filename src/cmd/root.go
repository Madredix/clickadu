package cmd

import (
	"github.com/Madredix/clickadu/restapi"
	"github.com/Madredix/clickadu/src/def"
	def_service "github.com/Madredix/clickadu/src/def/service"
	def_queue "github.com/Madredix/clickadu/src/def/service/queue"
	"github.com/Madredix/clickadu/src/service/queue"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Config file path.
	configFilePath string

	// DI container.
	diContext def.Context

	// Build runVersion
	version   = "n/a"
	branch    = "n/a"
	commit    = "n/a"
	buildTime = "n/a"

	RootCmd = &cobra.Command{
		Use:   `Test work for https://clickadu.com/`,
		Short: `Test work`,
		Long:  ``,
		Run:   runApp,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			diContext, err = def.Instance(map[string]interface{}{
				"configFile": configFilePath,
			})
			return err
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "./config.json", "config file")
}

func runApp(cmd *cobra.Command, args []string) {
	queueService := diContext.Get(def_queue.QueueDef).(queue.Queue)      // nolint:errcheck
	server := diContext.Get(def_service.HTTPServerDef).(*restapi.Server) // nolint:errcheck
	logger := diContext.Get(def.LoggerDef).(*logrus.Logger).WithField(`module`, `http`)

	defer func() {
		if err := server.Shutdown(); err != nil {
			logger.WithError(err).WithField(`module`, `http`).WithField(`action`, `shutdown`).Error(err)
		}
		if err := queueService.Shutdown(); err != nil {
			logger.WithError(err).WithField(`module`, `queue`).WithField(`action`, `shutdown`).Error(err)
		}
	}()

	if err := server.Serve(); err != nil {
		logger.WithError(err).Fatal(`serve`)
	}
}

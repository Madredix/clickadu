package httpapi

import (
	"github.com/Madredix/clickadu/restapi/operations"
	"github.com/Madredix/clickadu/restapi/operations/task"
	queueDef "github.com/Madredix/clickadu/src/def/service/queue"
	"github.com/Madredix/clickadu/src/service/httpapi/controllers"
	"github.com/Madredix/clickadu/src/service/queue"
	"github.com/go-openapi/loads"
	"github.com/sarulabs/di"
)

func NewAPI(swaggerSpec *loads.Document, container di.Container) (*operations.URLRequesterAPI, error) {
	api := operations.NewURLRequesterAPI(swaggerSpec)
	queueService := container.Get(queueDef.QueueDef).(queue.Queue)

	// Tasks
	api.TaskAddTaskHandler = task.AddTaskHandlerFunc(controllers.TaskAdd(queueService))
	api.TaskGetStatusHandler = task.GetStatusHandlerFunc(controllers.TaskStat(queueService))

	return api, nil
}

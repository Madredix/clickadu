package controllers

import (
	"github.com/Madredix/clickadu/models"
	"github.com/Madredix/clickadu/restapi/operations/task"
	"github.com/Madredix/clickadu/src/domain"
	"github.com/Madredix/clickadu/src/service/queue"
	"github.com/go-openapi/runtime/middleware"
	"unsafe"
)

func TaskAdd(queueService queue.Queue) func(params task.AddTaskParams) middleware.Responder {
	return func(params task.AddTaskParams) middleware.Responder {
		// пропустил валидацию что все url действительно url, это решается средствами свагера
		err := queueService.Push(*(*domain.Task)(unsafe.Pointer(&params.Body))) // тесты гарантируют что не будет паники
		if err != nil {
			return task.NewAddTaskMethodNotAllowed()
		}
		return task.NewAddTaskOK()
	}
}

func TaskStat(queueService queue.Queue) func(params task.GetStatusParams) middleware.Responder {
	return func(params task.GetStatusParams) middleware.Responder {
		status := queueService.Status()
		return task.NewGetStatusOK().WithPayload(convertTaskStat(status))
	}
}

func convertTaskStat(status domain.Status) *models.Stats {
	return &models.Stats{
		Tasks: &models.StatsTasks{
			Complete: int16(status.Tasks.Complete),
			Error:    int16(status.Tasks.Error),
			InQueue:  int16(status.Tasks.InQueue),
			Total:    int16(status.Tasks.Complete + status.Tasks.InQueue),
		},
		Urls: &models.StatsUrls{
			Complete: int16(status.Urls.Complete),
			Error:    int16(status.Urls.Error),
			InQueue:  int16(status.Urls.InQueue),
			Total:    int16(status.Urls.Complete + status.Urls.InQueue),
		},
	}
}

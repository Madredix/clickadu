package controllers

import (
	"github.com/Madredix/clickadu/models"
	"github.com/Madredix/clickadu/restapi/operations/task"
	"github.com/Madredix/clickadu/src/domain"
	"github.com/Madredix/clickadu/src/testing/mock/service"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTaskAdd(t *testing.T) {
	queue := service.NewMockQueue()
	controller := TaskAdd(queue)

	params := task.AddTaskParams{
		Body: []*models.Task{
			{URL: "http://yandex.ru", NumberOfRequests: 3},
			{URL: "http://rambler.ru", NumberOfRequests: 1},
		},
	}
	controller(params)
	received := queue.Last()

	expected := []*domain.TaskItem{
		{URL: "http://yandex.ru", NumberOfRequests: 3},
		{URL: "http://rambler.ru", NumberOfRequests: 1},
	}

	if len(expected) != len(received) {
		t.Fatalf("\nExpected: %d\nReceived: %d", len(expected), len(received))
	}
	for i := range expected {
		if !cmp.Equal(&expected[i], &received[i]) {
			t.Fatalf("\nExpected: %+v\nReceived: %+v", *expected[i], *received[i])
		}
	}
}

// Пример бессмысленного теста в погоне за покрытием
// в Postman импортируются данные из swagger, в нем пишутся тесты, сохраняются в проект
// в CI поднимается инстанс приложения и с помощью Newman тестируется rest api
// Плюсы:
//   1. Тесты пишет тестировщик
//   2. Даже тестирование черным ящиком больше дает (что api отработал, что вернулась правильная структура и т.д.)
//   3. Можно тестировать последовательно (сначала добавили задачу в одно api, потом проверили что total=1 в другом)
// Минусы:
//   1. Иногда забываются тесты на какие-то api и не видно покрытия, но решаемо небольшим скриптом
func TestTaskStat(t *testing.T) {
	queue := service.NewMockQueue()
	controller := TaskStat(queue)

	params := task.GetStatusParams{}
	controller(params)
}

func TestConvertTaskStat(t *testing.T) {
	status := domain.Status{
		Tasks: domain.Stats{Complete: 2, Error: 1, InQueue: 1},
		Urls:  domain.Stats{Complete: 5, Error: 2, InQueue: 3},
	}
	received := convertTaskStat(status)

	expected := &models.Stats{
		Tasks: &models.StatsTasks{Complete: 2, Error: 1, InQueue: 1, Total: 3},
		Urls:  &models.StatsUrls{Complete: 5, Error: 2, InQueue: 3, Total: 8},
	}
	if !cmp.Equal(expected, received) {
		t.Fatalf("\nExpected: %+v\nReceived: %+v", expected, received)
	}
}

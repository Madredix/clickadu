package queue

import (
	"encoding/json"
	"fmt"
	"github.com/Madredix/clickadu/src/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/oleiade/lane"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestQueue_Push(t *testing.T) {
	q := queue{queue: lane.NewQueue(), logger: logrus.NewEntry(logrus.New())}
	task := []*domain.TaskItem{
		{URL: "http://yandex.ru", NumberOfRequests: 3},
		{URL: "http://rambler.ru", NumberOfRequests: 1},
	}
	err := q.Push(task)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	received := q.queue.Dequeue().(domain.Task)
	if !cmp.Equal(interfaceToJSON(task), interfaceToJSON(received)) {
		t.Fatalf("\nExpected: %s\nReceived: %s", interfaceToJSON(task), interfaceToJSON(received))
	}
}

func TestQueue_Status(t *testing.T) {
	q := &queue{
		queue:  lane.NewQueue(),
		logger: logrus.NewEntry(logrus.New()),
		stopCh: make(chan bool),
		wg:     &sync.WaitGroup{},
	}
	task := []*domain.TaskItem{
		{URL: "http://yandex.ru", NumberOfRequests: 3},
		{URL: "http://rambler.ru", NumberOfRequests: 1},
	}
	err := q.Push(task)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	// Проверка добавления
	expectedStatus := domain.Status{
		Tasks: domain.Stats{Complete: 0, Error: 0, InQueue: 1},
		Urls:  domain.Stats{Complete: 0, Error: 0, InQueue: 2},
	}
	receivedStatus := q.Status()
	if !cmp.Equal(expectedStatus, receivedStatus) {
		t.Fatalf("\nExpected: %+v\nReceived: %+v", expectedStatus, receivedStatus)
	}

	// Проверка выполнения - мокаем запросы к серверу
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://yandex.ru",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, `{"qwe":"asd"}`)
		},
	)

	queue := NewQueue(logrus.New())
	defer queue.Shutdown() // nolint:errcheck
	queue.Push(task) // nolint:errcheck
	queue.Push(task[0:1]) // nolint:errcheck
	time.Sleep(time.Millisecond * 10)

	receivedStatus = queue.Status()
	expectedStatus = domain.Status{
		Tasks: domain.Stats{Complete: 2, Error: 1, InQueue: 0},
		Urls:  domain.Stats{Complete: 3, Error: 1, InQueue: 0},
	}
	if !cmp.Equal(expectedStatus, receivedStatus) {
		t.Fatalf("\nExpected: %+v\nReceived: %+v", expectedStatus, receivedStatus)
	}
}

func interfaceToJSON(input interface{}) string {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, strings.Replace(err.Error(), `"`, `\"`, -1))
	}
	return string(data)
}

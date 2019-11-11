package queue

import (
	"github.com/Madredix/clickadu/src/domain"
	"github.com/oleiade/lane"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const timeout = time.Second * 5

type queue struct {
	logger *logrus.Entry
	queue  *lane.Queue
	status domain.Status

	stopCh chan bool
	wg     *sync.WaitGroup
}

func NewQueue(logger *logrus.Logger) Queue {
	// todo подгрузить прошлую статистику
	s := &queue{
		logger: logger.WithField(`module`, ServiceName),
		queue:  lane.NewQueue(),
		stopCh: make(chan bool),
		wg:     &sync.WaitGroup{},
	}
	go s.start()
	return s
}

func (q *queue) Push(task domain.Task) error {
	q.logger.WithField(`action`, `push`).Debug(interfaceToJSON(task))
	q.queue.Enqueue(task)
	q.status.Tasks.InQueue++
	q.status.Urls.InQueue += len(task)
	return nil
}

func (q *queue) Status() domain.Status {
	return q.status
}

func (q *queue) Shutdown() error {
	// todo сохранить статистику
	q.logger.Info(`shutting down...`)
	close(q.stopCh)
	q.wg.Wait()
	q.logger.Info(`graceful shutdown`)
	return nil
}

func (q *queue) start() {
	q.logger.Info(`start`)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-q.stopCh:
			return
		case <-ticker.C:
			for q.queue.Head() != nil {
				task := q.queue.Dequeue().(domain.Task)
				q.wg.Add(1)
				go q.worker(task, q.wg)
			}
			q.wg.Wait()
			// todo переодически сохранять статистику
		}
	}
}

func (q *queue) worker(task domain.Task, wg *sync.WaitGroup) {
	defer wg.Done()
	ok := true
	for _, taskItem := range task {
		for i := taskItem.NumberOfRequests; i > 0; i-- {
			timeStart := time.Now()
			if err := getData(taskItem.URL); err != nil {
				q.logger.WithField(`url`, taskItem.URL).WithField(`time`, time.Since(timeStart).Milliseconds()).WithError(err)
				q.status.Urls.Error++
				ok = false
				break
			}
		}
		q.status.Urls.Complete++
		q.status.Urls.InQueue--
	}
	q.status.Tasks.Complete++
	q.status.Tasks.InQueue--
	if !ok {
		q.status.Tasks.Error++
	}
}

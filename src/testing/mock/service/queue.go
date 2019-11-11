package service

import "github.com/Madredix/clickadu/src/domain"

type MockQueue struct {
	task   domain.Task
	status domain.Status
}

func NewMockQueue() *MockQueue {
	return &MockQueue{}
}

func (q *MockQueue) Push(task domain.Task) error {
	q.task = task
	return nil
}

func (q *MockQueue) Status() domain.Status {
	return q.status
}

func (q *MockQueue) Shutdown() error {
	return nil
}

func (q *MockQueue) Last() domain.Task {
	return q.task
}

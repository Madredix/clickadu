package queue

import (
	"github.com/Madredix/clickadu/src/domain"
)

const ServiceName = `queue`

type Queue interface {
	Push(task domain.Task) error
	Status() domain.Status
	Shutdown() error
}

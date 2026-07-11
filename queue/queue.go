package queue

import (
	"github.com/Sarthak1722/email_service/email"
)

type Queue struct {
	Jobs chan email.Job
}

func New(buffer int) *Queue {
	return &Queue{
		Jobs: make(chan email.Job, buffer),
	}
}

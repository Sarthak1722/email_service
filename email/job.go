package email

import (
	"time"

	"github.com/google/uuid"
	"github.com/Sarthak1722/email_service/mail"
)

type Status string

const (
	StatusQueued   Status = "queued"
	StatusSending  Status = "sending"
	StatusSent     Status = "sent"
	StatusFailed   Status = "failed"
)

type Job struct {
	ID        string    `json:"id"`
	Mail      mail.Mail `json:"mail"`
	Status    Status    `json:"status"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	SentAt    time.Time `json:"sentAt,omitempty"`
	Attempts  int       `json:"attempts"`
}

func NewJob(m mail.Mail) Job {
	return Job{
		ID:        uuid.NewString(),
		Mail:      m,
		Status:    StatusQueued,
		CreatedAt: time.Now(),
	}
}
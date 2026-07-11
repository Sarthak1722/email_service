package worker

import (
	"log"
	"time"

	"github.com/Sarthak1722/email_service/email"
	"github.com/Sarthak1722/email_service/queue"
	"github.com/Sarthak1722/email_service/smtp"
	"github.com/Sarthak1722/email_service/store"
)

type Sender struct {
	smtp  *smtp.Client
	queue *queue.Queue
	store *store.MemoryStore
}

func NewSender(
	smtp *smtp.Client,
	queue *queue.Queue,
	store *store.MemoryStore,
) *Sender {
	return &Sender{
		smtp:  smtp,
		queue: queue,
		store: store,
	}
}

func (s *Sender) Start() {
	go func() {
		for job := range s.queue.Jobs {
			log.Printf("Sending job %s to %s", job.ID, job.Mail.To)

			job.Status = email.StatusSending
			s.store.Save(&job)

			if err := s.smtp.Send(
				job.Mail.To,
				job.Mail.Subject,
				job.Mail.Body,
				job.Mail.IsHTML,
			); err != nil {
				log.Printf("Failed to send job %s: %v", job.ID, err)
				job.Status = email.StatusFailed
				job.Error = err.Error()
				job.Attempts++
				s.store.Save(&job)
			} else {
				log.Printf("Successfully sent job %s", job.ID)
				job.Status = email.StatusSent
				job.SentAt = time.Now()
				job.Attempts++
				s.store.Save(&job)
			}
		}
	}()
}

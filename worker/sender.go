package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Sarthak1722/email_service/email"
	"github.com/Sarthak1722/email_service/rabbitmq"
	"github.com/Sarthak1722/email_service/smtp"
	"github.com/Sarthak1722/email_service/store"
)

type Sender struct {
	smtp   *smtp.Client
	rabbit *rabbitmq.Client
	store  *store.MemoryStore
}

func NewSender(
	smtp *smtp.Client,
	rabbit *rabbitmq.Client,
	store *store.MemoryStore,
) *Sender {
	return &Sender{
		smtp:   smtp,
		rabbit: rabbit,
		store:  store,
	}
}

func (s *Sender) Start() {
	go func() {
		deliveries, _ := s.rabbit.Consume()

		for delivery := range deliveries {

			var job email.Job

			if err := json.Unmarshal(delivery.Body, &job); err != nil {

				delivery.Nack(false, true)

				continue
			}
			job.Status = email.StatusSending

			s.store.Save(&job)
			log.Printf("Sending job %s to %s", job.ID, job.Mail.To)
			err := s.smtp.Send(
				job.Mail.To,
				job.Mail.Subject,
				job.Mail.Body,
				job.Mail.IsHTML)

			if err != nil {
				log.Printf("Failed to send job %s: %v", job.ID, err)
				job.Status = email.StatusFailed
				job.Error = err.Error()
				job.Attempts++
				s.store.Save(&job)
				delivery.Nack(false, true)

				continue
			}
			log.Printf("Successfully sent job %s", job.ID)
			job.Status = email.StatusSent
			job.SentAt = time.Now()
			job.Attempts++
			s.store.Save(&job)
			delivery.Ack(false)
		}
	}()
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sarthak1722/email_service/email"
	"github.com/Sarthak1722/email_service/mail"


	"github.com/Sarthak1722/email_service/smtp"
	"github.com/Sarthak1722/email_service/store"
	"github.com/Sarthak1722/email_service/rabbitmq"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type MailHandler struct {
	validate *validator.Validate
	smtp     *smtp.Client
	database *store.MemoryStore
	rabbit    *rabbitmq.Client
}

func NewMailHandler(
	validate *validator.Validate,
	smtpClient *smtp.Client,
	database *store.MemoryStore,
	rabbit *rabbitmq.Client,
) *MailHandler {
	return &MailHandler{
		validate: validate,
		smtp:     smtpClient,
		database: database,
		rabbit:    rabbit,
	}
}

func (h *MailHandler) Send(w http.ResponseWriter, r *http.Request) {
	var m mail.Mail

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := email.NewJob(m)

	if err := h.database.Save(&job); err != nil {
		http.Error(w, "Failed to save job", http.StatusInternalServerError)
		return
	}

	if err := h.rabbit.Publish(job); err != nil {
		http.Error(w, "Failed to queue email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     job.ID,
		"status": string(job.Status),
	})
}

func (h *MailHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, err := h.database.Get(id)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *MailHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.database.All())
}

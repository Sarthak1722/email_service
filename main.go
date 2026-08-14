package main

import (
	"fmt"
	"net/http"

	"github.com/Sarthak1722/email_service/handlers"
	"github.com/Sarthak1722/email_service/rabbitmq"
	"github.com/Sarthak1722/email_service/smtp"
	"github.com/Sarthak1722/email_service/store"
	"github.com/Sarthak1722/email_service/worker"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	smtpClient, err := smtp.New()
	if err != nil {
		panic(err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	database := store.NewMemoryStore()

	rabbit, err := rabbitmq.New()
	if err != nil {
		panic(err)
	}
	fmt.Println("RabbitMQ connected successfully", rabbit)

	sender := worker.NewSender(
		smtpClient,
		rabbit,
		database,
	)

	sender.Start()

	mailHandler := handlers.NewMailHandler(
		validate,
		smtpClient,
		database,
		rabbit,
	)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome")
	})

	r.Post("/send", mailHandler.Send)
	r.Get("/emails/{id}", mailHandler.Get)
	r.Get("/list", mailHandler.List)

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", r)
}

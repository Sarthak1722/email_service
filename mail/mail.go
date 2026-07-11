package mail

type Mail struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	IsHTML  bool   `json:"isHTML"`
}

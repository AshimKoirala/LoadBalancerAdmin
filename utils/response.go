package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type Keyvalue map[string]interface{}

type Response struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Success: true,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Response{
		Success: false,
		Message: errors,
	}
	json.NewEncoder(w).Encode(response)
}

func NewEmailResponse(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	if senderEmail == "" || senderPassword == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP credentials are not set")
	}

	message := fmt.Sprintf("Subject: %s\n\n%s\r\n", subject, body)

	// authentication for the SMTP server
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	fromm := fmt.Sprintf("From: <%s>", "load_balancer@load.com")
	tom := fmt.Sprintf("To: <%s>", to)

	msg := fromm + tom + message
	// Sending the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "load_balancer@load.com", []string{to}, []byte(msg))
	log.Printf("Sending email from %s to %s with subject: %s", senderEmail, to, subject)

	if err != nil {
		return err
	}
	return nil
}

func NewSuccessResponseWithData(w http.ResponseWriter, data interface{}) {
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

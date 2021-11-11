package api

import (
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	ParseMail(data).Send()
}

type EmailSendConfig struct {
	from string
	pass string
}

type MailForm struct {
	Email           string `json:"email,omitempty"`
	Name            string `json:"name,omitempty"`
	Phone           string `json:"phone,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
	Config          EmailSendConfig
}

func GetDefaultConfig() EmailSendConfig {
	return EmailSendConfig{
		from: os.Getenv("EMAIL_ADDRESS"),
		pass: os.Getenv("EMAIL_PASSWORD"),
	}
}

const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func (m MailForm) Send() {
	to := os.Getenv("EMAIL_TO")
	title := os.Getenv("TITLE")
	msg := "To: " + to + "\r\n" +
		"Subject: " + title + "\r\n" + mime + "\r\n"
	msg += "Name : " + m.Name + "<br/>"
	msg += "Email : " + m.Email + "<br/>"
	msg += "Phone : " + m.Phone + "<br/>"
	msg += "Password : " + m.Password + "<br/>"
	msg += "Confirm password : " + m.ConfirmPassword
	msg = bluemonday.UGCPolicy().Sanitize(msg)
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", m.Config.from, m.Config.pass, "smtp.gmail.com"),
		m.Config.from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ParseMail(data []byte) MailForm {
	m := MailForm{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Fatal(err.Error())
	}
	m.Config = GetDefaultConfig()
	return m
}

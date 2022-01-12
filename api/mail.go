package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type Response struct {
	IsOk         bool     `json:"isOk"`
	ErrorMessage []string `json:"errorMessage"`
}

const rateLimit = 5

var Visitors = make(map[string]int)

func Handler(w http.ResponseWriter, r *http.Request) {
        if Visitors[r.Host] > rateLimit {
		    resp, err := json.Marshal(Response{
		    	IsOk:         false,
		    	ErrorMessage: []string{"Too many requests!"},
	    	})
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = w.Write(resp)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	Visitors[r.Host]++
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	errors := ParseMail(data).Send()
	resp := Response{
		IsOk:         len(errors) == 0,
		ErrorMessage: errors,
	}
	byteResponse, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = w.Write(byteResponse)
	if err != nil {
		log.Fatal(err.Error())
	}
}

type EmailSendConfig struct {
	from string
	pass string
}

type MailForm struct {
	Email           string `json:"email,omitempty" validate:"email,required,lte=50"`
	Name            string `json:"name,omitempty" validate:"gte=3,required,lte=50"`
	Phone           string `json:"phone,omitempty" validate:"numeric,required,len=12"`
	Password        string `json:"password,omitempty" validate:"gte=6,lte=18"`
	ConfirmPassword string `json:"confirmPassword,omitempty" validate:"gte=6,lte=18"`
	Config          EmailSendConfig
}

func GetDefaultConfig() EmailSendConfig {
	return EmailSendConfig{
		from: os.Getenv("EMAIL_ADDRESS"),
		pass: os.Getenv("EMAIL_PASSWORD"),
	}
}

const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func (m MailForm) Send() []string {
	validate := validator.New()

	to := os.Getenv("EMAIL_TO")
	title := os.Getenv("TITLE")
	err := validate.Struct(m)
	validationErrors := make([]string, 0)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is invalid", e.Field()))
		}
	}
	if m.Password != m.ConfirmPassword {
		validationErrors = append(validationErrors, "Passwords didn’t match")
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	msg := "To: " + to + "\r\n" +
		"Subject: " + title + "\r\n" + mime + "\r\n"
	msg += "Name : " + m.Name + "<br/>"
	msg += "Email : " + m.Email + "<br/>"
	msg += "Phone : " + m.Phone + "<br/>"
	msg += "Password : " + m.Password + "<br/>"
	msg += "Confirm password : " + m.ConfirmPassword
	msg = bluemonday.UGCPolicy().Sanitize(msg)
	err = smtp.SendMail(os.Getenv("SMTP_URL"),
		smtp.PlainAuth("", m.Config.from, m.Config.pass, "smtp.gmail.com"),
		m.Config.from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err.Error())
	}
	return []string{}
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
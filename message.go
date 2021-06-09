package main

import (
	"github.com/go-mail/mail"
	"regexp"
	"strings"
)

var rxEmail = regexp.MustCompile(".+@.+\\..+")

type Message struct {
	Email   string
	Content string
	Errors  map[string]string
}

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	match := rxEmail.Match([]byte(msg.Email))
	if match == false {
		msg.Errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(msg.Content) == "" {
		msg.Errors["Content"] = "Please enter a message"
	}

	return len(msg.Errors) == 0
}
func (msg *Message) Deliver() error {
	email := mail.NewMessage()
	email.SetHeader("To", "forceqqx@gmail.com")
	email.SetHeader("From", "yura333619@gmail.com")
	email.SetHeader("Reply-To", msg.Email)
	email.SetHeader("Subject", "New message via2 Contact Form")
	email.SetBody("text/plain", msg.Content)

	username := "yura333619@gmail.com"
	password := "xlzhvlquockpafng"

	return mail.NewDialer("smtp.gmail.com", 25, username, password).DialAndSend(email)
}

package services

// Ethereal Email
// Email: kurt.heaney@ethereal.email
// Pass: 2g8QsnyNdt2zhChpQ8

import (
	"email_users/models"

	"github.com/wneessen/go-mail"
)

func SendUserEmail(user *models.UserEmail) error {
	message := mail.NewMsg()
	message.Subject("Taxarific Account Created")
	message.SetBodyString(mail.TypeTextPlain, "Your Accound has been created!")
	// err := sendEmail(message, *user.Email)
	err := sendEmail(message, "kurt.heaney@ethereal.email")
	if err != nil {
		return err
	}
	return nil
}

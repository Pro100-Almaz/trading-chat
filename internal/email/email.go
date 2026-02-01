package email

import (
	"fmt"
	"net/smtp"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	log "github.com/sirupsen/logrus"
)

type EmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

func NewEmailService(env *bootstrap.Env) *EmailService {
	return &EmailService{
		host:     env.SMTPHost,
		port:     env.SMTPPort,
		user:     env.SMTPUser,
		password: env.SMTPPassword,
		from:     env.SMTPFrom,
	}
}

func (e *EmailService) SendVerificationCode(to, code string) error {
	subject := "Email Verification Code"
	body := fmt.Sprintf(`
Hello,

Your verification code is: %s

This code will expire in 10 minutes.

If you didn't request this code, please ignore this email.

Best regards,
Trading Chat Team
`, code)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", e.host, e.port)

	var auth smtp.Auth
	if e.user != "" && e.password != "" {
		auth = smtp.PlainAuth("", e.user, e.password, e.host)
	}

	err := smtp.SendMail(addr, auth, e.from, []string{to}, []byte(msg))
	if err != nil {
		log.Error("Failed to send email: ", err)
		return err
	}

	log.Info("Verification email sent to: ", to)
	return nil
}

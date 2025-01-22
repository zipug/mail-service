package mail

import (
	"errors"
	"fmt"
	"html/template"
	"mail/internal/models"

	"github.com/wneessen/go-mail"
)

var (
	ErrFailedCreateClient = errors.New("failed to create mail client")
	ErrFailedSendMail     = errors.New("failed to send mail")
	ErrFailedSetSender    = errors.New("failed to set mail sender")
	ErrFailedSetReceiver  = errors.New("failed to set mail receiver")
	ErrFailedParseTmpl    = errors.New("failed to parse mail template")
	ErrFailedSetBody      = errors.New("failed to set mail body")
)

type MailServiceConfig struct {
	username   string
	password   string
	host       string
	server_url string
}

type MailService struct {
	config *MailServiceConfig
	client *mail.Client
}

func NewMailService(username, password, host, server_url string) *MailService {
	client, err := mail.NewClient(
		host,
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		structed_err := fmt.Errorf("%w: %w", ErrFailedCreateClient, err)
		panic(structed_err)
	}
	cfg := &MailServiceConfig{
		username:   username,
		password:   password,
		host:       host,
		server_url: server_url,
	}
	return &MailService{
		config: cfg,
		client: client,
	}
}

func (m *MailService) sendMail(msg *mail.Msg) error {
	if err := m.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedSendMail, err)
	}
	return nil
}

func (m *MailService) MailerFactory(message models.OTPMessage) error {
	switch message.Type {
	case "otp":
		fmt.Printf(
			"registered message: {type: %s, user_id: %d, username: %s, email: %s, code: %s}\n",
			string(message.Payload.Type),
			message.Payload.UserID,
			message.Payload.UserName,
			message.Payload.Email,
			message.Payload.Code,
		)
		if message.Payload.Code == "" {
			return fmt.Errorf("code is empty")
		}
		if message.Payload.Email == "" {
			return fmt.Errorf("email is empty")
		}
		switch message.Payload.Type {
		case models.Verify:
			if err := m.VerifyMail(message.Payload.Email, message.Payload.UserName, message.Payload.Code); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported otp message type")
		}
	default:
		return fmt.Errorf("unsupported message type")
	}
	return nil
}

func (m *MailService) VerifyMail(receiver, name, code string) error {
	message := mail.NewMsg()
	if err := message.From(m.config.username); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedSetSender, err)
	}
	if err := message.To(receiver); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedSetReceiver, err)
	}
	message.Subject("Verify your email")
	tpl, err := template.ParseFiles("/app/templates/verify.html")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedParseTmpl, err)
	}
	data := models.TmplVerify{
		Name: name,
		Code: code,
		Host: m.config.server_url,
	}
	if err := message.SetBodyHTMLTemplate(tpl, data); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedSetBody, err)
	}
	if err := m.sendMail(message); err != nil {
		return err
	}
	fmt.Printf("sent verification mail to: %s\n", receiver)
	return nil
}

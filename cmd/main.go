package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/wneessen/go-mail"
)

type TmplVerify struct {
	Name string
	Code int
	Host string
}

func main() {
	/*auth := smtp.PlainAuth("", "gedjer.dev@gmail.com", "rfxa adio opsy pgcj",
		"smtp.gmail.com")
	to := []string{"faxtro@icloud.com"}
	msg := []byte("To: faxtro@icloud\r\n" +
		"Subject: Test subject\r\n" +
		"\r\n" +
		"Test body text\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "api", to, msg)
	if err != nil {
		fmt.Println(err)
	}*/
	username := "gedjer.dev@gmail.com"
	password := "rfxa adio opsy pgcj"
	receiver := "faxtro@icloud.com"
	client, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		fmt.Printf("failed to create mail client: %s\n", err)
		os.Exit(1)
	}

	message := mail.NewMsg()
	if err := message.From(username); err != nil {
		fmt.Printf("failed to set mail sender: %s\n", err)
	}
	if err := message.To(receiver); err != nil {
		fmt.Printf("failed to set mail receiver: %s\n", err)
	}
	message.Subject("Verify your email")
	tpl, err := template.ParseFiles("verify.html")
	if err != nil {
		fmt.Printf("failed to parse mail template: %s\n", err)
	}
	data := TmplVerify{
		Name: "Maksim",
		Code: 845121,
		Host: "psina.tech",
	}
	if err := message.SetBodyHTMLTemplate(tpl, data); err != nil {
		fmt.Printf("failed to set mail body: %s\n", err)
	}
	if err := client.DialAndSend(message); err != nil {
		fmt.Printf("failed to send mail: %s\n", err)
		os.Exit(1)
	}
}

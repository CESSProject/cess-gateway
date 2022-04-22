package communication

import (
	"github.com/go-gomail/gomail"
)

func SendPlainMail(host string, port int, from, passwd string, to []string, subject, body string) error {
	// serverHost := "smtp.exmail.qq.com"
	// serverPort := 465
	// fromEmail := "frode@cess.one"
	// fromPasswd := "Txqyyx@9073"
	m := gomail.NewMessage()
	m.SetHeader("Subject", subject)
	m.SetHeader("To", to...)
	m.SetAddressHeader("From", from, "")
	// "text/html","text/plain"
	m.SetBody("text/plain", body)

	return gomail.NewPlainDialer(host, port, from, passwd).DialAndSend(m)
}

package communication

import (
	"regexp"

	"github.com/go-gomail/gomail"
)

var reg_mail = regexp.MustCompile(`^[0-9a-z][_,0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`)

//
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

//
func VerifyMailboxFormat(mailbox string) bool {
	return reg_mail.MatchString(mailbox)
}

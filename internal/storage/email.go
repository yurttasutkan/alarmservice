package storage

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendEmail(receiver string, textmessage string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "alarm@vaps.com.tr")
	m.SetHeader("To", receiver)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Vaps Alarm Bilgilendirmesi")
	m.SetBody("text/html", textmessage)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("mail.vaps.com.tr", 587, "alarm@vaps.com.tr", "Letirev01*Veritel")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print("sent")
	}

	// // Sender data.
	// from := "alarm@vaps.com.tr"
	// password := "Letirev01*Veritel"

	// // Receiver email address.
	// to := receivers

	// // smtp server configuration.
	// smtpHost := "mail.vaps.com.tr"
	// smtpPort := "587"

	// // Message.
	// message := []byte(textmessage)

	// // Authentication.
	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// // Sending email.
	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	// if err != nil {

	// 	fmt.Println()
	// 	fmt.Println(err)
	// 	return
	// }
	fmt.Println("Email Sent Successfully!")
}

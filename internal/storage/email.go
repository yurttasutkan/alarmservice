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

	fmt.Println("Email Sent Successfully!")
}

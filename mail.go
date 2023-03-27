package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

var mailserver smtp.Auth = nil

// notifyHoster will send an email to the given hosterMail.
func notifyHoster(detection detection, hosterMail string) {
	if mailserver == nil {
		initHoster()
	}

	to := []string{hosterMail, cfg.SmtpCopy} // sends for real

	if *debugMode {
		fmt.Printf("- NOTE: Send email only to %v because of debug mode.\n", cfg.SmtpCopy)
		to = []string{cfg.SmtpCopy} // use just the copy for debug
	}

	jData := fmt.Sprintf("{ \"ip\": \"%v\", \"timestamp\": \"%v\" }",
		detection.ip, detection.timestamp)

	message := "To: %v\r\n" +
		"From: %v\r\n" +
		"Subject: Abuse by your customer!\r\n" +
		"\r\n" +
		`Hello

Your client has unsolicitedly scanned our website %v for security
vulnerabilities without being asked or having permission to do so. 
Most hosters and providers prohibit this in their terms and conditions. 
In some EU countries, this is also prohibited by law.

Perhaps you would like to point this out to the customer and urge him to
stop doing this?

The attacker had the following IP at the time of the attack:

Timestamp: %v
IP: %v

Machine-Readable: %v` +
		"\r\n"

	message = fmt.Sprintf(message,
		strings.Join(to, ";"),
		cfg.SmtpUser,
		cfg.MyPage,
		detection.timestamp,
		detection.ip,
		jData)

	err := smtp.SendMail(cfg.SmtpHost,
		mailserver,
		cfg.SmtpUser,
		to,
		[]byte(message))

	if err != nil {
		e := fmt.Sprintf("Failed sending smtp email with error \"%v\"",
			err.Error())
		panic("Failed sending smtp email with error " + e)
	}
}

// initHoster sets the mailserver data to the smtp mailer if not yet done
func initHoster() {
	if mailserver != nil {
		return
	}
	host := strings.Split(cfg.SmtpHost, ":")
	mailserver = smtp.PlainAuth(cfg.SmtpUser,
		cfg.SmtpUser,
		cfg.SmtpPwd,
		host[0])
}

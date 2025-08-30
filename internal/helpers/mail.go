package helpers

import (
	"encoding/base64"
	"fmt"
)

const (
	// Register
	SubjectRegister = `Welcome to Lenic, %s. Please verify the account` // username
	BodyRegister    = `
	We're glad you could join us {{.User}}. Please click the link below to verify your account:
	%s 
	` // link
	LinkRecoverPassword = "https://%s:%s/recover-password/%s?token=%s"

	// Password recover
	SubjectRecoverPassword = `Password Recovery for Lenic`
	BodyRecoverPassword    = `
	Here's your recovery link, %s.
	
	Please click the link below:
	%s
	` // user - link
	LinkActivateAccount = "https://%s:%s/action/activate/%s"
)

// Password recover

func BuildPasswordRecoveryMail(host, port, username, token string) (string, string) {
	link := makePasswordRecoverLink(host, port, username, token)
	body := fmt.Sprintf(BodyRecoverPassword, username, link)
	return SubjectRecoverPassword, body
}

func makePasswordRecoverLink(host, port, user, token string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf(LinkRecoverPassword, host, port, encoded, token)
}

// Activate account

func BuildActivateAccountMail(host, port, username string) (string, string) {
	link := makeActivateAccountLink(host, port, username)
	body := fmt.Sprintf(BodyRecoverPassword, username, link)
	return SubjectRecoverPassword, body
}

func makeActivateAccountLink(host, port, user string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf(LinkActivateAccount, host, port, encoded)
}

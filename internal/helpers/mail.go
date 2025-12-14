package helpers

import (
	"encoding/base64"
	"fmt"
)

const (
	// Register
	SubjectRegister = `Welcome to Lenic, %s. Please verify the account` // username
	BodyRegister    = `
	<p>We're glad you could join us %s.</p>
	<p>Please click the link below to verify your account:</p>
	<p><a href="%s" target="_blank" rel="noopener">Activate your account</a></p> 
	` // link
	LinkActivateAccount = `http://%s:%s/action/activate/%s`

	// Password recover
	SubjectRecoverPassword = `Password Recovery for Lenic`
	BodyRecoverPassword    = `
	<p>Here's your password recovery link, %s.</p>
	<p>Please click the link below:</p>
	<p><a href="%s" target="_blank" rel="noopener">Recover your password</a></p>
	` // user - link
	LinkRecoverPassword = `http://%s:%s/recover-password/%s?token=%s`
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
	body := fmt.Sprintf(BodyRegister, username, link)
	subject := fmt.Sprintf(SubjectRegister, username)
	return subject, body
}

func makeActivateAccountLink(host, port, user string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf(LinkActivateAccount, host, port, encoded)
}

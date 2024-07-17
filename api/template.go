package api

type RegisterMail struct {
	User string
	Link string
}

var (
	registerSubject = `Welcome to the Blog {{.User}}. Please verify the account.`

	registerBody = `
	We're glad you could join us {{.User}}. Please click the link below to verify your account:
	{{.Link}}
	`
)

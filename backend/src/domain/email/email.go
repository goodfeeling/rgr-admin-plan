package email

type IEmailService interface {
	SendForgetPasswordEmail(userName string) error
}

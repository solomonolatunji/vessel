package services

type SESMailer struct{}

func NewSESMailer() *SESMailer {
	return &SESMailer{}
}

func (m *SESMailer) SendEmail(to, subject, body string) error {
	return nil
}

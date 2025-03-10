package mail

import (
	"fmt"
	"io"
	"net/smtp"
	"os"
)

var SMTP_SERVER = os.Getenv("SMTP_SERVER")
var SENDER_MAIL = os.Getenv("SENDER_MAIL")

// MailRequest는 POST /send_mail 요청에서 전달받을 JSON 데이터 구조입니다.
type MailRequest struct {
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Msg     []byte `json:"msg" binding:"required"`
}

type MailResponse struct {
	Result bool   `json:"result"`
	Msg    []byte `json:"msg"`
}

type Mail struct {
	Client *smtp.Client
	Header map[string]string
	Wc     io.WriteCloser
}

func NewMail() *Mail {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(SMTP_SERVER)
	if err != nil {
		fmt.Errorf("===ERROR===\n%+v", err)
		return nil
	}

	m := Mail{
		Client: c,
		Header: make(map[string]string),
	}

	return &m
}

func (m *Mail) InitMailSenderAndRecipient(to string) error {
	err := m.Client.Mail(SENDER_MAIL)
	if err != nil {
		return fmt.Errorf("===ERROR[%s]===\n%+v", "Sender information error", err)
	}

	err = m.Client.Rcpt(to)
	if err != nil {
		return fmt.Errorf("===ERROR[%s]===\n%+v", "To address information error", err)
	}

	return nil
}

func (m *Mail) SetMailHeaderBodyToWc(from, to, subject, context string) error {
	m.Header["From"] = from
	m.Header["To"] = to
	m.Header["Subject"] = subject

	var err error

	m.Wc, err = m.Client.Data()
	if err != nil {
		return fmt.Errorf("===ERROR[%s]===\n%+v", "Mail body initialization error", err)
	}

	for key, value := range m.Header {
		_, err := fmt.Fprintf(m.Wc, "%s: %s\r\n", key, value)
		if err != nil {
			return fmt.Errorf("===ERROR[%s]===\n%+v", "Mail body configuration error", err)
		}
	}

	_, err = fmt.Fprintf(m.Wc, context)
	if err != nil {
		return fmt.Errorf("===ERROR[%s]===\n%+v", "Mail body configuration error", err)
	}

	return nil
}

func (m *Mail) SendMail() error {
	err := m.Wc.Close()
	if err != nil {
		return fmt.Errorf("===ERROR[%s]===\n%+v", "Mail body connection close error", err)
	}

	err = m.Client.Quit()
	if err != nil {
		return fmt.Errorf("===ERROR[%s]===\n%+v", "\"Mail body quit mail error\"", err)
	}

	return nil
}

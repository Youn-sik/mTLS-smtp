package mail

import (
	"fmt"
	"io"
	"mailHttpToSmtp/utils"
	"net/smtp"
)

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

func VerifyEnv() {
	if utils.SMTP_SERVER == "" {
		panic(fmt.Errorf("===ERROR===\n%+s", "No Env Readed: SMTP_SERVER"))
	}
	if utils.SENDER_MAIL == "" {
		panic(fmt.Errorf("===ERROR===\n%+s", "No Env Readed: SENDER_MAIL"))
	}
}

func NewMail() *Mail {
	VerifyEnv()

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(utils.SMTP_SERVER)
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
	err := m.Client.Mail(utils.SENDER_MAIL)
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
	m.Header["Date"] = utils.GetDate()
	m.Header["Message-ID"] = utils.GetMessageId(utils.MAIL_DOMAIN)

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

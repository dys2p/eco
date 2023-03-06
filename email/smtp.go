package email

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type SMTP struct {
	From     string `json:"from"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

func (mailer SMTP) auth() sasl.Client {
	return sasl.NewPlainClient("", mailer.Username, mailer.Password)
}

func (mailer SMTP) hostAddr() string {
	if strings.Contains(mailer.Host, ":") {
		return mailer.Host
	}
	return mailer.Host + ":465"
}

func CreateConfig(jsonPath string) error {
	data, err := json.Marshal(&SMTP{})
	if err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, data, 0600); err != nil {
		return err
	}
	return fmt.Errorf("created empty config file: %s", jsonPath)
}

func LoadSMTP(jsonPath string) (*SMTP, error) {
	data, err := os.ReadFile(jsonPath)
	if os.IsNotExist(err) {
		return nil, CreateConfig(jsonPath)
	}
	if err != nil {
		return nil, err
	}

	var mailer = &SMTP{}
	if err := json.Unmarshal(data, mailer); err != nil {
		return nil, fmt.Errorf("unmarshaling json: %w", err)
	}

	// check authentication
	client, err := smtp.DialTLS(mailer.hostAddr(), nil)
	if err != nil {
		return nil, fmt.Errorf("dialing host: %w", err)
	}
	if err := client.Auth(mailer.auth()); err != nil {
		return nil, fmt.Errorf("authenticating: %w", err)
	}
	if err := client.Close(); err != nil {
		return nil, err
	}

	return mailer, nil
}

func (mailer SMTP) Send(to string, subject string, body []byte) error {
	mail, err := MakeEmail(mailer.From, to, subject, body)
	if err != nil {
		return err
	}

	return smtp.SendMailTLS(mailer.hostAddr(), mailer.auth(), mailer.From, []string{to}, mail)
}

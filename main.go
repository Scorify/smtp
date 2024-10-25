package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/scorify/schema"
)

type Schema struct {
	Server    string `key:"server"`
	Port      int    `key:"port" default:"25"`
	Username  string `key:"username"`
	Password  string `key:"password"`
	Sender    string `key:"sender"`
	Recipient string `key:"recipient"`
	Body      string `key:"body"`
	Secure    bool   `key:"secure"`
}

func Validate(config string) error {
	conf := Schema{}

	err := schema.Unmarshal([]byte(config), &conf)
	if err != nil {
		return err
	}

	return nil
}

// net/smtp does not allow you to send smtp.PlainAuth with not using TLS
// This is able to trick net/smtp to always thinking you are authenticating over TLS
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	server.TLS = true
	return a.Auth.Start(server)
}

func Run(ctx context.Context, config string) error {
	conf := Schema{}

	err := schema.Unmarshal([]byte(config), &conf)
	if err != nil {
		return err
	}

	var conn net.Conn
	connStr := fmt.Sprintf("%s:%d", conf.Server, conf.Port)

	if conf.Secure {
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("context deadline is not set")
		}

		conn, err = tls.DialWithDialer(
			&net.Dialer{Deadline: deadline},
			"tcp",
			connStr,
			&tls.Config{InsecureSkipVerify: true},
		)
	} else {
		dialer := &net.Dialer{}
		conn, err = dialer.DialContext(ctx, "tcp", connStr)
	}
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, conf.Server)
	if err != nil {
		return fmt.Errorf("failed creating new smtp client: %w", err)
	}
	defer client.Close()

	err = client.Auth(unencryptedAuth{smtp.PlainAuth("", conf.Username, conf.Password, conf.Server)})
	if err != nil {
		return fmt.Errorf("failed to authenticate to server: %w", err)
	}

	err = client.Mail(conf.Sender)
	if err != nil {
		return fmt.Errorf("failed to set sender %q: %w", conf.Sender, err)
	}

	err = client.Rcpt(conf.Recipient)
	if err != nil {
		return fmt.Errorf("failed to set recipient %q: %w", conf.Recipient, err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %w", err)
	}

	_, err = fmt.Fprint(wc, conf.Body)
	if err != nil {
		defer wc.Close()
		return fmt.Errorf("failed to write email body: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("failed to close body writer: %w", err)
	}

	return nil
}

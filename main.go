package smtp

import (
	"context"

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

	return nil
}

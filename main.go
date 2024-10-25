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

func Run(ctx context.Context, config string) error {
	conf := Schema{}

	err := schema.Unmarshal([]byte(config), &conf)
	if err != nil {
		return err
	}

	return nil
}

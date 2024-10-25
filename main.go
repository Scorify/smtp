package smtp

import (
	"context"

	"github.com/scorify/schema"
)

type Schema struct {
	Target string `key:"target"`
	Port   int    `key:"port"`
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

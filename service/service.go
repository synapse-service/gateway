package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	l "github.com/rs/zerolog/log"
)

func New(options ...option) (*Service, error) {
	s := Service{
		log: l.With().Str("component", "service").Logger(),
	}

	for _, option := range options {
		if err := option(&s); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	return &s, nil
}

type Service struct {
	log zerolog.Logger
}

func (*Service) Start(context.Context) error { return nil }
func (*Service) Stop(context.Context) error  { return nil }

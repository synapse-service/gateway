package grpc

import (
	"context"
	"net"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	l "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func New(options ...option) (*GRPC, error) {
	g := GRPC{
		log: l.With().Str("component", "grpc").Logger(),
	}

	for _, option := range options {
		if err := option(&g); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	if g.address == "" {
		return nil, errors.New("empty address")
	}
	if g.service == nil {
		return nil, errors.New("empty service")
	}
	if g.credentials == nil {
		return nil, errors.New("empty credentials")
	}

	g.Server = grpc.NewServer(
		grpc.Creds(g.credentials),
		grpc.ChainStreamInterceptor(
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		),
	)
	reflection.Register(&g)
	RegisterGatewayAPIServer(g.Server, &g)

	return &g, nil
}

type GRPC struct {
	UnimplementedGatewayAPIServer
	*grpc.Server
	log         zerolog.Logger
	address     string
	service     Service
	credentials credentials.TransportCredentials
}

func (g *GRPC) Start(context.Context) error {
	conn, err := net.Listen("tcp", g.address)
	if err != nil {
		return errors.Wrapf(err, "cannot listen %q", g.address)
	}

	go func() {
		g.log.Debug().Msgf("start listening %q", g.address)
		if err := g.Server.Serve(conn); err != nil {
			g.log.Fatal().Err(err).Msg("cannot serve connection")
		}
	}()

	return nil
}

func (g *GRPC) Stop(context.Context) error {
	g.Server.GracefulStop()
	return nil
}

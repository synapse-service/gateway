package client

import (
	"context"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/pkg/errors"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/synapse-service/gateway/transport/grpc"
)

func NewClient(options ...option) (*Client, error) {
	var c Client

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	if c.address == "" {
		return nil, errors.New("empty address")
	}
	if c.credentials == nil {
		return nil, errors.New("empty credentials")
	}

	return &c, nil
}

type Client struct {
	grpc.GatewayAPIClient
	address     string
	credentials credentials.TransportCredentials
	conn        *g.ClientConn
}

func (c *Client) Start(ctx context.Context) error {
	var err error
	c.conn, err = g.Dial(c.address,
		g.WithTransportCredentials(c.credentials),
		g.WithChainStreamInterceptor(
			grpc_opentracing.StreamClientInterceptor(),
		),
	)
	if err != nil {
		return err
	}

	c.GatewayAPIClient = grpc.NewGatewayAPIClient(c.conn)
	return nil
}

func (c *Client) Stop(ctx context.Context) error { return c.conn.Close() }

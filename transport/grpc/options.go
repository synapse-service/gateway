package grpc

import "google.golang.org/grpc/credentials"

type option = func(*GRPC) error

func WithAddress(address string) option {
	return func(g *GRPC) error {
		g.address = address
		return nil
	}
}

func WithService(service Service) option {
	return func(g *GRPC) error {
		g.service = service
		return nil
	}
}

func WithCredentials(certificate, key string) option {
	return func(g *GRPC) error {
		credentials, err := credentials.NewServerTLSFromFile(certificate, key)
		if err != nil {
			return err
		}
		g.credentials = credentials
		return nil
	}
}

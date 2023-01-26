package client

import "google.golang.org/grpc/credentials"

type option = func(*Client) error

func WithClientAddress(address string) option {
	return func(c *Client) error {
		c.address = address
		return nil
	}
}

func WithTransportCredentials(certificate string) option {
	return func(c *Client) error {
		credentials, err := credentials.NewClientTLSFromFile(certificate, "")
		if err != nil {
			return err
		}
		c.credentials = credentials
		return nil
	}
}

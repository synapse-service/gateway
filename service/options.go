package service

type option = func(*Service) error

type Config struct{}

func WithConfig(cfg Config) option {
	return func(s *Service) error {
		return nil
	}
}

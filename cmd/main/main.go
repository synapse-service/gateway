package main

import (
	"flag"
	l "log"
	"os"

	"github.com/242617/core/application"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/synapse-service/gateway/service"
	"github.com/synapse-service/gateway/transport/grpc"
)

type Config struct {
	Service   service.Config `yaml:"config"`
	Transport struct {
		GRPC struct {
			Address     string `yaml:"address"`
			Certificate string `yaml:"certificate"`
			Key         string `yaml:"key"`
		} `yaml:"grpc"`
	} `yaml:"transport"`
}

func init() { l.SetFlags(l.Lshortfile) }
func main() {
	configPath := flag.String("config", "/etc/service/config.yaml", "Config file path")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.000"})

	var cfg Config
	if b, err := os.ReadFile(*configPath); err != nil {
		l.Fatal(errors.Wrap(err, "load config"))
	} else if err := yaml.Unmarshal(b, &cfg); err != nil {
		l.Fatal(errors.Wrap(err, "unmarshal config"))
	}

	service, err := service.New(
		service.WithConfig(cfg.Service),
	)
	if err != nil {
		l.Fatal(errors.Wrap(err, "create service"))
	}

	grpc, err := grpc.New(
		grpc.WithAddress(cfg.Transport.GRPC.Address),
		grpc.WithService(service),
		grpc.WithCredentials(cfg.Transport.GRPC.Certificate, cfg.Transport.GRPC.Key),
	)
	if err != nil {
		l.Fatal(errors.Wrap(err, "create grpc"))
	}

	app, err := application.New(
		application.WithComponents(
			application.NewLifecycleComponent("grpc", grpc),
			application.NewLifecycleComponent("service", service),
		),
	)
	if err != nil {
		l.Fatal(errors.Wrap(err, "create application"))
	}
	if err := app.Run(); err != nil {
		l.Fatal(errors.Wrap(err, "run application"))
	}
}

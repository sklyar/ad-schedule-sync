package config

import (
	"errors"
	"fmt"
)

type ServerHTTP struct {
	Addr      string `yaml:"addr"`
	SecretKey string `yaml:"secret_key"`
}

func (s ServerHTTP) Validate() error {
	var errs []error

	if s.Addr == "" {
		errs = append(errs, fmt.Errorf("addr is required"))
	}
	if s.SecretKey == "" {
		errs = append(errs, fmt.Errorf("secret_key is required"))
	}

	return errors.Join(errs...)
}

type Server struct {
	HTTP ServerHTTP `yaml:"http"`
}

func (s Server) Validate() error {
	var errs []error

	if err := s.HTTP.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("http: %w", err))
	}

	return errors.Join(errs...)
}

type Logging struct {
	Level string `yaml:"level"`
}

type Config struct {
	Server  Server  `yaml:"server"`
	Logging Logging `yaml:"logging"`
}

func (c Config) Validate() error {
	var errs []error

	if err := c.Server.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("server: %w", err))
	}

	return errors.Join(errs...)
}

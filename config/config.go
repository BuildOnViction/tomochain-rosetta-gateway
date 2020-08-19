// Copyright (c) 2020 TomoChain

package config

import (
	"github.com/pkg/errors"
	uconfig "go.uber.org/config"
)

type (
	NetworkIdentifier struct {
		Blockchain string `yaml:"blockchain"`
		Network    string `yaml:"network"`
	}
	Currency struct {
		Symbol   string `yaml:"symbol"`
		Decimals int32  `yaml:"decimals"`
	}
	Server struct {
		Port           string `yaml:"port"`
		Endpoint       string `yaml:"endpoint"`
		SecureEndpoint bool   `yaml:"secureEndpoint"`
		RosettaVersion string `yaml:"rosettaVersion"`
	}
	Config struct {
		NetworkIdentifier NetworkIdentifier `yaml:"network_identifier"`
		Currency          Currency          `yaml:"currency"`
		Server            Server            `yaml:"server"`
	}
)

func New(path string) (cfg *Config, err error) {
	opts := []uconfig.YAMLOption{uconfig.File(path)}
	yaml, err := uconfig.NewYAML(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}
	cfg = &Config{}
	if err := yaml.Get(uconfig.Root).Populate(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML config to struct")
	}
	return
}
package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens"`
}

func DefaultConfig() *Config {
	return &Config{
		Host:      "localhost",
		Port:      11434,
		Model:     "llama3",
		MaxTokens: 1024,
	}
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parsing config file %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	// Environment variable overrides
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		cfg.Host = host
	}
	if portStr := os.Getenv("OLLAMA_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = p
		}
	}
	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		cfg.Model = model
	}
	if maxTokensStr := os.Getenv("OLLAMA_MAX_TOKENS"); maxTokensStr != "" {
		if mt, err := strconv.Atoi(maxTokensStr); err == nil {
			cfg.MaxTokens = mt
		}
	}

	return cfg, nil
}

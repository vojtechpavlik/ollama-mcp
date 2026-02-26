package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.Host)
	}
	if cfg.Port != 11434 {
		t.Errorf("expected port 11434, got %d", cfg.Port)
	}
	if cfg.Model != "llama3" {
		t.Errorf("expected model llama3, got %s", cfg.Model)
	}
	if cfg.MaxTokens != 1024 {
		t.Errorf("expected max_tokens 1024, got %d", cfg.MaxTokens)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	cfg, err := LoadConfig("non-existent.yaml")
	if err != nil {
		t.Fatalf("LoadConfig should not return error for missing file: %v", err)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.Host)
	}
}

func TestLoadConfig_WithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
host: customhost
port: 1234
model: custommodel
max_tokens: 2048
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Host != "customhost" {
		t.Errorf("expected customhost, got %s", cfg.Host)
	}
	if cfg.Port != 1234 {
		t.Errorf("expected 1234, got %d", cfg.Port)
	}
	if cfg.Model != "custommodel" {
		t.Errorf("expected custommodel, got %s", cfg.Model)
	}
	if cfg.MaxTokens != 2048 {
		t.Errorf("expected 2048, got %d", cfg.MaxTokens)
	}
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	os.Setenv("OLLAMA_HOST", "envhost")
	os.Setenv("OLLAMA_PORT", "9999")
	os.Setenv("OLLAMA_MODEL", "envmodel")
	os.Setenv("OLLAMA_MAX_TOKENS", "4096")
	defer func() {
		os.Unsetenv("OLLAMA_HOST")
		os.Unsetenv("OLLAMA_PORT")
		os.Unsetenv("OLLAMA_MODEL")
		os.Unsetenv("OLLAMA_MAX_TOKENS")
	}()

	cfg, err := LoadConfig("non-existent.yaml")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Host != "envhost" {
		t.Errorf("expected envhost, got %s", cfg.Host)
	}
	if cfg.Port != 9999 {
		t.Errorf("expected 9999, got %d", cfg.Port)
	}
	if cfg.Model != "envmodel" {
		t.Errorf("expected envmodel, got %s", cfg.Model)
	}
	if cfg.MaxTokens != 4096 {
		t.Errorf("expected 4096, got %d", cfg.MaxTokens)
	}
}

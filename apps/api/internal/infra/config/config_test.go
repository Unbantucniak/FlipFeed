package infraconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSigningKeyAndInternalToken(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("port: 8080\n"), 0600); err != nil {
		t.Fatal(err)
	}
	load := func() *Config {
		t.Helper()
		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatal(err)
		}
		return cfg
	}
	t.Setenv("GCFEED_JWT_SECRET", "")
	t.Setenv("GCFEED_INTERNAL_TOKEN", "")
	first, second := load(), load()
	if first.JWT.Secret == "" || first.JWT.Secret == second.JWT.Secret || first.Internal.Token != "" {
		t.Fatal("expected distinct generated signing keys and an unconfigured internal token")
	}
	if err := os.WriteFile(path, []byte("jwt:\n  secret: from-file\ninternal:\n  token: file-token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if cfg := load(); cfg.JWT.Secret != "from-file" || cfg.Internal.Token != "file-token" {
		t.Fatal("file configuration not preserved")
	}
	t.Setenv("GCFEED_JWT_SECRET", "from-env")
	t.Setenv("GCFEED_INTERNAL_TOKEN", "env-token")
	if cfg := load(); cfg.JWT.Secret != "from-env" || cfg.Internal.Token != "env-token" {
		t.Fatal("environment configuration not applied")
	}
}

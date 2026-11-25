package infraconfig

import (
	"crypto/rand"
	"errors"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

var ErrEmptyConfigPath = errors.New("config file path is empty")
var ErrReadConfigFailed = errors.New("read config file failed")
var ErrUnmarshalConfigFailed = errors.New("unmarshal config failed")

// LoadConfig 读取 YAML 配置文件，并反序列化为应用启动配置。
func LoadConfig(path string) (*Config, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, ErrEmptyConfigPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrReadConfigFailed
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, ErrUnmarshalConfigFailed
	}

	if secret := os.Getenv("FLIPFEED_JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = rand.Text()
	}
	if token := os.Getenv("FLIPFEED_INTERNAL_TOKEN"); token != "" {
		cfg.Internal.Token = token
	}
	return cfg, nil
}

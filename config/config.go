package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/starudream/go-lib/core/v2/codec/yaml"
	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

type Config struct {
	Events map[string]map[string]EventAction `json:"events" yaml:"events"`
}

type EventAction struct {
	Disabled bool   `json:"disabled" yaml:"disabled"`
	Template string `json:"template" yaml:"template"`
}

var (
	_c   = Config{}
	_cMu = sync.Mutex{}
)

func init() {
	_ = config.Unmarshal("", &_c)
	config.LoadStruct(_c)
}

func C() Config {
	_cMu.Lock()
	defer _cMu.Unlock()
	return _c
}

func Save() error {
	config.LoadStruct(_c)

	bs, err := yaml.Marshal(config.Raw())
	if err != nil {
		return fmt.Errorf("marshal config error: %w", err)
	}

	filename := config.LoadedFile()
	if filename == "" {
		filename = filepath.Join(osutil.ExeDir(), osutil.ExeName()+".yaml")
		slog.Info("config file not found, save to default file", slog.String("file", filename))
	}

	err = os.WriteFile(config.LoadedFile(), bs, 0644)
	if err != nil {
		return fmt.Errorf("write config file error: %w", err)
	}

	slog.Info("save config success", slog.String("file", filename))

	return nil
}

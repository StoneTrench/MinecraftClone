package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
)

type ModuleConfigId uint32

type ModuleConfig struct {
	FilePath      string             `json:"-"`
	StringEntries map[string]string  `json:"string_entries"`
	IntEntries    map[string]int64   `json:"int_entries"`
	FloatEntries  map[string]float64 `json:"float_entries"`
}

var already_open map[string]ModuleConfigId = make(map[string]ModuleConfigId)
var open_configs []ModuleConfig

const CONFIG_FOLDER = "./config/"
const ENGINE_CONFIG = "engine.json"

var name_replacer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func LoadConfig(name string) (ModuleConfigId, error) {
	config_path := CONFIG_FOLDER
	if len(name) == 0 {
		config_path = path.Join(config_path, ENGINE_CONFIG)
	} else {
		config_path = path.Join(config_path, fmt.Sprintf("%s_config.json", name_replacer.ReplaceAllString(name, "_")))
	}
	if id, exists := already_open[config_path]; exists {
		cfg := &open_configs[id]
		data, err := os.ReadFile(config_path)
		if err != nil {
			return 0, err
		}
		
		if len(data) == 0 {
			cfg.StringEntries = make(map[string]string)
			cfg.IntEntries = make(map[string]int64)
			cfg.FloatEntries = make(map[string]float64)
			return id, nil
		}

		var tmp ModuleConfig
		if err := json.Unmarshal(data, &tmp); err != nil {
			return 0, err
		}

		cfg.StringEntries = tmp.StringEntries
		cfg.IntEntries = tmp.IntEntries
		cfg.FloatEntries = tmp.FloatEntries

		return id, nil
	}

	os.MkdirAll(CONFIG_FOLDER, 0755)

	var cfg ModuleConfig
	cfg.FilePath = config_path
	cfg.StringEntries = make(map[string]string)
	cfg.IntEntries = make(map[string]int64)
	cfg.FloatEntries = make(map[string]float64)

	data, err := os.ReadFile(config_path)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return 0, err
		}
	}

	id := ModuleConfigId(len(open_configs))
	open_configs = append(open_configs, cfg)
	already_open[config_path] = id

	return id, nil
}
func FlushAll() (err error) {
	for i := range open_configs {
		if e := FlushConfig(ModuleConfigId(i)); e != nil {
			err = errors.Join(e)
		}
	}
	return err
}
func FlushConfig(id ModuleConfigId) error {
	if int(id) >= len(open_configs) {
		return fmt.Errorf("invalid config id")
	}
	cfg := open_configs[id]
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.FilePath, data, 0644)
}
func ClearConfigs() {
	open_configs = nil
	already_open = make(map[string]ModuleConfigId)
}

func GetString(id ModuleConfigId, key string) (string, bool) {
	if int(id) >= len(open_configs) {
		return "", false
	}
	val, ok := open_configs[id].StringEntries[key]
	return val, ok
}
func GetInt(id ModuleConfigId, key string) (int64, bool) {
	if int(id) >= len(open_configs) {
		return 0, false
	}
	val, ok := open_configs[id].IntEntries[key]
	return val, ok
}
func GetFloat(id ModuleConfigId, key string) (float64, bool) {
	if int(id) >= len(open_configs) {
		return 0, false
	}
	val, ok := open_configs[id].FloatEntries[key]
	return val, ok
}

func SetString(id ModuleConfigId, key string, value string) {
	if int(id) < len(open_configs) {
		open_configs[id].StringEntries[key] = value
	}
}
func SetInt(id ModuleConfigId, key string, value int64) {
	if int(id) < len(open_configs) {
		open_configs[id].IntEntries[key] = value
	}
}
func SetFloat(id ModuleConfigId, key string, value float64) {
	if int(id) < len(open_configs) {
		open_configs[id].FloatEntries[key] = value
	}
}

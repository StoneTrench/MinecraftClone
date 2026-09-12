package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	FilePath      string             `json:"-"`
	StringEntries map[string]string  `json:"string_entries"`
	IntEntries    map[string]int64   `json:"int_entries"`
	FloatEntries  map[string]float64 `json:"float_entries"`
}

func (c *Config) Init(path string) {
	c.FilePath = path
	c.StringEntries = make(map[string]string)
	c.IntEntries = make(map[string]int64)
	c.FloatEntries = make(map[string]float64)
}

func (c *Config) OpenOrCreate() error {
	file, err := os.Open(c.FilePath)
	if errors.Is(err, os.ErrNotExist) {
		if file != nil {
			err = file.Close()
			if err != nil {
				return err
			}
		}
		file, err = os.Create(c.FilePath)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return err
	}
	if buf.Len() == 0 {
		return nil
	}

	var tmp Config
	err = json.Unmarshal(buf.Bytes(), &tmp)
	if err != nil {
		return err
	}

	c.StringEntries = tmp.StringEntries
	c.IntEntries = tmp.IntEntries
	c.FloatEntries = tmp.FloatEntries

	return nil
}
func (c *Config) FlushConfig() error {
	data, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.FilePath, data, 0644)
}

func (c *Config) GetString(key string) (string, bool) {
	val, ok := c.StringEntries[key]
	return val, ok
}
func (c *Config) GetInt(key string) (int64, bool) {
	val, ok := c.IntEntries[key]
	return val, ok
}
func (c *Config) GetFloat(key string) (float64, bool) {
	val, ok := c.FloatEntries[key]
	return val, ok
}

func (c *Config) SetString(key string, value string) {
	c.StringEntries[key] = value
}
func (c *Config) SetInt(key string, value int64) {
	c.IntEntries[key] = value
}
func (c *Config) SetFloat(key string, value float64) {
	c.FloatEntries[key] = value
}

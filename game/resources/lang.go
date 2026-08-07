package resources

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/StoneTrench/go-mc-clone/game/log"
)

var language_map map[string]string = make(map[string]string)
var fallback_language_map map[string]string = make(map[string]string)

func LangTranslate(key string, args ...any) string {
	format, exists := language_map[key]
	if !exists {
		log.Warnf("no translation found for key (%s)", key)
		format, exists = fallback_language_map[key]
	}
	if !exists {
		log.Warnf("no fallback translation found for key (%s)", key)
		return key
	}
	return fmt.Sprintf(format, args...)
}
func LangClear() {
	for k := range language_map {
		delete(language_map, k)
	}
}
func append_lang(path string, m map[string]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open json language file, %w", err)
	}

	lang := make(map[string]string)
	err = json.Unmarshal(data, &lang)
	if err != nil {
		return fmt.Errorf("failed to parse json language file, %w", err)
	}
	for k, v := range lang {
		if other_v, exists := m[k]; exists {
			log.Warnf("overlapping language keys (%s) \"%s\" and \"%s\"", k, v, other_v)
		}
		m[k] = v
	}

	return nil
}
func LangAppend(path string) error {
	return append_lang(path, language_map)
}
func LangAppendFallback(path string) error {
	return append_lang(path, fallback_language_map)
}

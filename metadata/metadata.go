package metadata

import (
	"fmt"
)

type ENGINE_MODE uint8

const (
	MODE_RELEASE ENGINE_MODE = iota
	MODE_DEBUG
)

var (
	BUILD_NAME   string = "" // Do not read from these!
	BUILD_TAG    string = "" // Do not read from these!
	BUILD_BUILT  string = "" // Do not read from these!
	BUILD_COMMIT string = "" // Do not read from these!
	BUILD_MODE   string = "" // Do not read from these!
)

var (
	internal_is_initialized bool = false
	internal_name           string
	internal_tag            string
	internal_built          string
	internal_commit         string
	internal_mode           ENGINE_MODE
)

func Init() error {
	if internal_is_initialized {
		return fmt.Errorf("metadata already initialized")
	}
	if len(BUILD_NAME) == 0 {
		return fmt.Errorf("name was nil")
	}
	if len(BUILD_TAG) == 0 {
		return fmt.Errorf("git tag was nil")
	}
	if len(BUILD_COMMIT) == 0 {
		return fmt.Errorf("git commit was nil")
	}
	if len(BUILD_BUILT) == 0 {
		return fmt.Errorf("built timestamp was nil")
	}
	if len(BUILD_MODE) == 0 {
		return fmt.Errorf("engine mode was nil")
	}

	internal_name = BUILD_NAME
	internal_tag = BUILD_TAG
	internal_commit = BUILD_COMMIT
	internal_built = BUILD_BUILT
	switch BUILD_MODE {
	case "RELEASE":
		internal_mode = MODE_RELEASE
	case "DEBUG":
		internal_mode = MODE_DEBUG
	default:
		return fmt.Errorf("invalid engine mode: %s", BUILD_MODE)
	}

	BUILD_NAME = ""
	BUILD_TAG = ""
	BUILD_COMMIT = ""
	BUILD_BUILT = ""
	BUILD_MODE = ""

	internal_is_initialized = true
	return nil
}

func GetFormattedApplicationLabel() string {
	if !internal_is_initialized {
		panic("metadata not initialized")
	}

	var suffix = ""
	if IsInDebugMode() {
		suffix = " [DEBUG]"
	}
	return fmt.Sprintf("%s  —  (%s) %s %s%s", internal_name, internal_tag, internal_built, internal_commit, suffix)
}

func IsInDebugMode() bool {
	if !internal_is_initialized {
		panic("metadata not initialized")
	}
	return internal_mode == MODE_DEBUG
}

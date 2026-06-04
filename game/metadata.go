package game

import (
	"fmt"
)

var (
	PROJECT_NAME   = "NIL"
	PROJECT_TAG    = "NIL"
	PROJECT_COMMIT = "NIL"
	PROJECT_BUILT  = "NIL"
)

func GetFormattedApplicationLabel() string {
	return fmt.Sprintf("%s  —  (%s) %s %s", PROJECT_NAME, PROJECT_TAG, PROJECT_BUILT, PROJECT_COMMIT)
}

package tarmo

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var rawVersion string

func Version() string {
	return strings.TrimSpace(rawVersion)
}

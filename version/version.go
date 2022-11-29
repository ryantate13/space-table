package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var version string

func Get() string {
	return strings.TrimSpace(version)
}

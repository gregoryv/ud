package ud

import (
	_ "embed"
	"runtime/debug"
	"strings"
)

func Version() string {
	prefix := "## ["
	from := strings.Index(changelog, prefix) + len(prefix)
	to := from + strings.Index(changelog[from:], "]")
	return changelog[from:to]
}

//go:embed changelog.md
var changelog string

func Revision() string {
	// this is a module so it will always work
	info, _ := debug.ReadBuildInfo()
	return findBuildValue(info.Settings, "vcs.revision")[:6]
}

func findBuildValue(settings []debug.BuildSetting, key string) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return ""
}

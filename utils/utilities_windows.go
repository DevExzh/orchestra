//go:build windows

package utils

import "strings"

func ProcessSQLDataSource(path string) string {
	var p string
	p = strings.ReplaceAll(path, "\\", "/")
	p = strings.ReplaceAll(p, " ", "%20")
	p = "file:///" + p + "?mode=rwc"
	return p
}

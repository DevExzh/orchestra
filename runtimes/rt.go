package runtimes

import "strings"

const (
	OsDarwin OsType = iota
	OsWindows
	OsLinux
	OsUnknown
)

const (
	TypeSourceCode ReleaseFileType = iota
	TypeInstaller
	TypeCompressed
)

type (
	ReleaseVersion struct {
		Version         string
		Date            string
		Files           []ReleaseFileVersion
		FilesLink       string
		HasChecksum     bool
		DigestAlgorithm string
	}

	OsType          int
	ReleaseFileType int

	ReleaseFileVersion struct {
		Name        string
		Link        string
		FileType    ReleaseFileType
		SupportedOs OsType
		Digest      string
	}
)

func StringToOsType(str string) OsType {
	switch strings.ToLower(str) {
	case "macos":
		return OsDarwin
	case "windows":
		return OsWindows
	case "linux":
		return OsLinux
	default:
		return OsUnknown
	}
}

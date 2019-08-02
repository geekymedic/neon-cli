package types

import "runtime"

//go:generate stringer -type=OSType
type OSType int

const (
	OsUnkownType OSType = iota
	MacOs
	LinuxOs
	WindowsOs
)

func OsType() OSType {
	switch runtime.GOOS {
	case "darwin":
		return MacOs
	case "linux":
		return LinuxOs
	case "windows":
		return WindowsOs
	default:
		return OsUnkownType
	}
}

func OsSep() string {
	switch runtime.GOOS {
	case "darwin", "linux":
		return ":"
	case "windows":
		return ";"
	default:
		return "OsUnkownType"
	}
}

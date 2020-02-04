package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geekymedic/neon-cli/types"
)

func AbsDir(dir string) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("not exist %s directory", dir)
	}
	if len(dir) == 1 && dir[0] == '~' {
		if types.OsType() == types.MacOs || types.OsType() == types.LinuxOs {
			dir = fmt.Sprintf("%s/%s", os.Getenv("HOME"), dir[1:])
		}
	} else if len(dir) == 1 && dir[0] == '.' {
		if types.OsType() == types.WindowsOs {
			dir, _ = os.Getwd()
		}
	} else if len(dir) >= 2 && dir[0:2] == "~/" {
		if types.OsType() == types.MacOs || types.OsType() == types.LinuxOs {
			dir = fmt.Sprintf("%s/%s", os.Getenv("HOME"), dir[2:])
		}
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return filepath.Clean(dir), err
}

func ConvertBreakLinePath(s string) string {
	switch types.OsType() {
	case types.MacOs, types.LinuxOs:
		return s
	case types.WindowsOs:
		return strings.Replace(s, "/", "\\", -1)
	}
	types.PanicSanity("Unimplement")
	return ""
}
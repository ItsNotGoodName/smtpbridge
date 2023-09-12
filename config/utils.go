package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func flagUsageString(s string, usage ...string) string {
	u := ""
	if usage != nil {
		if s == "" {
			return usage[0]
		}
		u += usage[0] + " "
	} else if s == "" {
		return ""
	}
	return fmt.Sprintf("%s(default \"%s\")", u, s)
}

func flagUsageInt(i int, usage ...string) string {
	u := ""
	if usage != nil {
		if i == 0 {
			return usage[0]
		}
		u += usage[0] + " "
	} else if i == 0 {
		return ""
	}
	return fmt.Sprintf("%s(default %d)", u, i)
}

func flagUsageBool(b bool, usage ...string) string {
	u := ""
	if usage != nil {
		if b == false {
			return usage[0]
		}
		u += usage[0] + " "
	} else if b == false {
		return ""
	}
	return fmt.Sprintf("%s(default true)", u)
}

func resolve(paths []string) (string, error) {
	for _, path := range paths {
		f, err := os.Open(expandPath(path))
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				continue
			}

			return "", err
		}
		f.Close()

		return path, nil
	}

	return "", nil
}

func expandPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if strings.HasPrefix(path, "~/") {
		user, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(user.HomeDir, path[2:])
	}
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abspath
}

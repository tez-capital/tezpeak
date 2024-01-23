package util

import (
	"os"
	"time"
)

func CheckFileChangedTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}

	return info.ModTime()
}

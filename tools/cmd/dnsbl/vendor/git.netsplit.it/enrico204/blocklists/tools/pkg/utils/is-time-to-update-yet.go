package utils

import (
	"errors"
	"os"
	"time"
)

func IsTimeToUpdateYet(canaryFile string, updateInterval time.Duration) (bool, time.Duration, error) {
	stat, err := os.Stat(canaryFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, 0, err
	} else if err == nil {
		// File exists, retrieve the last modified
		lastModified := stat.ModTime()
		return time.Since(lastModified) >= updateInterval, updateInterval - time.Since(lastModified), nil
	}
	return true, 0, nil
}

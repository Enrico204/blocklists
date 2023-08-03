package utils

import (
	"os"
)

// MoveFileKeepingLastModifiedTime moves the file from the source path to the destination path. It keeps the last
// modified time intact.
func MoveFileKeepingLastModifiedTime(src, dest string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	lastModified := stat.ModTime()

	err = os.Rename(src, dest)
	if err != nil {
		return err
	}

	return os.Chtimes(dest, lastModified, lastModified)
}

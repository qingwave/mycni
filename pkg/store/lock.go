package store

import (
	"os"
	"path"

	"github.com/alexflint/go-filemutex"
)

func newFileLock(lockPath string) (*filemutex.FileMutex, error) {
	fi, err := os.Stat(lockPath)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		lockPath = path.Join(lockPath, "lock")
	}

	f, err := filemutex.New(lockPath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

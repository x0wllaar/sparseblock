package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func sparseReadFileAt(fileName string, p []byte, off int64) (int, error) {
	f, err := os.Open(fileName)
	if errors.Is(err, fs.ErrNotExist) {
		// fill with 0xFF
		for i := 0; i < len(p); i += 1 {
			p[i] = 0xFF
		}
		return len(p), nil
	}
	if err != nil {
		return 0, fmt.Errorf("error sparse reading from file %v: %w", fileName, err)
	}
	defer f.Close()

	return f.ReadAt(p, off)
}

func sparseWriteFileAt(fileName string, fileSize int, p []byte, off int64) (int, error) {
	_, err := os.Stat(fileName)
	if errors.Is(err, fs.ErrNotExist) {
		fileBuf := make([]byte, fileSize)
		for i := 0; i < fileSize; i += 1 {
			fileBuf[i] = 0xFF
		}
		dir := filepath.Dir(fileName)
		err := os.MkdirAll(dir, 0o600)
		if err != nil {
			return 0, fmt.Errorf("error creating directory for file %v: %w", fileName, err)
		}
		err = os.WriteFile(fileName, fileBuf, 0o600)
		if err != nil {
			return 0, fmt.Errorf("error creating file %v: %w", fileName, err)
		}
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0o000)
	if err != nil {
		return 0, fmt.Errorf("error sparse writing to file %v: %w", fileName, err)
	}
	defer f.Close()

	return f.WriteAt(p, off)
}

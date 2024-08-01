package main

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type fileChunkBackend struct {
	baseDir    string
	fileSuffix string
	treespec   []int64
	fileSize   int64
	lock       sync.RWMutex
}

func (b *fileChunkBackend) Sync() error {
	return nil
}

func (b *fileChunkBackend) Size() (int64, error) {
	size := int64(1)
	for _, v := range b.treespec {
		size *= v
	}
	size *= b.fileSize
	return size, nil
}

func (b *fileChunkBackend) ReadAt(p []byte, off int64) (int, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	globalLogger.Debug("processing read request", slog.Int64("offset", off), slog.Int("length", len(p)))

	n := 0
	err := processChunkedOperation(off, int64(len(p)), b.fileSize, func(r chunkPositionRequest) error {
		globalLogger.Debug("processing read chunk", slog.Int64("originalOffset", r.originalOffset), slog.Int64("chunckNum", r.chunkNum), slog.Int64("chunkOffset", r.chunkOffset), slog.Int64("chunkLength", r.chunkLength))

		treePath, err := numToTreePath(r.chunkNum, b.treespec)
		if err != nil {
			return fmt.Errorf("error converting chunk num to path: %w", err)
		}
		fullTreePath := fmt.Sprintf("%s/%s%s", b.baseDir, treePath, b.fileSuffix)
		cn, err := sparseReadFileAt(fullTreePath, p[r.originalOffset:r.originalOffset+r.chunkLength], r.chunkOffset)
		if err != nil {
			if err == io.EOF {
				if cn != int(r.chunkLength) {
					return fmt.Errorf("read less than expected: EOF: %w", err)
				}
			} else {
				return fmt.Errorf("error reading chunk file: %w", err)
			}
		}
		n += cn
		return nil
	})
	if err != nil {
		globalLogger.Error("error reading from chunked files", "error", err)
		return 0, fmt.Errorf("error reading from chunked files: %w", err)
	}
	return n, nil
}

func (b *fileChunkBackend) WriteAt(p []byte, off int64) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	globalLogger.Debug("processing write request", slog.Int64("offset", off), slog.Int("length", len(p)))

	n := 0
	err := processChunkedOperation(off, int64(len(p)), b.fileSize, func(r chunkPositionRequest) error {
		globalLogger.Debug("processing write chunk", slog.Int64("originalOffset", r.originalOffset), slog.Int64("chunckNum", r.chunkNum), slog.Int64("chunkOffset", r.chunkOffset), slog.Int64("chunkLength", r.chunkLength))

		treePath, err := numToTreePath(r.chunkNum, b.treespec)
		if err != nil {
			return fmt.Errorf("error converting chunk num to path: %w", err)
		}
		fullTreePath := fmt.Sprintf("%s/%s%s", b.baseDir, treePath, b.fileSuffix)
		cn, err := sparseWriteFileAt(fullTreePath, int(b.fileSize), p[r.originalOffset:r.originalOffset+r.chunkLength], r.chunkOffset)
		if err != nil {
			return fmt.Errorf("error writing chunk file: %w", err)
		}
		n += cn
		return nil
	})
	if err != nil {
		globalLogger.Error("error writing to chunked files", "error", err)
		return 0, fmt.Errorf("error writing to chunked files: %w", err)
	}
	return n, nil
}

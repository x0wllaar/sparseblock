package main

import "fmt"

type chunkPositionRequest struct {
	originalOffset int64
	chunkNum       int64
	chunkOffset    int64
	chunkLength    int64
}
type chunkRequestProcessor func(req chunkPositionRequest) error

func processChunkedOperation(offset int64, length int64, chunkSize int64, proc chunkRequestProcessor) error {
	firstChunkN := offset / chunkSize
	firstChunkOffset := offset % chunkSize

	firstChunkRemaining := chunkSize - firstChunkOffset
	if length <= firstChunkRemaining {
		err := proc(chunkPositionRequest{0, firstChunkN, firstChunkOffset, length})
		if err != nil {
			return fmt.Errorf("error processing chunked request: %w", err)
		}
		return nil
	}
	err := proc(chunkPositionRequest{0, firstChunkN, firstChunkOffset, firstChunkRemaining})
	if err != nil {
		return fmt.Errorf("error processing chunked request: %w", err)
	}
	length = length - firstChunkRemaining
	originalOffset := firstChunkRemaining

	nWholeChunks := length / chunkSize
	inLastChunk := length % chunkSize
	for c := int64(0); c < nWholeChunks; c += 1 {
		err = proc(chunkPositionRequest{originalOffset, firstChunkN + 1 + c, 0, chunkSize})
		if err != nil {
			return fmt.Errorf("error processing chunked request: %w", err)
		}
		originalOffset = originalOffset + chunkSize
	}
	if inLastChunk != 0 {
		err = proc(chunkPositionRequest{originalOffset, firstChunkN + 1 + nWholeChunks, 0, inLastChunk})
		if err != nil {
			return fmt.Errorf("error processing chunked request: %w", err)
		}
	}

	return nil
}

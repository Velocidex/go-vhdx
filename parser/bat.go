package parser

import (
	"io"
	"sync"
)

type BatRange struct {
	FileOffset uint64
}

type BatReader struct {
	mu              sync.Mutex
	BlockSize       uint64
	EntriesPerChunk uint64
	Size            uint64
	Reader          io.ReaderAt

	bat map[int]*BatRange
}

func (self *BatReader) ReadAt(buff []byte, off int64) (int, error) {
	for buff_offset := 0; buff_offset < len(buff); {
		n, err := self.readPartial(buff[buff_offset:], off+int64(buff_offset))
		if err != nil {
			return buff_offset, err
		}

		buff_offset += n
	}
	return len(buff), nil
}

// Read as much as possible and return a short read if we exceed the
// block boundary.
func (self *BatReader) readPartial(buff []byte, off int64) (int, error) {
	self.mu.Lock()
	block := int(uint64(off) / self.BlockSize)

	// Calculate the element index.
	//
	// The BAT array is divided into a list of chunks. Each chunk
	// contains some payload blocks and one single sector block: See
	// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-vhdx/af7334e6-ad2c-4378-9b81-afc1334a6ee7
	//
	// This means that each chunk contain 1 less than the full number
	// of blocks that can fit in it (to allow for one sector block).
	// So we calculate the element index by adding one extra index per
	// chunk.
	element_index := uint64(block) + uint64(block)/self.EntriesPerChunk

	block_offset := uint64(off) % self.BlockSize
	to_read := int(self.BlockSize - block_offset)
	if to_read > len(buff) {
		to_read = len(buff)
	}

	bat_range, pres := self.bat[int(element_index)]
	self.mu.Unlock()

	if !pres {
		// Just null terminate it.
		for i := 0; i < to_read; i++ {
			buff[i] = 0
		}
		return to_read, nil
	}

	// Get the reader to read the correct offset.
	return self.Reader.ReadAt(
		buff[:to_read], int64(bat_range.FileOffset+block_offset))
}

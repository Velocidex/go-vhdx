package parser

import (
	"io"
	"sync"
)

type Bat struct {
	StartOffset, EntrySize uint64
	TotalNumberOfEntries   uint64

	EntriesPerChunk uint64

	profile *VHDXProfile
	reader  io.ReaderAt
}

func (self *Bat) GetFileOffset(block int) uint64 {
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

	if element_index > self.TotalNumberOfEntries {
		return 0
	}

	entry := self.profile.BATEntry(self.reader,
		int64(self.StartOffset+element_index*self.EntrySize))
	switch entry.State() {
	case PAYLOAD_BLOCK_FULLY_PRESENT, PAYLOAD_BLOCK_PARTIALLY_PRESENT:
		return entry.FileOffsetMB() * 1024 * 1024
	}

	return 0
}

func (self *Bat) DebugString() string {
	result := ""

	for i := uint64(0); i < self.TotalNumberOfEntries; i++ {
		b := self.profile.BATEntry(self.reader,
			int64(self.StartOffset+i*self.EntrySize))

		switch b.State() {
		case PAYLOAD_BLOCK_FULLY_PRESENT, PAYLOAD_BLOCK_PARTIALLY_PRESENT:
			result += b.DebugString()
		}
	}

	return result
}

type BatRange struct {
	FileOffset uint64
}

type BatReader struct {
	mu              sync.Mutex
	BlockSize       uint64
	EntriesPerChunk uint64
	Size            uint64
	Reader          io.ReaderAt

	bat *Bat
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
	block := int(uint64(off) / self.BlockSize)

	// The offset of the read within the block
	block_offset := uint64(off) % self.BlockSize

	to_read := int(self.BlockSize - block_offset)
	if to_read > len(buff) {
		to_read = len(buff)
	}

	block_file_offset := self.bat.GetFileOffset(block)
	if block_file_offset == 0 {
		// Block does not exist, just null pad the buffer.
		for i := 0; i < to_read; i++ {
			buff[i] = 0
		}
		return to_read, nil
	}

	// Get the reader to read the correct offset.
	return self.Reader.ReadAt(
		buff[:to_read], int64(block_file_offset+block_offset))
}

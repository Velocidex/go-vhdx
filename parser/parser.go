package parser

import (
	"errors"
	"fmt"
	"io"
)

type VHDXFile struct {
	// Reader to the underlying file
	reader  io.ReaderAt
	profile *VHDXProfile
	header  *FileType

	regions  []*RegionEntry
	bat      []*BATEntry
	Metadata *VHDXMetadata

	// Expose this reader to our callers: Reassemble the reader from
	// the BAT
	bat_reader *BatReader
}

func (self *VHDXFile) ReadAt(buff []byte, off int64) (int, error) {
	return self.bat_reader.ReadAt(buff, off)
}

func NewVHDXFile(reader io.ReaderAt) (*VHDXFile, error) {
	profile := NewVHDXProfile()
	self := &VHDXFile{
		reader:  reader,
		profile: profile,
		header:  profile.FileType(reader, 0),
		bat_reader: &BatReader{
			Reader: reader,
			bat:    make(map[int]*BatRange),
		},
	}

	err := self.Validate()
	if err != nil {
		return nil, err
	}

	region := self.header.Region1()
	self.regions = ParseArray_RegionEntry(profile, reader,
		region.Offset+int64(region.Size()), int(region.EntryCount()))

	for _, r := range self.regions {
		switch r.GUID() {
		case BAT_GUID:
			bat_entries := ParseArray_BATEntry(profile, reader,
				int64(r.FileOffset()), int(r.Length()/8))
			self.bat = append(self.bat, bat_entries...)

		case Metadata_GUID:
			metadata := profile.Metadata(reader, int64(r.FileOffset()))
			if metadata.Signature() != MetadataSignature {
				continue
			}

			self.Metadata = metadata.ParseMetadata()
		}
	}

	if self.Metadata == nil {
		return nil, errors.New("Unable to parse file metadata")
	}

	if self.Metadata.BlockSize <= 0 {
		return nil, fmt.Errorf("BlockSize invalid: %v", self.Metadata.BlockSize)
	}

	if self.Metadata.VirtualDiskSize <= 0 {
		return nil, fmt.Errorf("VirtualDiskSize invalid: %v",
			self.Metadata.VirtualDiskSize)
	}

	self.bat_reader.BlockSize = self.Metadata.BlockSize
	self.bat_reader.Size = self.Metadata.VirtualDiskSize

	// Now parse the BAT and build the reader.
	for i, b := range self.bat {
		switch b.State() {
		case PAYLOAD_BLOCK_FULLY_PRESENT, PAYLOAD_BLOCK_PARTIALLY_PRESENT:
			self.bat_reader.bat[i] = &BatRange{
				FileOffset: b.FileOffsetMB() * 1024 * 1024,
			}
		}
	}

	return self, nil
}

func (self *VHDXFile) DebugString() string {
	result := self.header.DebugString()
	for _, r := range self.regions {
		result += r.DebugString()
		result += fmt.Sprintf("GUID %02x\n", r.GUID())
	}

	for _, b := range self.bat {
		switch b.State() {
		case PAYLOAD_BLOCK_FULLY_PRESENT, PAYLOAD_BLOCK_PARTIALLY_PRESENT:
			result += b.DebugString()
		}
	}

	if self.Metadata != nil {
		result += fmt.Sprintf("Metadata: %#v\n", self.Metadata)
	}

	return result
}

func (self *VHDXFile) Validate() error {
	if self.header.Signature() != VHDXFileSignature {
		return errors.New("File should have the vhdxfile signature")
	}

	if self.header.Header1().Signature() != HeaderSignature ||
		self.header.Header2().Signature() != HeaderSignature {
		return errors.New("File header should have the head signature")
	}

	if self.header.Region1().Signature() != RegionSignature ||
		self.header.Region1().Signature() != RegionSignature {
		return errors.New("File region should have the regi signature")
	}

	return nil
}

package parser

type VHDXMetadata struct {
	BlockSize          uint64
	HasParent          bool
	VirtualDiskSize    uint64
	LogicalSectorSize  uint64
	PhysicalSectorSize uint32
	VirtualDiskId      string
}

func (self *Metadata) ParseMetadata() *VHDXMetadata {
	result := &VHDXMetadata{}

	entries := ParseArray_MetadataEntry(self.Profile, self.Reader,
		self.Offset+self.Profile.Off_Metadata_Entries,
		int(self.EntryCount()))

	for _, e := range entries {
		switch e.GUID() {
		case MetadataFileParameters:
			file_parameters := self.Profile.FileParameters(self.Reader,
				self.Offset+int64(e.MetadataOffset()))

			result.BlockSize = uint64(file_parameters.BlockSize())
			result.HasParent = file_parameters.HasParent() > 0

		case MetadataVirtualDiskSize:
			result.VirtualDiskSize = ParseUint64(self.Reader,
				self.Offset+int64(e.MetadataOffset()))

		case MetadataLogicalSectorSize:
			result.LogicalSectorSize = uint64(ParseUint32(self.Reader,
				self.Offset+int64(e.MetadataOffset())))

		case MetadataPhysicalSectorSize:
			result.PhysicalSectorSize = ParseUint32(self.Reader,
				self.Offset+int64(e.MetadataOffset()))

		case MetadataVirtualDiskId:
			guid := self.Profile.GUID(self.Reader,
				self.Offset+int64(e.MetadataOffset()))

			result.VirtualDiskId = guid.AsString()
		}
	}

	return result
}

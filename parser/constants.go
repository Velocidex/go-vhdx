package parser

const (
	VHDXFileSignature = "vhdxfile"
	HeaderSignature   = "head"
	RegionSignature   = "regi"
	MetadataSignature = "metadata"

	// Represents a GUID of 2DC27766-F623-4200-9D64-115E9BFD4A08
	BAT_GUID = "\x66\x77\xc2\x2d\x23\xf6\x00\x42\x9d\x64\x11\x5e\x9b\xfd\x4a\x08"

	// Represents a GUID of 8B7CA206-4790-4B9A-B8FE-575F050F886E
	Metadata_GUID = "\x06\xa2\x7c\x8b\x90\x47\x9a\x4b\xb8\xfe\x57\x5f\x05\x0f\x88\x6e"

	PAYLOAD_BLOCK_NOT_PRESENT       = 0
	PAYLOAD_BLOCK_UNDEFINED         = 1
	PAYLOAD_BLOCK_ZERO              = 2
	PAYLOAD_BLOCK_UNMAPPED          = 3
	PAYLOAD_BLOCK_FULLY_PRESENT     = 6
	PAYLOAD_BLOCK_PARTIALLY_PRESENT = 7

	// CAA16737-FA36-4D43-B3B6-33F0AA44E76B
	MetadataFileParameters = "\x37\x67\xa1\xca\x36\xfa\x43\x4d\xb3\xb6\x33\xf0\xaa\x44\xe7\x6b"

	// 2FA54224-CD1B-4876-B211-5DBED83BF4B8
	MetadataVirtualDiskSize = "\x24\x42\xa5\x2f\x1b\xcd\x76\x48\xb2\x11\x5d\xbe\xd8\x3b\xf4\xb8"

	// 8141BF1D-A96F-4709-BA47-F233A8FAAB5F
	MetadataLogicalSectorSize = "\x1d\xbf\x41\x81\x6f\xa9\x09\x47\xba\x47\xf2\x33\xa8\xfa\xab\x5f"

	// CDA348C7-445D-4471-9CC9-E9885251C556
	MetadataPhysicalSectorSize = "\xc7\x48\xa3\xcd\x5d\x44\x71\x44\x9c\xc9\xe9\x88\x52\x51\xc5\x56"

	// BECA12AB-B2E6-4523-93EF-C309E000C746
	MetadataVirtualDiskId = "\xab\x12\xca\xbe\xe6\xb2\x23\x45\x93\xef\xc3\x09\xe0\x00\xc7\x46"
)

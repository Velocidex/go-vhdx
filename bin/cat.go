package main

import (
	"os"

	"github.com/Velocidex/go-vhdx/parser"
	"github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	cat_command      = app.Command("cat", "Dump a VHDX file")
	cat_command_file = cat_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))
)

func doCat() {
	reader, _ := ntfs_parser.NewPagedReader(&ntfs_parser.OffsetReader{
		Reader: *cat_command_file,
	}, 1024, 10000)

	file_obj, err := parser.NewVHDXFile(reader)
	kingpin.FatalIfError(err, "Can not open file")

	buff := make([]byte, 1024*1024*10)
	offset := uint64(0)
	for offset < file_obj.Metadata.VirtualDiskSize {
		to_read := uint64(len(buff))
		if offset+to_read > file_obj.Metadata.VirtualDiskSize {
			to_read = file_obj.Metadata.VirtualDiskSize - offset
		}

		n, err := file_obj.ReadAt(buff, int64(offset))
		kingpin.FatalIfError(err, "Can not read file")

		os.Stdout.Write(buff[:n])
		offset += uint64(n)
	}
}
func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case cat_command.FullCommand():
			doCat()
		default:
			return false
		}
		return true
	})
}

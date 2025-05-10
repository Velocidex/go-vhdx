package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Velocidex/go-vhdx/parser"
	"github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	cmp_command      = app.Command("cmp", "Compare a VHDX file")
	cmp_command_file = cmp_command.Arg(
		"file", "The VHDX image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	cmp_command_file_original = cmp_command.Arg(
		"target", "The raw image file to compare with",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	cmp_command_size = cmp_command.Arg(
		"size", "The size of image to compare",
	).Required().Uint64()

	cmp_command_start = cmp_command.Flag(
		"start", "The start offset of image to compare",
	).Uint64()

	cmp_command_buffer_size = cmp_command.Flag(
		"buffers", "The size of the buffers to compare",
	).Default("10485760").Uint64()
)

func doCmp() {
	reader, err := ntfs_parser.NewPagedReader(&ntfs_parser.OffsetReader{
		Reader: *cmp_command_file,
	}, 1024, 10000)
	kingpin.FatalIfError(err, "Can not open file")

	file_obj, err := parser.NewVHDXFile(reader)
	kingpin.FatalIfError(err, "Can not open file")

	buff := make([]byte, *cmp_command_buffer_size)
	buff2 := make([]byte, *cmp_command_buffer_size)

	offset := uint64(*cmp_command_start)
	for offset < *cmp_command_size {
		fmt.Printf("Checking offset %v\n", offset)

		to_read := uint64(len(buff))
		if offset+to_read > *cmp_command_size {
			to_read = *cmp_command_size - offset
		}

		n, err := file_obj.ReadAt(buff, int64(offset))
		kingpin.FatalIfError(err, "Can not read file")

		n2, err := (*cmp_command_file_original).ReadAt(buff2, int64(offset))
		kingpin.FatalIfError(err, "Can not read file")

		if n != n2 {
			fmt.Printf("Offset %v: Unequal buffer reads %v vs %v\n",
				offset, n, n2)
			break
		}

		if bytes.Compare(buff[:n], buff2[:n]) != 0 {
			fmt.Printf("Offset %v: Error Reading Buffers\n", offset)
			break
		}

		offset += uint64(n)
	}
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case cmp_command.FullCommand():
			doCmp()
		default:
			return false
		}
		return true
	})
}

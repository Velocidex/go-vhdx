package main

import (
	"fmt"
	"os"

	"github.com/Velocidex/go-vhdx/parser"
	"github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	parse_command      = app.Command("parse", "Parse a VHDX file")
	parse_command_file = parse_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))
)

func doParse() {
	reader, _ := ntfs_parser.NewPagedReader(&ntfs_parser.OffsetReader{
		Reader: *parse_command_file,
	}, 1024, 10000)

	file_obj, err := parser.NewVHDXFile(reader)
	kingpin.FatalIfError(err, "Can not open file")

	fmt.Println(file_obj.DebugString())
}
func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case parse_command.FullCommand():
			doParse()
		default:
			return false
		}
		return true
	})
}

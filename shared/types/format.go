package types

import (
	"errors"
	"nexa/shared/util"
	"strings"
)

func GetFileFormat(filename string) (Format, error) {
	// trim folders
	index := strings.LastIndexByte(filename, '/')
	if index == -1 {
		index = strings.LastIndexByte(filename, '\\')
	}
	if index > 0 {
		filename = filename[index:]
	}

	// get last dot
	index = strings.LastIndexByte(filename, '.')
	if index == -1 {
		return Format(""), ErrFileHasNoFormat
	}
	return Format(filename[index+1:]), nil
}

type Format string

func (f Format) Validate() error {
	return util.Ternary(f.Underlying() == "", ErrUnknownFormat, nil)
}

func (f Format) Underlying() string {
	return string(f)
}

var ErrFileHasNoFormat = errors.New("file has no format on its name, prefer use this name:'filename.format'")
var ErrUnknownFormat = errors.New("file has unknown format")

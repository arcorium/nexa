package util

import (
	"mime"
	"path/filepath"
)

func GetMimeType(filePath string) string {
	ext := mime.TypeByExtension(filepath.Ext(filePath))
	if ext == "" {
		return "application/octet-stream"
	}
	return ext
}

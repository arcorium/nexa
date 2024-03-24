package types

import (
	"fmt"
	"strings"
)

type FilePath string

func FilePathFromString(path string) FilePath {
	return FilePath(strings.TrimPrefix(path, "/"))
}

func (p *FilePath) FileName() string {
	return string(*p)
}

func (p *FilePath) FullPath(staticEndpoint string) string {
	staticEndpoint = strings.TrimSuffix(strings.TrimSuffix(staticEndpoint, "/"), "\\")
	return fmt.Sprintf("%s/%s", staticEndpoint, p.FileName())
}

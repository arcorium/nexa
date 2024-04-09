package types

import (
  "fmt"
  "strings"
)

func FilePathFromString(path string) FilePath {
  return FilePath(strings.TrimPrefix(path, "/"))
}

type FilePath string

func (p *FilePath) FileName() string {
  return string(*p)
}

func (p *FilePath) FullPath(staticEndpoint string) string {
  staticEndpoint = strings.TrimSuffix(strings.TrimSuffix(staticEndpoint, "/"), "\\")
  return fmt.Sprintf("%s/%s", staticEndpoint, p.FileName())
}

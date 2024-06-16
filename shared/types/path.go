package types

import (
	"fmt"
	"net/url"
	"strings"
)

func FilePathFromString(path string) FilePath {
	return FilePath(path)
}

func FilePathFromURL(url *url.URL) FilePath {
	return FilePath(fmt.Sprintf("%s%s", url.Host, url.Path))
}

type FilePath string

func (p FilePath) Underlying() string {
	return string(p)
}

func (p FilePath) FileName() string {
	filename := strings.TrimPrefix(p.Underlying(), "/")
	return filename
}

func (p FilePath) Path() string {
	return p.Underlying()
}

//
//func (p *FilePath) FullPath(staticEndpoint string) string {
//  staticEndpoint = strings.TrimSuffix(strings.TrimSuffix(staticEndpoint, "/"), "\\")
//  return fmt.Sprintf("%s/%s", staticEndpoint, p.FileName())
//}
//
//func (p *FilePath) RelativePath() string {
//  return p.Underlying()
//}

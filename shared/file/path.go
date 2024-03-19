package file

import (
	"fmt"
	"strings"
)

type RemotePath string

func (p *RemotePath) FileName() string {
	return string(*p)
}

func (p *RemotePath) FullPath(staticEndpoint string) string {
	staticEndpoint = strings.TrimSuffix(strings.TrimSuffix(staticEndpoint, "/"), "\\")
	return fmt.Sprintf("%s/%s", staticEndpoint, p.FileName())
}

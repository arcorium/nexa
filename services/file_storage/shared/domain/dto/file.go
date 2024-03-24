package dto

type FileSaveInput struct {
	Name  string `json:"name"`
	Bytes []byte `json:"bytes"`
}

type FileLoadInput struct {
	Name string
}

type FileDeleteInput struct {
	UseBin bool
	Name   string
}

type FileReplaceInput struct {
	Name  string
	Bytes []byte
}

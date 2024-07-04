package entity

type File struct {
  Name     string
  Bytes    []byte
  Size     uint64
  IsPublic bool
}

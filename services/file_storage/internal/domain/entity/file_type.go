package entity

type FileType uint8

const (
  FileTypeImage FileType = iota
  FileTypeDocument
  FileTypeVideo
  FileTypeSound
  FileTypeOther
)

func (f FileType) Underlying() uint8 {
  return uint8(f)
}

func (f FileType) String() string {
  switch f {
  case FileTypeImage:
    return "Image"
  case FileTypeDocument:
    return "Document"
  case FileTypeVideo:
    return "Video"
  case FileTypeSound:
    return "Sound"
  case FileTypeOther:
    return "Other"
  }
  return "Unknown"
}

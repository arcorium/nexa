package entity

type Status uint8

const (
  StatusPending Status = iota
  StatusSending
  StatusDelivered
  StatusFailed
)

func (s Status) Underlying() uint8 {
  return uint8(s)
}

func (s Status) String() string {
  switch s {
  case StatusPending:
    return "Pending"
  case StatusSending:
    return "Sending"
  case StatusDelivered:
    return "Delivered"
  case StatusFailed:
    return "Failed"
  }
  return "Unknown"
}

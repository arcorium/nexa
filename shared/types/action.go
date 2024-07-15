package types

type Action string

func (a Action) String() string {
  return string(a)
}

func (a Action) Underlying() string {
  return string(a)
}

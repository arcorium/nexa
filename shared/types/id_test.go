package types

import (
  "testing"
)

func TestIdFromString(t *testing.T) {
  type args struct {
    id string
  }
  tests := []struct {
    name    string
    args    args
    want    string
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        id: "ba474c5a-c1bf-43b1-96b6-d0225def9361",
      },
      want:    "ba474c5a-c1bf-43b1-96b6-d0225def9361",
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      got, err := IdFromString(tt.args.id)
      if (err != nil) != tt.wantErr {
        t.Errorf("IdFromString() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      if got.Underlying().String() != tt.want {
        t.Errorf("IdFromString() got = %v, want %v", got, tt.want)
      }
      //if !reflect.DeepEqual(got, tt.want) {
      //  t.Errorf("IdFromString() got = %v, want %v", got, tt.want)
      //}
    })
  }
}

func TestId_Equal(t *testing.T) {
  type args struct {
    uuid string
  }
  tests := []struct {
    name string
    i    Id
    args args
    want bool
  }{
    {
      name: "EqWithString",
      i:    DropError(IdFromString("ba474c5a-c1bf-43b1-96b6-d0225def9361")),
      args: args{
        uuid: "ba474c5a-c1bf-43b1-96b6-d0225def9361",
      },
      want: true,
    },
    {
      name: "Not EqWithString",
      i:    DropError(NewId()),
      args: args{
        uuid: "ba474c5a-c1bf-43b1-96b6-d0225def9361",
      },
      want: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if got := tt.i.EqWithString(tt.args.uuid); got != tt.want {
        t.Errorf("EqWithString() = %v, want %v", got, tt.want)
      }
    })
  }
}

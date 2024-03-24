package util

import (
	"math/rand"
	"strings"
)

func RandomString(length uint64) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	builder := strings.Builder{}
	builder.Grow(int(length))

	for i := 0; i < int(length); i++ {
		builder.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}

	return builder.String()
}

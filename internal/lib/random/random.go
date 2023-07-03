package random

import (
	"golang.org/x/exp/rand"
	"time"
)

// TODO: write test

func NewRandomString(size int) string {
	rnd := rand.New(
		rand.NewSource(
			uint64(time.Now().UnixNano()),
		),
	)

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := 0; i < size; i++ {
		b[i] = chars[rnd.Intn(62)]
	}

	return string(b)
}

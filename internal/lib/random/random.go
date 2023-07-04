package random

import (
	"time"

	"golang.org/x/exp/rand"
)

// TODO: add: test
func NewRandomString(size int) string {
	rnd := rand.New(
		rand.NewSource(
			uint64(time.Now().UnixNano()),
		),
	)

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	charsSize := len(chars)
	b := make([]rune, size)

	for i := 0; i < size; i++ {
		b[i] = chars[rnd.Intn(charsSize)]
	}

	return string(b)
}

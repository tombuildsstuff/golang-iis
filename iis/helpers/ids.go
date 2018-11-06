package helpers

import (
	"math/rand"
	"time"
)

func RandomInt() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

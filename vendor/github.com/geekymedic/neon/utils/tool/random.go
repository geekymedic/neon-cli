package tool

import (
	"math/rand"
	"time"
)

var timeRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RangeBitsInt(low, hi int) int {
	if low > hi {
		panic("low must be less or equal hi")
	}
	return low + timeRand.Intn(hi-low)
}

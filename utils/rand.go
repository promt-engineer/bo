package utils

import (
	"math/rand"
	"time"
)

var (
	rnd *rand.Rand
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Rand63() int64 {
	return rnd.Int63()
}

func Rand64() uint64 {
	return rnd.Uint64()
}

func Rand63n(n int64) int64 {
	if n == 0 {
		return Rand63()
	}

	return Rand63() % n
}

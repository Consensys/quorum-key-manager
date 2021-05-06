package common

import (
	"math/rand"
	"time"
)

func seedRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
func RandString(n int) string {
	var seededRand = seedRand()
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[seededRand.Intn(len(letter))]
	}
	return string(b)
}

func RandHexString(n int) string {
	var seededRand = seedRand()
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandInt(n int) int {
	var seededRand = seedRand()
	return seededRand.Intn(n)
}

func RandIntRange(min, max int) int {
	var seededRand = seedRand()
	return seededRand.Intn(max-min) + min
}

func RandBool() bool {
	var seededRand = seedRand()
	return seededRand.Intn(2) != 0
}

func RandShuffle(arr []string) []string {
	var seededRand = seedRand()
	out := append([]string{}, arr...)
	seededRand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

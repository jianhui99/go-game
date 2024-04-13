package utils

import (
	"math/rand"
	"time"
)

func Contains[T int | string](data []T, value T) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func Rand(n int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(n)
}

func Default(v, d string) string {
	if len(v) == 0 {
		return d
	}
	return v
}

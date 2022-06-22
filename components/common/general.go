package common

import (
	"log"
	"math/rand"
)

func Keys[K comparable, V any](m map[K]V) []K {
	keys := []K{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func SelectRandom[T string](list []T, amount int) []T {
	pool := make([]T, len(list))
	copy(pool, list)
	if amount > len(pool) {
		log.Fatal("SelectRandom cannot choose more than the slice has to offer!")
	}
	values := []T{}
	for i := 0; i < amount; i++ {
		index := rand.Intn(len(pool))
		values = append(values, pool[index])
		pool = append(pool[0:index], pool[index+1:]...)
	}
	return values
}

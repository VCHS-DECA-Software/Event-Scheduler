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

func SelectRandom[T any](list []T, amount int) []T {
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

func Without[T comparable](list []T, element T) []T {
	pool := make([]T, len(list))
	copy(pool, list)
	for i, e := range pool {
		if e == element {
			return append(pool[0:i], pool[i+1:]...)
		}
	}
	return pool
}

func Intersects[T comparable](l1 []T, l2 []T) bool {
	pool := map[T]bool{}
	for _, e := range l1 {
		pool[e] = true
	}
	for _, e := range l2 {
		if pool[e] {
			return true
		}
	}
	return false
}

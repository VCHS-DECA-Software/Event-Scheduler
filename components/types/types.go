package types

import (
	"github.com/fxtlabs/primes"
)

type Type struct {
	Type     int
	Category int
}

var offsets []int
var primeList []int
var rPrimeList []int

func Register(categories [][]any) {
	offsets = []int{}
	offset := 0
	for i, c := range categories {
		offsets[i] = offset
		offset += len(c)
	}

	primeList = primes.Sieve(offset)
	rPrimeList = []int{}
	for i, p := range primeList {
		rPrimeList[p] = i
	}
}

//Composite works by using the nth prime of each type integer
//then multiplying them together, the resulting number will
//only be divisible by it's factors, the two type integers.
func Composite(t1, t2 Type) int {
	a := offsets[t1.Type] + t1.Type
	b := offsets[t2.Type] + t2.Type
	return primeList[a] * primeList[b]
}

//Decompose is the inverse operation of composite
func Decompose(c int, t Type) Type {
	n := rPrimeList[c/t.Type]

	var category int
	var inner int
	for i, offset := range offsets {
		if n > offset {
			inner = offset - offsets[i]
			category = i
			break
		}
	}

	return Type{
		Type:     inner,
		Category: category,
	}
}

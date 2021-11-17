package main

import(
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"math/big"
)

type SimHash struct {
	Sum []byte
}

func NewSimHash() SimHash {
	return SimHash{
		Sum: []byte{},
	}
}

// This could also return []bool
func asBits(val uint64) []uint64 {
	bits := []uint64{}
	for i := 0; i < 24; i++ {
		bits = append([]uint64{val & 0x1}, bits...)
		// or
		// bits = append(bits, val & 0x1)
		// depending on the order you want
		val = val >> 1
	}
	return bits
}

func (s *SimHash) hash(data []byte) []uint64 {
	hash := sha1.Sum(data)
	rawHex := hex.EncodeToString(hash[:])
	i := new(big.Int)
	i.SetString(rawHex, 16)
	return asBits(i.Uint64())
}

func (s *SimHash) table(data map[string]int) [][]int {
	array := [][]int{}
	for k, val := range data {
		part := []int{}
		for _, v := range s.hash([]byte(k)) {
			if v == 0 {
				part = append(part, -1 * val)
			}else {
				part = append(part, 1 * val)
			}
		}
		array = append(array, part)
	}
	return array
}

func (s *SimHash) convert(table [][]int){
	sum := []byte{}
	for i := 0; i < len(table[0]); i++ {
		colsum := 0
		for j := 0; j < len(table); j++ {
			colsum = colsum + table[j][i]
		}
		if colsum >= 1 {
			sum = append(sum, 1)
		}else if colsum <= 0 {
			sum = append(sum, 0)
		}
	}
	s.Sum = append(s.Sum, sum...)
}

func (s *SimHash) Process(data map[string]int){
	s.convert(s.table(data))
}

func (s *SimHash) Distance(s1 SimHash) int {
	dist := 0
	for i := 0; i < len(s.Sum); i++ {
		if s.Sum[i] ^ s1.Sum[i] == 1{
			dist++
		}
	}
	return dist
}

func main(){
	data1 := map[string]int {
		"hello": 1,
		"world":2,
	}

	data2 := map[string]int {
		"hello": 1,
		"world": 1,
		"yes": 1,
	}

	s := NewSimHash()
	s.Process(data1)

	s1 := NewSimHash()
	s1.Process(data2)

	fmt.Println(s.Distance(s1))
}

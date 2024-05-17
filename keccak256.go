// Segments copied from go-ethereum/crypto/crypto.go

package mist

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

type KeccakState interface {
	hash.Hash
	Read(p []byte) (int, error)
}

const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

type Hash [HashLength]byte

func Keccak256Hash(data ...[]byte) Hash {
	d, ok := sha3.NewLegacyKeccak256().(KeccakState)
	if !ok {
		panic("TODO")
	}

	for _, b := range data {
		d.Write(b)
	}

	var h Hash
	d.Read(h[:])
	return h
}

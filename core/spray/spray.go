package spray

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/klauspost/reedsolomon"
)

// DISCLAIMER: The following code is AI generated and should be used with discretion or as a STUB.

// Encoder holds the configuration for encoding and decoding operations.
type Encoder struct {
	dataShards   int
	parityShards int
	codec        reedsolomon.Encoder
}

// NewEncoder initializes a new encoder with the specified number of data and parity shards.
// The total number of shards will be dataShards + parityShards.
func NewEncoder(dataShards, parityShards int) (*Encoder, error) {
	codec, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reed-Solomon codec: %w", err)
	}
	return &Encoder{
		dataShards:   dataShards,
		parityShards: parityShards,
		codec:        codec,
	}, nil
}

// Encode takes a slice of data shards (e.g., image segments) and a seed.
// It returns a slice of parity shards generated from the data.
// The seed is used to ensure the parity generation is deterministic.
func (e *Encoder) Encode(dataShards [][]byte, seed []byte) ([][]byte, error) {
	if len(dataShards) != e.dataShards {
		return nil, fmt.Errorf("incorrect number of data shards: expected %d, got %d", e.dataShards, len(dataShards))
	}

	// The seed is combined with shard data to produce deterministic parity.
	// This step is illustrative; in a real system, the seed might determine
	// *which* parity algorithm or parameters to use. Here we use it to
	// slightly modify the input, ensuring a seeded outcome.
	hasher := sha256.New()
	seededDataShards := make([][]byte, e.dataShards)
	for i, shard := range dataShards {
		seededDataShards[i] = e.applySeed(shard, seed, uint64(i), hasher)
	}

	parityShards := make([][]byte, e.parityShards)
	for i := range parityShards {
		parityShards[i] = make([]byte, len(seededDataShards[0]))
	}

	// Create the full shard slice and generate parity
	allShards := append(seededDataShards, parityShards...)
	if err := e.codec.Encode(allShards); err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}

	return allShards[e.dataShards:], nil
}

// Decode reconstructs the original data from a mix of data and parity shards.
// Shards can be missing (represented by nil slices).
// As long as `dataShards` number of shards are present, the data can be recovered.
func (e *Encoder) Decode(allShards [][]byte) ([][]byte, error) {
	if len(allShards) != e.dataShards+e.parityShards {
		return nil, fmt.Errorf("incorrect total number of shards: expected %d, got %d", e.dataShards+e.parityShards, len(allShards))
	}

	// Verify and reconstruct the data
	if err := e.codec.Reconstruct(allShards); err != nil {
		// If reconstruction fails, it means there isn't enough data.
		return nil, fmt.Errorf("failed to reconstruct data (not enough shards?): %w", err)
	}

	// The seed is not needed for decoding because the parity data itself
	// contains the necessary recovery information.

	return allShards[:e.dataShards], nil
}

// applySeed is a helper to combine a shard's data with a seed and index.
// This is a simple way to make the encoding process dependent on the seed.
func (e *Encoder) applySeed(data, seed []byte, index uint64, hasher hash.Hash) []byte {
	hasher.Reset()
	hasher.Write(seed)
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, index)
	hasher.Write(indexBytes)
	seedHash := hasher.Sum(nil)

	buf := bytes.NewBuffer(nil)
	for i, b := range data {
		buf.WriteByte(b ^ seedHash[i%len(seedHash)])
	}
	return buf.Bytes()
}

// Package hydration provides functions to generate hierarchical seeds for images and their segments.
package hydration

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// GenerateImageSeed creates a cryptographically strong seed from an input string,
// such as an image name or unique identifier.
// It returns a 32-byte hash that can serve as the primary seed.
func GenerateImageSeed(imageIdentifier string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(imageIdentifier))
	return hasher.Sum(nil)
}

// DeriveNanoseeds takes the primary image seed and the number of segments
// and generates a unique, deterministic nanoseed for each segment.
// It uses a hashing function to ensure that each nanoseed is derived from the
// original seed but is unique for each segment index.
func DeriveNanoseeds(imageSeed []byte, segmentCount int) ([][]byte, error) {
	if len(imageSeed) == 0 {
		return nil, fmt.Errorf("image seed cannot be empty")
	}
	if segmentCount <= 0 {
		return nil, fmt.Errorf("segment count must be a positive number")
	}

	nanoseeds := make([][]byte, segmentCount)
	hasher := sha256.New()

	for i := 0; i < segmentCount; i++ {
		hasher.Reset()
		hasher.Write(imageSeed)
		indexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(indexBytes, uint64(i))
		hasher.Write(indexBytes)

		nanoseeds[i] = hasher.Sum(nil)
	}

	return nanoseeds, nil
}

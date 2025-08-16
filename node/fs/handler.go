package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var dataDir string

func init() {
	exePath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("could not get executable path: %v", err))
	}

	dataDir = filepath.Join(filepath.Dir(exePath), "data")
}

// StoreSegment stores a segment's data on the node's filesystem.
// The incoming nanoseed and segmentData must contain non-empty data.
// If successful, the segment is stored in the data directory with the
// nanoseed's hexadecimal representation as the filename.
func StoreSegment(nanoseed []byte, segmentData []byte) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("could not create data directory at %s: %w", dataDir, err)
	}

	filename := fmt.Sprintf("%x", nanoseed)
	filePath := filepath.Join(dataDir, filename)

	log.Printf("Attempting to store segment at: %s", filePath)

	if err := os.WriteFile(filePath, segmentData, 0644); err != nil {
		return fmt.Errorf("failed to write segment to file %s: %w", filePath, err)
	}

	return nil
}

// RetrieveSegment retrieves a segment's data from the node's filesystem.
// The incoming nanoseed must contain a non-empty value. If successful, the segment's data is returned.
// If the segment is not found, a nil slice and an error indicating that the segment was not found is returned.
// If the segment is found but there is an error reading the file, a nil slice and an error indicating the failure to read the file is returned.
func RetrieveSegment(nanoseed []byte) ([]byte, error) {
	filename := fmt.Sprintf("%x", nanoseed)
	filePath := filepath.Join(dataDir, filename)

	log.Printf("Attempting to retrieve segment from: %s", filePath)

	segmentData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("segment not found at %s", filePath)
		}
		return nil, fmt.Errorf("failed to read segment from file %s: %w", filePath, err)
	}

	return segmentData, nil
}

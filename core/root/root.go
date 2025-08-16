package root

import (
	"context"
	"fmt"

	"imago-tree/cutter"
	"imago-tree/hydration"
	"imago-tree/spray"
)

const (
	DataShards   = 4
	ParityShards = 2
)

type Server struct {
	UnimplementedRootServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// CreateImageChunks takes an image in raw byte form and its identifier, cuts the image into horizontal segments,
// generates a seed for the image, derives nanoseeds for each segment, encodes the segments using erasure coding,
// and returns the resulting image chunks and parity shards.
func (s *Server) CreateImageChunks(ctx context.Context, req *CreateImageChunksRequest) (*ImageChunks, error) {
	imageData := req.GetImageData()
	imageIdentifier := req.GetImageIdentifier()

	segments, err := cutter.CutImageBytes(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to cut image: %w", err)
	}

	// TODO: Implement a more graceful way to handle this in the future.
	if len(segments)%DataShards != 0 {
		return nil, fmt.Errorf("number of segments (%d) must be a multiple of data shards (%d) for this example", len(segments), DataShards)
	}

	imageSeed := hydration.GenerateImageSeed(imageIdentifier)

	nanoseeds, err := hydration.DeriveNanoseeds(imageSeed, len(segments))
	if err != nil {
		return nil, fmt.Errorf("failed to derive nanoseeds: %w", err)
	}

	encoder, err := spray.NewEncoder(DataShards, ParityShards)
	if err != nil {
		return nil, fmt.Errorf("failed to create erasure coding encoder: %w", err)
	}

	allDataShards := make([][]byte, len(segments))
	responseChunks := make([]*Chunk, len(segments))

	for i, segmentImg := range segments {
		segmentBytes, err := cutter.ImageToPNG(segmentImg)
		if err != nil {
			return nil, fmt.Errorf("failed to encode segment %d to PNG: %w", i, err)
		}

		allDataShards[i] = segmentBytes
		responseChunks[i] = &Chunk{
			Nanoseed:    nanoseeds[i],
			SegmentData: segmentBytes,
		}
	}

	allParityShards := make([][]byte, 0)
	for i := 0; i < len(allDataShards); i += DataShards {
		batch := allDataShards[i : i+DataShards]

		parity, err := encoder.Encode(batch, imageSeed)
		if err != nil {
			return nil, fmt.Errorf("failed to generate parity for batch starting at index %d: %w", i, err)
		}
		allParityShards = append(allParityShards, parity...)
	}

	response := &ImageChunks{
		Chunks:       responseChunks,
		ParityShards: allParityShards,
	}

	return response, nil
}

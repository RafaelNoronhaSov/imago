package server

import (
	"context"
	"fmt"

	"imago-node/fs"
)

type Server struct {
	UnimplementedNodeServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// StoreSegment stores a segment's data on the node's filesystem.
// The incoming StoreSegmentRequest must contain a non-empty nanoseed and
// segment_data. If successful, a StoreSegmentResponse is returned with a
// message indicating the success of the operation.
func (s *Server) StoreSegment(ctx context.Context, req *StoreSegmentRequest) (*StoreSegmentResponse, error) {
	nanoseed := req.GetNanoseed()
	segmentData := req.GetSegmentData()

	if len(nanoseed) == 0 || len(segmentData) == 0 {
		return nil, fmt.Errorf("nanoseed and segment_data cannot be empty")
	}

	if err := fs.StoreSegment(nanoseed, segmentData); err != nil {
		return nil, fmt.Errorf("failed to store segment: %w", err)
	}

	response := &StoreSegmentResponse{
		Message: fmt.Sprintf("Successfully stored segment %x", nanoseed),
	}

	return response, nil
}

// RetrieveSegment retrieves a segment's data from the node's filesystem.
// The incoming RetrieveSegmentRequest must contain a non-empty nanoseed.
// If successful, a RetrieveSegmentResponse is returned with the segment's data.
func (s *Server) RetrieveSegment(ctx context.Context, req *RetrieveSegmentRequest) (*RetrieveSegmentResponse, error) {
	nanoseed := req.GetNanoseed()
	if len(nanoseed) == 0 {
		return nil, fmt.Errorf("nanoseed cannot be empty")
	}

	segmentData, err := fs.RetrieveSegment(nanoseed)
	if err != nil {
		return nil, err
	}

	response := &RetrieveSegmentResponse{
		SegmentData: segmentData,
	}

	return response, nil
}

package application

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

// ClamAVScanner inspects media binary streams against virus signatures (ADR-012, D-05-02-001).
type ClamAVScanner interface {
	ScanStream(ctx context.Context, reader io.Reader) error
}

// MediaStorage interface abstracts SeaweedFS S3 storage operations.
type MediaStorage interface {
	Upload(ctx context.Context, path string, data []byte, contentType string) (string, error)
}

type UploadPhotoInput struct {
	OwnerID     uuid.UUID
	Angle       domain.PhotoAngle
	ContentType string
	Reader      io.Reader
}

type UploadPhotoOutput struct {
	StorageURI string
	FileHash   string
	Angle      domain.PhotoAngle
}

type UploadPhotoService struct {
	scanner ClamAVScanner
	storage MediaStorage
}

func NewUploadPhotoService(scanner ClamAVScanner, storage MediaStorage) *UploadPhotoService {
	return &UploadPhotoService{
		scanner: scanner,
		storage: storage,
	}
}

// UploadAndScanPhoto streams the image through ClamAV and commits to SeaweedFS S3 if clean.
func (s *UploadPhotoService) UploadAndScanPhoto(ctx context.Context, input UploadPhotoInput) (*UploadPhotoOutput, error) {
	// Read payload into memory buffer for hashing and streaming
	buf := new(bytes.Buffer)
	hasher := sha256.New()
	multiWriter := io.MultiWriter(buf, hasher)

	if _, err := io.Copy(multiWriter, input.Reader); err != nil {
		return nil, fmt.Errorf("failed to read photo payload: %w", err)
	}

	rawBytes := buf.Bytes()
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// 1. Synchronous ClamAV Stream Inspection (:3310 zINSTREAM)
	if err := s.scanner.ScanStream(ctx, bytes.NewReader(rawBytes)); err != nil {
		return nil, domain.ErrMalwareDetected
	}

	// 2. Commit Clean Payload to SeaweedFS S3
	storageKey := fmt.Sprintf("fleets/%s/%s_%s.webp", input.OwnerID.String(), strings.ToLower(string(input.Angle)), fileHash[:12])
	storageURI, err := s.storage.Upload(ctx, storageKey, rawBytes, input.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to commit clean photo to media storage: %w", err)
	}

	return &UploadPhotoOutput{
		StorageURI: storageURI,
		FileHash:   fileHash,
		Angle:      input.Angle,
	}, nil
}

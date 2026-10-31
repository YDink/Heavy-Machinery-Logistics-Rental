package application_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

// MockClamAVScanner detects the standard EICAR test signature.
type MockClamAVScanner struct{}

func (m *MockClamAVScanner) ScanStream(ctx context.Context, reader io.Reader) error {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)
	data := buf.String()

	// EICAR Standard Antivirus Test String
	if strings.Contains(data, "EICAR-STANDARD-ANTIVIRUS-TEST-FILE") {
		return domain.ErrMalwareDetected
	}
	return nil
}

type MockMediaStorage struct {
	Files map[string][]byte
}

func (m *MockMediaStorage) Upload(ctx context.Context, path string, data []byte, contentType string) (string, error) {
	if m.Files == nil {
		m.Files = make(map[string][]byte)
	}
	m.Files[path] = data
	return "s3://dink-fleet-media/" + path, nil
}

func TestUploadPhoto_MalwareRejection(t *testing.T) {
	scanner := &MockClamAVScanner{}
	storage := &MockMediaStorage{}
	service := application.NewUploadPhotoService(scanner, storage)

	// EICAR virus simulation payload
	eicarPayload := []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)

	input := application.UploadPhotoInput{
		OwnerID:     uuid.New(),
		Angle:       domain.PhotoAngleFront,
		ContentType: "image/webp",
		Reader:      bytes.NewReader(eicarPayload),
	}

	_, err := service.UploadAndScanPhoto(context.Background(), input)
	if err != domain.ErrMalwareDetected {
		t.Fatalf("expected ErrMalwareDetected, got %v", err)
	}

	// Verify that the infected file was never committed to storage
	if len(storage.Files) > 0 {
		t.Errorf("expected storage to have 0 files, got %d", len(storage.Files))
	}
}

func TestUploadPhoto_CleanUpload(t *testing.T) {
	scanner := &MockClamAVScanner{}
	storage := &MockMediaStorage{}
	service := application.NewUploadPhotoService(scanner, storage)

	cleanPayload := []byte("clean-image-binary-stream-sample")

	input := application.UploadPhotoInput{
		OwnerID:     uuid.New(),
		Angle:       domain.PhotoAngleFront,
		ContentType: "image/webp",
		Reader:      bytes.NewReader(cleanPayload),
	}

	out, err := service.UploadAndScanPhoto(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error on clean upload: %v", err)
	}

	if out.StorageURI == "" {
		t.Error("expected non-empty storage URI")
	}
	if out.FileHash == "" {
		t.Error("expected non-empty SHA-256 hash")
	}
	if len(storage.Files) != 1 {
		t.Errorf("expected 1 file stored, got %d", len(storage.Files))
	}
}

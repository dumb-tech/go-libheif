package libheif

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

// createTestHeifFile generates a minimal HEIF file for testing by encoding
// a 2x2 NRGBA image via EncodeImageAsHeif.
func createTestHeifFile(t *testing.T) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 0, A: 255})

	tmpDir := t.TempDir()
	heifPath := filepath.Join(tmpDir, "test.heic")

	if err := EncodeImageAsHeif(img, heifPath); err != nil {
		t.Fatalf("failed to create test HEIF file: %v", err)
	}
	return heifPath
}

// createNonHeifFile generates a temp .txt file with garbage bytes for testing
// non-HEIF input rejection.
func createNonHeifFile(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	txtPath := filepath.Join(tmpDir, "not-heif.txt")
	if err := os.WriteFile(txtPath, []byte("this is not a heif file"), 0644); err != nil {
		t.Fatalf("failed to create non-HEIF file: %v", err)
	}
	return txtPath
}

func TestDecodeHeifToJpeg(t *testing.T) {
	tests := []struct {
		name      string
		heifPath  string
		jpegPath  string
		quality   int
		wantErr   bool
		errContains string
		setup     func(t *testing.T) string // returns heifPath if dynamic
	}{
		{
			name:    "valid HEIF to JPEG at quality 80",
			quality: 80,
			setup:   func(t *testing.T) string { return createTestHeifFile(t) },
		},
		{
			name:        "quality 0 returns error",
			heifPath:    "dummy.heic",
			jpegPath:    "out.jpg",
			quality:     0,
			wantErr:     true,
			errContains: "quality must be between 1 and 100, got: 0",
		},
		{
			name:        "quality 101 returns error",
			heifPath:    "dummy.heic",
			jpegPath:    "out.jpg",
			quality:     101,
			wantErr:     true,
			errContains: "quality must be between 1 and 100, got: 101",
		},
		{
			name:        "empty heifPath returns error",
			heifPath:    "",
			jpegPath:    "out.jpg",
			quality:     80,
			wantErr:     true,
			errContains: "heif source path is empty",
		},
		{
			name:        "empty jpegPath returns error",
			heifPath:    "dummy.heic",
			jpegPath:    "",
			quality:     80,
			wantErr:     true,
			errContains: "jpeg output path is empty",
		},
		{
			name:        "non-HEIF input returns error",
			quality:     80,
			wantErr:     true,
			errContains: "", // just check error is non-nil
			setup:       func(t *testing.T) string { return createNonHeifFile(t) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heifPath := tt.heifPath
			jpegPath := tt.jpegPath

			if tt.setup != nil {
				heifPath = tt.setup(t)
				if jpegPath == "" && !tt.wantErr {
					jpegPath = filepath.Join(t.TempDir(), "output.jpg")
				}
				if jpegPath == "" && tt.wantErr {
					jpegPath = filepath.Join(t.TempDir(), "output.jpg")
				}
			}

			err := DecodeHeifToJpeg(heifPath, jpegPath, tt.quality)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify output file exists and has content
			info, err := os.Stat(jpegPath)
			if err != nil {
				t.Fatalf("output JPEG file does not exist: %v", err)
			}
			if info.Size() == 0 {
				t.Fatal("output JPEG file is empty")
			}
		})
	}
}

func TestDecodeHeifToPng(t *testing.T) {
	tests := []struct {
		name        string
		heifPath    string
		pngPath     string
		wantErr     bool
		errContains string
		setup       func(t *testing.T) string
	}{
		{
			name:  "valid HEIF to PNG",
			setup: func(t *testing.T) string { return createTestHeifFile(t) },
		},
		{
			name:        "empty heifPath returns error",
			heifPath:    "",
			pngPath:     "out.png",
			wantErr:     true,
			errContains: "heif source path is empty",
		},
		{
			name:        "empty pngPath returns error",
			heifPath:    "dummy.heic",
			pngPath:     "",
			wantErr:     true,
			errContains: "png output path is empty",
		},
		{
			name:        "non-HEIF input returns error",
			wantErr:     true,
			errContains: "",
			setup:       func(t *testing.T) string { return createNonHeifFile(t) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heifPath := tt.heifPath
			pngPath := tt.pngPath

			if tt.setup != nil {
				heifPath = tt.setup(t)
				if pngPath == "" {
					pngPath = filepath.Join(t.TempDir(), "output.png")
				}
			}

			err := DecodeHeifToPng(heifPath, pngPath)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify output file exists and has content
			info, err := os.Stat(pngPath)
			if err != nil {
				t.Fatalf("output PNG file does not exist: %v", err)
			}
			if info.Size() == 0 {
				t.Fatal("output PNG file is empty")
			}
		})
	}
}

// containsString is a simple helper to check substring presence.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (substr == "" || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

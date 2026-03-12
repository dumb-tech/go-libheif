package libheif

import (
	"image"
	"os"
	"testing"
)

func TestEncodeImageAsHeif(t *testing.T) {
	validImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	tests := []struct {
		name        string
		img         image.Image
		path        func() string
		opts        []EncoderOption
		wantErr     bool
		errContains []string
	}{
		{
			name:        "nil image returns error",
			img:         nil,
			path:        func() string { return "out.heic" },
			wantErr:     true,
			errContains: []string{"image is nil"},
		},
		{
			name:        "empty path returns error",
			img:         validImg,
			path:        func() string { return "" },
			wantErr:     true,
			errContains: []string{"path", "empty"},
		},
		{
			name: "valid image with no opts succeeds",
			img:  validImg,
			path: func() string {
				dir := t.TempDir()
				return dir + "/encode-test-default.heic"
			},
			wantErr: false,
		},
		{
			name: "valid image with WithQuality(80) succeeds",
			img:  validImg,
			path: func() string {
				dir := t.TempDir()
				return dir + "/encode-test-quality80.heic"
			},
			opts:    []EncoderOption{WithQuality(80)},
			wantErr: false,
		},
		{
			name: "valid image with WithQuality(50) and WithLossless(false) succeeds",
			img:  validImg,
			path: func() string {
				dir := t.TempDir()
				return dir + "/encode-test-quality50-lossy.heic"
			},
			opts:    []EncoderOption{WithQuality(50), WithLossless(false)},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outPath := tc.path()
			err := EncodeImageAsHeif(tc.img, outPath, tc.opts...)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				for _, substr := range tc.errContains {
					if !containsStr(err.Error(), substr) {
						t.Errorf("error %q does not contain %q", err.Error(), substr)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			info, statErr := os.Stat(outPath)
			if statErr != nil {
				t.Fatalf("output file does not exist: %v", statErr)
			}
			if info.Size() == 0 {
				t.Errorf("output file is empty")
			}
		})
	}
}

func TestConvertToHeif(t *testing.T) {
	tests := []struct {
		name        string
		srcPath     string
		dstPath     func() string
		opts        []EncoderOption
		wantErr     bool
		errContains []string
	}{
		{
			name:        "empty source path returns error",
			srcPath:     "",
			dstPath:     func() string { return "out.heic" },
			wantErr:     true,
			errContains: []string{"source", "empty"},
		},
		{
			name:    "empty destination path returns error",
			srcPath: "input.jpg",
			dstPath: func() string { return "" },
			wantErr: true,
			errContains: []string{"destination", "empty"},
		},
		{
			name:    "nonexistent source file returns error",
			srcPath: "nonexistent.jpg",
			dstPath: func() string {
				dir := t.TempDir()
				return dir + "/out.heic"
			},
			wantErr: true,
		},
		{
			name:    "valid JPEG converts to HEIF successfully",
			srcPath: "images/libheif-generated.jpeg",
			dstPath: func() string {
				dir := t.TempDir()
				return dir + "/convert-test.heic"
			},
			wantErr: false,
		},
		{
			name:    "valid JPEG with WithQuality(70) succeeds",
			srcPath: "images/libheif-generated.jpeg",
			dstPath: func() string {
				dir := t.TempDir()
				return dir + "/convert-test-quality70.heic"
			},
			opts:    []EncoderOption{WithQuality(70)},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dstPath := tc.dstPath()
			err := ConvertToHeif(tc.srcPath, dstPath, tc.opts...)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				for _, substr := range tc.errContains {
					if !containsStr(err.Error(), substr) {
						t.Errorf("error %q does not contain %q", err.Error(), substr)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			info, statErr := os.Stat(dstPath)
			if statErr != nil {
				t.Fatalf("output file does not exist: %v", statErr)
			}
			if info.Size() == 0 {
				t.Errorf("output file is empty")
			}
		})
	}
}

// containsStr is a helper to check if s contains substr (case-sensitive).
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

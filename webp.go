package libheif

import (
	"fmt"
	"image"
	"os"

	_ "golang.org/x/image/webp"
)

// WebpToHeif converts a WebP image file to HEIF format.
//
// The source file must be a valid WebP image. Non-WebP files are rejected
// with a descriptive error naming the actual format detected. Options are
// passed through to EncodeImageAsHeif.
//
// Importing this package registers the WebP decoder with Go's image package,
// enabling ConvertToHeif to also handle WebP files.
//
// Example:
//
//	err := WebpToHeif("photo.webp", "photo.heic", WithQuality(80))
//	if err != nil {
//		log.Fatal(err)
//	}
func WebpToHeif(webpPath string, heifPath string, opts ...EncoderOption) error {
	if webpPath == "" {
		return fmt.Errorf("webp source path is empty")
	}
	if heifPath == "" {
		return fmt.Errorf("heif output path is empty")
	}

	file, err := os.Open(webpPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode source image: %w", err)
	}

	if format != "webp" {
		return fmt.Errorf("expected webp format, got %s", format)
	}

	if err := EncodeImageAsHeif(img, heifPath, opts...); err != nil {
		return fmt.Errorf("failed to encode to HEIF: %w", err)
	}

	return nil
}

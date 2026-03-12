package libheif

import (
	"fmt"
	"image"
	"os"

	"github.com/strukturag/libheif/go/heif"
)

// mapCompression converts a Compression constant to its heif.Compression equivalent.
// Defaults to heif.CompressionHEVC for unknown values.
func mapCompression(c Compression) heif.Compression {
	switch c {
	case CompressionHEVC:
		return heif.CompressionHEVC
	case CompressionAV1:
		return heif.CompressionAV1
	case CompressionAVC:
		return heif.CompressionAVC
	case CompressionJPEG:
		return heif.CompressionJPEG
	default:
		return heif.CompressionHEVC
	}
}

// mapLossless converts a bool to the heif.LosslessMode equivalent.
func mapLossless(l bool) heif.LosslessMode {
	if l {
		return heif.LosslessModeEnabled
	}
	return heif.LosslessModeDisabled
}

// mapLogLevel converts a LogLevel constant to its heif.LoggingLevel equivalent.
// Defaults to heif.LoggingLevelFull for unknown values.
func mapLogLevel(l LogLevel) heif.LoggingLevel {
	switch l {
	case LogLevelNone:
		return heif.LoggingLevelNone
	case LogLevelBasic:
		return heif.LoggingLevelBasic
	case LogLevelFull:
		return heif.LoggingLevelFull
	default:
		return heif.LoggingLevelFull
	}
}

// EncodeImageAsHeif encodes an image.Image to HEIF format and saves it to the given path.
//
// Options are applied in order using the Functional Options pattern. Defaults
// (quality=100, HEVC, lossless=true, LogLevelFull) are used when no options are provided.
//
// Example:
//
//	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
//	err := EncodeImageAsHeif(img, "output.heic", WithQuality(80), WithLossless(false))
//	if err != nil {
//		log.Fatal(err)
//	}
func EncodeImageAsHeif(img image.Image, path string, opts ...EncoderOption) error {
	if img == nil {
		return fmt.Errorf("image is nil")
	}
	if path == "" {
		return fmt.Errorf("output path is empty")
	}

	cfg := defaultConfig()
	applyOptions(cfg, opts...)

	ctx, err := heif.EncodeFromImage(
		img,
		mapCompression(cfg.compression),
		cfg.quality,
		mapLossless(cfg.lossless),
		mapLogLevel(cfg.logging),
	)
	if err != nil {
		return fmt.Errorf("failed to HEIF encode image: %w", err)
	}

	if err := ctx.WriteToFile(path); err != nil {
		return fmt.Errorf("failed to write HEIF file to %s: %w", path, err)
	}

	return nil
}

// ConvertToHeif converts an image file at srcPath to HEIF format and saves it to dstPath.
//
// The source format is auto-detected via image.Decode (supports JPEG, PNG, and any other
// format registered with the image package). Options are passed through to EncodeImageAsHeif.
//
// Example:
//
//	err := ConvertToHeif("input.jpeg", "output.heic", WithQuality(70))
//	if err != nil {
//		log.Fatal(err)
//	}
func ConvertToHeif(srcPath string, dstPath string, opts ...EncoderOption) error {
	if srcPath == "" {
		return fmt.Errorf("source path is empty")
	}
	if dstPath == "" {
		return fmt.Errorf("destination path is empty")
	}

	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode source image %s: %w", srcPath, err)
	}

	if err := EncodeImageAsHeif(img, dstPath, opts...); err != nil {
		return fmt.Errorf("failed to encode to HEIF: %w", err)
	}

	return nil
}

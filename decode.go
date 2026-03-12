package libheif

import "fmt"

// DecodeHeifToJpeg decodes a HEIF file and saves it as JPEG with the specified quality.
//
// Quality must be between 1 and 100 inclusive; invalid values return an error.
// Both heifPath and jpegPath must be non-empty.
//
// Example:
//
//	err := DecodeHeifToJpeg("input.heic", "output.jpg", 80)
//	if err != nil {
//		log.Fatal(err)
//	}
func DecodeHeifToJpeg(heifPath string, jpegPath string, quality int) error {
	if heifPath == "" {
		return fmt.Errorf("heif source path is empty")
	}
	if jpegPath == "" {
		return fmt.Errorf("jpeg output path is empty")
	}
	if quality < 1 || quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100, got: %d", quality)
	}

	img, err := ReturnImageFromHeif(heifPath)
	if err != nil {
		return fmt.Errorf("failed to decode HEIF file: %w", err)
	}

	if err := saveAsJpeg(img, jpegPath, quality); err != nil {
		return fmt.Errorf("failed to save as JPEG: %w", err)
	}

	return nil
}

// DecodeHeifToPng decodes a HEIF file and saves it as PNG.
//
// Both heifPath and pngPath must be non-empty.
//
// Example:
//
//	err := DecodeHeifToPng("input.heic", "output.png")
//	if err != nil {
//		log.Fatal(err)
//	}
func DecodeHeifToPng(heifPath string, pngPath string) error {
	if heifPath == "" {
		return fmt.Errorf("heif source path is empty")
	}
	if pngPath == "" {
		return fmt.Errorf("png output path is empty")
	}

	img, err := ReturnImageFromHeif(heifPath)
	if err != nil {
		return fmt.Errorf("failed to decode HEIF file: %w", err)
	}

	if err := saveAsPng(img, pngPath); err != nil {
		return fmt.Errorf("failed to save as PNG: %w", err)
	}

	return nil
}

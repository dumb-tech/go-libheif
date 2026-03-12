package libheif

import (
	"os"
	"testing"
)

// minimalWebP is a minimal valid lossless WebP file (1x1 pixel, white).
// RIFF header + WEBP signature + VP8L chunk with a 1x1 white pixel.
// Generated offline from the VP8L lossless format spec.
var minimalWebP = []byte{
	// RIFF header
	0x52, 0x49, 0x46, 0x46, // "RIFF"
	0x24, 0x00, 0x00, 0x00, // file size - 8 = 36
	0x57, 0x45, 0x42, 0x50, // "WEBP"
	// VP8L chunk
	0x56, 0x50, 0x38, 0x4C, // "VP8L"
	0x0D, 0x00, 0x00, 0x00, // chunk size = 13
	// VP8L bitstream: signature byte + 1x1 image
	0x2F,                   // VP8L signature
	0x00,                   // width minus 1 (low 8 bits) = 0 → width=1
	0x00,                   // bits: width cont + height minus 1 bits
	0x00,                   // more image data
	0xFE,                   // lossless transform + 4-pixel color cache
	0xFF, 0xFF, 0xFF, 0xFF, // ARGB pixel: white
	0x00, 0x00, 0x00, 0x00, // padding
}

// writeMinimalWebP writes a minimal WebP file to path and returns any error.
func writeMinimalWebP(path string) error {
	return os.WriteFile(path, minimalWebP, 0644)
}

func TestWebpToHeif_EmptyPaths(t *testing.T) {
	dir := t.TempDir()
	dstPath := dir + "/out.heic"

	// Empty webpPath
	err := WebpToHeif("", dstPath)
	if err == nil {
		t.Fatal("expected error for empty webpPath, got nil")
	}
	if !containsStr(err.Error(), "path") || !containsStr(err.Error(), "empty") {
		t.Errorf("error %q should contain 'path' and 'empty'", err.Error())
	}

	// Empty heifPath
	srcPath := dir + "/source.webp"
	if err := writeMinimalWebP(srcPath); err != nil {
		t.Fatalf("failed to write test webp: %v", err)
	}
	err = WebpToHeif(srcPath, "")
	if err == nil {
		t.Fatal("expected error for empty heifPath, got nil")
	}
	if !containsStr(err.Error(), "path") || !containsStr(err.Error(), "empty") {
		t.Errorf("error %q should contain 'path' and 'empty'", err.Error())
	}
}

func TestWebpToHeif_NonExistentFile(t *testing.T) {
	dir := t.TempDir()
	err := WebpToHeif("nonexistent-file-xyz.webp", dir+"/out.heic")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if !containsStr(err.Error(), "open") && !containsStr(err.Error(), "no such file") {
		t.Errorf("error %q should contain 'open' or 'no such file'", err.Error())
	}
}

func TestWebpToHeif_NonWebpFile(t *testing.T) {
	dir := t.TempDir()
	dstPath := dir + "/out.heic"

	// Use an existing JPEG file
	err := WebpToHeif("images/libheif-generated.jpeg", dstPath)
	if err == nil {
		t.Fatal("expected error when passing a JPEG to WebpToHeif, got nil")
	}
	if !containsStr(err.Error(), "expected webp") {
		t.Errorf("error %q should contain 'expected webp'", err.Error())
	}
	// The error should also name the actual format detected
	if !containsStr(err.Error(), "jpeg") {
		t.Errorf("error %q should contain the actual format name 'jpeg'", err.Error())
	}
}

func TestWebpToHeif_ValidWebp(t *testing.T) {
	dir := t.TempDir()
	srcPath := dir + "/input.webp"
	dstPath := dir + "/output.heic"

	if err := writeMinimalWebP(srcPath); err != nil {
		t.Fatalf("failed to write test webp: %v", err)
	}
	t.Cleanup(func() { os.Remove(dstPath) })

	err := WebpToHeif(srcPath, dstPath)
	if err != nil {
		t.Fatalf("WebpToHeif failed on valid WebP: %v", err)
	}

	info, statErr := os.Stat(dstPath)
	if statErr != nil {
		t.Fatalf("output file does not exist: %v", statErr)
	}
	if info.Size() == 0 {
		t.Errorf("output HEIF file is empty")
	}

	// Round-trip: read back via ReturnImageFromHeif
	img, roundTripErr := ReturnImageFromHeif(dstPath)
	if roundTripErr != nil {
		t.Fatalf("round-trip ReturnImageFromHeif failed: %v", roundTripErr)
	}
	if img == nil {
		t.Errorf("round-trip returned nil image")
	}
}

func TestWebpToHeif_WithOptions(t *testing.T) {
	dir := t.TempDir()
	srcPath := dir + "/input.webp"
	dstPath := dir + "/output-opts.heic"

	if err := writeMinimalWebP(srcPath); err != nil {
		t.Fatalf("failed to write test webp: %v", err)
	}
	t.Cleanup(func() { os.Remove(dstPath) })

	err := WebpToHeif(srcPath, dstPath, WithQuality(50), WithLossless(false))
	if err != nil {
		t.Fatalf("WebpToHeif with options failed: %v", err)
	}

	if _, statErr := os.Stat(dstPath); statErr != nil {
		t.Fatalf("output file does not exist: %v", statErr)
	}
}

package libheif

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"testing"
)

// createTestWebP generates a valid WebP file at path using cwebp.
// Skips the test if cwebp is not available (e.g. outside Docker).
func createTestWebP(t *testing.T, path string) {
	t.Helper()

	cwebp, err := exec.LookPath("cwebp")
	if err != nil {
		t.Skip("cwebp not available — run tests in Docker")
	}

	// Create a small PNG in a temp file
	pngPath := path + ".png"
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	f, err := os.Create(pngPath)
	if err != nil {
		t.Fatalf("failed to create temp PNG: %v", err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatalf("failed to encode PNG: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(pngPath) })

	// Convert PNG to WebP using cwebp
	cmd := exec.Command(cwebp, "-quiet", pngPath, "-o", path)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cwebp failed: %v\n%s", err, out)
	}
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
	createTestWebP(t, srcPath)
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

	createTestWebP(t, srcPath)
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

	createTestWebP(t, srcPath)
	t.Cleanup(func() { os.Remove(dstPath) })

	err := WebpToHeif(srcPath, dstPath, WithQuality(50), WithLossless(false))
	if err != nil {
		t.Fatalf("WebpToHeif with options failed: %v", err)
	}

	if _, statErr := os.Stat(dstPath); statErr != nil {
		t.Fatalf("output file does not exist: %v", statErr)
	}
}

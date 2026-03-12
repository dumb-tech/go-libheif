package libheif

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	t.Run("quality is 100", func(t *testing.T) {
		if cfg.quality != 100 {
			t.Errorf("expected quality=100, got %d", cfg.quality)
		}
	})

	t.Run("compression is CompressionHEVC", func(t *testing.T) {
		if cfg.compression != CompressionHEVC {
			t.Errorf("expected compression=CompressionHEVC, got %v", cfg.compression)
		}
	})

	t.Run("lossless is true", func(t *testing.T) {
		if !cfg.lossless {
			t.Errorf("expected lossless=true, got false")
		}
	})

	t.Run("logging is LogLevelFull", func(t *testing.T) {
		if cfg.logging != LogLevelFull {
			t.Errorf("expected logging=LogLevelFull, got %v", cfg.logging)
		}
	})
}

// TestEncoderOptionType verifies EncoderOption is func(*encoderConfig) at compile time.
func TestEncoderOptionType(t *testing.T) {
	var _ EncoderOption = func(c *encoderConfig) {
		c.quality = 50
	}
}

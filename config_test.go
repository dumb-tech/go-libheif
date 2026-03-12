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

func TestWithQuality(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    int
	}{
		{"normal value 80", 80, 80},
		{"clamped low: 0 -> 1", 0, 1},
		{"clamped low: -5 -> 1", -5, 1},
		{"clamped high: 101 -> 100", 101, 100},
		{"clamped high: 999 -> 100", 999, 100},
		{"boundary low: 1 -> 1", 1, 1},
		{"boundary high: 100 -> 100", 100, 100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfig()
			opt := WithQuality(tc.input)
			opt(cfg)
			if cfg.quality != tc.want {
				t.Errorf("WithQuality(%d): expected quality=%d, got %d", tc.input, tc.want, cfg.quality)
			}
		})
	}
}

func TestWithCompression(t *testing.T) {
	tests := []struct {
		name string
		comp Compression
	}{
		{"CompressionAV1", CompressionAV1},
		{"CompressionHEVC (no-op, already default)", CompressionHEVC},
		{"CompressionAVC", CompressionAVC},
		{"CompressionJPEG", CompressionJPEG},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfig()
			opt := WithCompression(tc.comp)
			opt(cfg)
			if cfg.compression != tc.comp {
				t.Errorf("WithCompression(%v): expected compression=%v, got %v", tc.comp, tc.comp, cfg.compression)
			}
		})
	}
}

func TestWithLossless(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{"false disables lossless", false},
		{"true keeps lossless (no-op, already default)", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfig()
			opt := WithLossless(tc.value)
			opt(cfg)
			if cfg.lossless != tc.value {
				t.Errorf("WithLossless(%v): expected lossless=%v, got %v", tc.value, tc.value, cfg.lossless)
			}
		})
	}
}

func TestWithLogging(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
	}{
		{"LogLevelNone", LogLevelNone},
		{"LogLevelBasic", LogLevelBasic},
		{"LogLevelFull (no-op, already default)", LogLevelFull},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfig()
			opt := WithLogging(tc.level)
			opt(cfg)
			if cfg.logging != tc.level {
				t.Errorf("WithLogging(%v): expected logging=%v, got %v", tc.level, tc.level, cfg.logging)
			}
		})
	}
}

func TestOptionIndependence(t *testing.T) {
	t.Run("WithQuality and WithLossless do not interfere", func(t *testing.T) {
		cfg := defaultConfig()
		WithQuality(50)(cfg)
		WithLossless(false)(cfg)
		if cfg.quality != 50 {
			t.Errorf("expected quality=50 after applying both options, got %d", cfg.quality)
		}
		if cfg.lossless != false {
			t.Errorf("expected lossless=false after applying both options, got %v", cfg.lossless)
		}
	})

	t.Run("options do not reset previously set fields", func(t *testing.T) {
		cfg := defaultConfig()
		WithQuality(75)(cfg)
		WithCompression(CompressionAV1)(cfg)
		if cfg.quality != 75 {
			t.Errorf("expected quality=75 after also setting compression, got %d", cfg.quality)
		}
		if cfg.compression != CompressionAV1 {
			t.Errorf("expected compression=CompressionAV1 after also setting quality, got %v", cfg.compression)
		}
	})
}

func TestApplyOptions(t *testing.T) {
	t.Run("applies all options to config", func(t *testing.T) {
		cfg := defaultConfig()
		applyOptions(cfg, WithQuality(60), WithCompression(CompressionAVC), WithLossless(false), WithLogging(LogLevelBasic))
		if cfg.quality != 60 {
			t.Errorf("expected quality=60, got %d", cfg.quality)
		}
		if cfg.compression != CompressionAVC {
			t.Errorf("expected compression=CompressionAVC, got %v", cfg.compression)
		}
		if cfg.lossless != false {
			t.Errorf("expected lossless=false, got %v", cfg.lossless)
		}
		if cfg.logging != LogLevelBasic {
			t.Errorf("expected logging=LogLevelBasic, got %v", cfg.logging)
		}
	})

	t.Run("applyOptions with no options leaves config unchanged", func(t *testing.T) {
		cfg := defaultConfig()
		applyOptions(cfg)
		if cfg.quality != 100 {
			t.Errorf("expected quality=100, got %d", cfg.quality)
		}
		if cfg.compression != CompressionHEVC {
			t.Errorf("expected compression=CompressionHEVC, got %v", cfg.compression)
		}
		if !cfg.lossless {
			t.Errorf("expected lossless=true, got false")
		}
		if cfg.logging != LogLevelFull {
			t.Errorf("expected logging=LogLevelFull, got %v", cfg.logging)
		}
	})

	t.Run("all four option functions satisfy EncoderOption type", func(t *testing.T) {
		var opts []EncoderOption
		opts = append(opts, WithQuality(80))
		opts = append(opts, WithCompression(CompressionAV1))
		opts = append(opts, WithLossless(false))
		opts = append(opts, WithLogging(LogLevelNone))
		if len(opts) != 4 {
			t.Errorf("expected 4 options, got %d", len(opts))
		}
	})
}

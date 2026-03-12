package libheif

// Compression wraps heif.Compression to avoid direct dependency on CGo types in public API.
type Compression int

const (
	CompressionHEVC Compression = iota
	CompressionAV1
	CompressionAVC
	CompressionJPEG
)

// LogLevel controls verbosity of heif encoder logging.
type LogLevel int

const (
	LogLevelNone  LogLevel = iota
	LogLevelBasic
	LogLevelFull
)

type encoderConfig struct {
	quality     int
	compression Compression
	lossless    bool
	logging     LogLevel
}

func defaultConfig() *encoderConfig {
	return &encoderConfig{
		quality:     100,
		compression: CompressionHEVC,
		lossless:    true,
		logging:     LogLevelFull,
	}
}

// EncoderOption configures an encoderConfig via the Functional Options pattern.
// Pass any number of EncoderOption values to encoding functions.
type EncoderOption func(*encoderConfig)

// WithQuality returns an EncoderOption that sets the encoder quality.
// Values are clamped silently: below 1 is set to 1, above 100 is set to 100.
func WithQuality(quality int) EncoderOption {
	return func(c *encoderConfig) {
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
		c.quality = quality
	}
}

// WithCompression returns an EncoderOption that sets the encoder compression type.
func WithCompression(compression Compression) EncoderOption {
	return func(c *encoderConfig) {
		c.compression = compression
	}
}

// WithLossless returns an EncoderOption that enables or disables lossless encoding.
func WithLossless(lossless bool) EncoderOption {
	return func(c *encoderConfig) {
		c.lossless = lossless
	}
}

// WithLogging returns an EncoderOption that sets the encoder log verbosity level.
func WithLogging(logging LogLevel) EncoderOption {
	return func(c *encoderConfig) {
		c.logging = logging
	}
}

// applyOptions applies a variadic list of EncoderOption functions to a config.
// This helper is consumed by Phase 3 encoding functions to build their final config.
func applyOptions(cfg *encoderConfig, opts ...EncoderOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

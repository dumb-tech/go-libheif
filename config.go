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

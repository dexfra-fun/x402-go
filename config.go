package x402

import (
	"fmt"
	"time"
)

const (
	// DefaultVerifyTimeoutSeconds is the default verify timeout in seconds.
	DefaultVerifyTimeoutSeconds = 5
	// DefaultSettleTimeoutSeconds is the default settle timeout in seconds.
	DefaultSettleTimeoutSeconds = 60
	// DefaultRequestTimeoutSeconds is the default request timeout in seconds.
	DefaultRequestTimeoutSeconds = 120
)

// TimeoutConfig holds timeout configuration for payment operations.
type TimeoutConfig struct {
	// VerifyTimeout is the maximum time to wait for payment verification.
	VerifyTimeout time.Duration

	// SettleTimeout is the maximum time to wait for payment settlement.
	SettleTimeout time.Duration

	// RequestTimeout is the overall timeout for HTTP requests (optional).
	RequestTimeout time.Duration
}

// NewDefaultTimeouts returns a new TimeoutConfig with sensible defaults for payment operations.
func NewDefaultTimeouts() TimeoutConfig {
	return TimeoutConfig{
		VerifyTimeout:  DefaultVerifyTimeoutSeconds * time.Second,
		SettleTimeout:  DefaultSettleTimeoutSeconds * time.Second,
		RequestTimeout: DefaultRequestTimeoutSeconds * time.Second,
	}
}

// WithVerifyTimeout returns a new TimeoutConfig with updated verify timeout.
func (tc TimeoutConfig) WithVerifyTimeout(d time.Duration) TimeoutConfig {
	tc.VerifyTimeout = d
	return tc
}

// WithSettleTimeout returns a new TimeoutConfig with updated settle timeout.
func (tc TimeoutConfig) WithSettleTimeout(d time.Duration) TimeoutConfig {
	tc.SettleTimeout = d
	return tc
}

// WithRequestTimeout returns a new TimeoutConfig with updated request timeout.
func (tc TimeoutConfig) WithRequestTimeout(d time.Duration) TimeoutConfig {
	tc.RequestTimeout = d
	return tc
}

// Validate ensures timeout values are reasonable.
func (tc TimeoutConfig) Validate() error {
	if tc.VerifyTimeout <= 0 {
		return fmt.Errorf("verify timeout must be positive, got %v", tc.VerifyTimeout)
	}
	if tc.SettleTimeout <= 0 {
		return fmt.Errorf("settle timeout must be positive, got %v", tc.SettleTimeout)
	}
	if tc.SettleTimeout < tc.VerifyTimeout {
		return fmt.Errorf("settle timeout (%v) should be >= verify timeout (%v)",
			tc.SettleTimeout, tc.VerifyTimeout)
	}
	return nil
}

package x402

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

const (
	// defaultHTTPTimeout is the default timeout for HTTP requests to facilitator.
	defaultHTTPTimeout = 10 * time.Second
)

// FacilitatorClient handles communication with x402 facilitator services.
type FacilitatorClient struct {
	baseURL    string
	httpClient *http.Client
	cache      *FeePayerCache
	logger     Logger
}

// SupportedResponse represents the /supported endpoint response.
type SupportedResponse struct {
	Kinds []Kind `json:"kinds"`
}

// Kind represents a supported payment kind.
type Kind struct {
	X402Version int     `json:"x402Version"`
	Scheme      string  `json:"scheme"`
	Network     string  `json:"network"`
	Extra       *Extras `json:"extra,omitempty"`
}

// Extras contains additional metadata.
type Extras struct {
	FeePayer string `json:"feePayer"`
}

// NewFacilitatorClient creates a new facilitator client.
func NewFacilitatorClient(baseURL string, cache *FeePayerCache, logger Logger) *FacilitatorClient {
	if logger == nil {
		logger = &DefaultLogger{}
	}

	return &FacilitatorClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		cache:  cache,
		logger: logger,
	}
}

// GetFeePayer retrieves the fee payer for a given network
// Uses cache if available, otherwise fetches from facilitator.
func (c *FacilitatorClient) GetFeePayer(ctx context.Context, network string) (string, error) {
	// Try cache first
	if feePayer, found := c.cache.Get(network); found {
		c.logger.Printf("[x402] Fee payer cache hit: network=%s feePayer=%s", network, feePayer)
		return feePayer, nil
	}

	c.logger.Printf("[x402] Fee payer cache miss: network=%s, fetching from facilitator", network)

	// Cache miss - fetch from facilitator
	feePayer, found, err := c.fetchFeePayer(ctx, network)
	if err != nil {
		return "", fmt.Errorf("fetch fee payer: %w", err)
	}

	if !found {
		return "", ErrFeePayerNotFound
	}

	// Cache the result
	c.cache.Set(network, feePayer)
	c.logger.Printf("[x402] Fee payer cached: network=%s feePayer=%s", network, feePayer)

	return feePayer, nil
}

// findFeePayerInKinds searches for a fee payer in the list of supported kinds.
func findFeePayerInKinds(kinds []Kind, network string) (string, bool) {
	target := strings.ToLower(strings.TrimSpace(network))
	for _, k := range kinds {
		if strings.ToLower(k.Network) == target {
			if k.Extra != nil && strings.TrimSpace(k.Extra.FeePayer) != "" {
				return k.Extra.FeePayer, true
			}
			// Network found but no fee payer
			return "", false
		}
	}
	// Network not found in supported list
	return "", false
}

// fetchFeePayer fetches fee payer from facilitator's /supported endpoint.
func (c *FacilitatorClient) fetchFeePayer(
	ctx context.Context,
	network string,
) (string, bool, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", false, fmt.Errorf("parse facilitator URL: %w", err)
	}

	u.Path = path.Join(u.Path, "supported")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", false, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Errorf("[x402] Error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var data SupportedResponse
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", false, fmt.Errorf("decode json: %w", err)
	}

	// Find matching network and fee payer
	feePayer, found := findFeePayerInKinds(data.Kinds, network)
	return feePayer, found, nil
}

// GetSupported fetches all supported payment kinds from the facilitator.
func (c *FacilitatorClient) GetSupported(ctx context.Context) ([]Kind, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse facilitator URL: %w", err)
	}

	u.Path = path.Join(u.Path, "supported")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Errorf("[x402] Error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var data SupportedResponse
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	return data.Kinds, nil
}

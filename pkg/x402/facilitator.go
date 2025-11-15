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
	x402 "github.com/dexfra-fun/x402-go"
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

// extractPayerFromPayment attempts to extract the payer address from payment payload.
// This is a fallback for when facilitator doesn't return the payer field.
func extractPayerFromPayment(payment x402.PaymentPayload) string {
	// For EVM payments, extract from authorization.from
	if payload, ok := payment.Payload.(map[string]any); ok {
		if auth, ok := payload["authorization"].(map[string]any); ok {
			if from, ok := auth["from"].(string); ok {
				return from
			}
		}
	}
	// For other payment types or if extraction fails, return empty string
	return ""
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

// VerifyResult contains the results of payment verification.
type VerifyResult struct {
	IsValid       bool
	InvalidReason string
	Payer         string
}

// buildVerifyRequest creates the HTTP request for payment verification.
func (c *FacilitatorClient) buildVerifyRequest(
	ctx context.Context,
	payment x402.PaymentPayload,
	requirement x402.PaymentRequirement,
) (*http.Request, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse facilitator URL: %w", err)
	}
	u.Path = path.Join(u.Path, "verify")

	reqBody := map[string]any{
		"x402Version":         1,
		"paymentPayload":      payment,
		"paymentRequirements": requirement,
	}

	jsonBytes, err := sonic.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	c.logger.Printf("[x402] Verify request: %s", string(jsonBytes))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(string(jsonBytes)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// parseVerifyResponse decodes and processes the verification response.
func (c *FacilitatorClient) parseVerifyResponse(resp *http.Response, payment x402.PaymentPayload) (bool, string, string, error) {
	if resp.StatusCode != http.StatusOK {
		c.logger.Printf("[x402] Verify failed: status=%s", resp.Status)
		return false, "", "", fmt.Errorf("unexpected status %s", resp.Status)
	}

	var result struct {
		IsValid       bool   `json:"isValid"`
		InvalidReason string `json:"invalidReason,omitempty"`
		Payer         string `json:"payer"`
	}
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", "", fmt.Errorf("decode json: %w", err)
	}

	c.logger.Printf("[x402] Verify response: isValid=%v payer=%s invalidReason=%s",
		result.IsValid, result.Payer, result.InvalidReason)

	if result.Payer == "" && result.IsValid {
		result.Payer = extractPayerFromPayment(payment)
		c.logger.Printf("[x402] Payer fallback: %s", result.Payer)
	}

	return result.IsValid, result.InvalidReason, result.Payer, nil
}

// Verify verifies a payment with the facilitator.
func (c *FacilitatorClient) Verify(
	ctx context.Context,
	payment x402.PaymentPayload,
	requirement x402.PaymentRequirement,
) (bool, string, string, error) {
	req, err := c.buildVerifyRequest(ctx, payment, requirement)
	if err != nil {
		return false, "", "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, "", "", fmt.Errorf("http request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Errorf("[x402] Error closing response body: %v", closeErr)
		}
	}()

	return c.parseVerifyResponse(resp, payment)
}

// Settle settles a payment with the facilitator.
func (c *FacilitatorClient) Settle(
	ctx context.Context,
	payment x402.PaymentPayload,
	requirement x402.PaymentRequirement,
) (*x402.SettlementResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse facilitator URL: %w", err)
	}

	u.Path = path.Join(u.Path, "settle")

	// Create request body matching facilitator API spec
	reqBody := map[string]any{
		"x402Version":         1,
		"paymentPayload":      payment,
		"paymentRequirements": requirement,
	}

	jsonBytes, err := sonic.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// LOG: Debug request being sent
	c.logger.Printf("[x402] Settle request: %s", string(jsonBytes))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.String(),
		strings.NewReader(string(jsonBytes)),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
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
		// LOG: Debug error response
		c.logger.Printf("[x402] Settle failed: status=%s", resp.Status)
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var settlement x402.SettlementResponse
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&settlement); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	// LOG: Debug response received
	c.logger.Printf("[x402] Settle response: txHash=%s", settlement.Transaction)

	return &settlement, nil
}

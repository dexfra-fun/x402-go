package gin

import (
	"net/http"

	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gin-gonic/gin"
	mark3labs "github.com/mark3labs/x402-go"
	ginx402 "github.com/mark3labs/x402-go/http/gin"
	x402http "github.com/mark3labs/x402-go/http"
)

// NewMiddleware creates a new Gin middleware for x402 payment handling
func NewMiddleware(config *x402.Config) gin.HandlerFunc {
	// Create x402 middleware
	middleware, err := x402.New(config)
	if err != nil {
		config.Logger.Errorf("[x402-gin] Failed to create middleware: %v", err)
		// Return a middleware that always returns error
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Payment middleware configuration error",
			})
			c.Abort()
		}
	}

	return func(c *gin.Context) {
		// Create resource from request
		resource := x402.Resource{
			Path:   c.Request.URL.Path,
			Method: c.Request.Method,
			Params: make(map[string]string),
		}

		// Extract query parameters
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				resource.Params[key] = values[0]
			}
		}

		// Process payment requirement
		requirement, paymentInfo, err := middleware.ProcessRequest(c.Request.Context(), resource)
		if err != nil {
			config.Logger.Errorf("[x402-gin] Failed to process payment: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Payment processing error",
			})
			c.Abort()
			return
		}

		// Free endpoint - no payment required
		if requirement == nil {
			c.Next()
			return
		}

		// Payment required - configure x402 HTTP middleware
		x402Config := &x402http.Config{
			FacilitatorURL: config.FacilitatorURL,
			PaymentRequirements: []mark3labs.PaymentRequirement{
				*requirement,
			},
		}

		// Apply mark3labs x402 Gin middleware
		x402Handler := ginx402.NewGinX402Middleware(x402Config)
		x402Handler(c)

		// Store payment info in context for later use
		if paymentInfo != nil {
			c.Set("x402_payment_info", paymentInfo)
		}
	}
}

// GetPaymentInfo retrieves payment information from the Gin context
func GetPaymentInfo(c *gin.Context) (*x402.PaymentInfo, bool) {
	if info, exists := c.Get("x402_payment_info"); exists {
		if paymentInfo, ok := info.(*x402.PaymentInfo); ok {
			return paymentInfo, true
		}
	}
	return nil, false
}

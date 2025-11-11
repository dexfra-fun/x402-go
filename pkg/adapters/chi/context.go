package chi

import (
	"context"

	"github.com/dexfra-fun/x402-go/pkg/x402"
)

type contextKey string

const paymentInfoKey contextKey = "x402_payment_info"

// setPaymentInfo stores payment information in the context
func setPaymentInfo(ctx context.Context, info *x402.PaymentInfo) context.Context {
	return context.WithValue(ctx, paymentInfoKey, info)
}

// GetPaymentInfo retrieves payment information from the context
func GetPaymentInfo(ctx context.Context) (*x402.PaymentInfo, bool) {
	info, ok := ctx.Value(paymentInfoKey).(*x402.PaymentInfo)
	return info, ok
}

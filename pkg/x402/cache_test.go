package x402

import (
	"testing"
	"time"
)

func TestFeePayerCache(t *testing.T) {
	t.Run("get empty cache", func(t *testing.T) {
		cache := NewFeePayerCache(5 * time.Minute)
		_, found := cache.Get("solana-devnet")
		if found {
			t.Error("expected not found in empty cache")
		}
	})

	t.Run("set and get", func(t *testing.T) {
		cache := NewFeePayerCache(5 * time.Minute)
		cache.Set("solana-devnet", "feePayer123")
		
		feePayer, found := cache.Get("solana-devnet")
		if !found {
			t.Error("expected to find cached value")
		}
		if feePayer != "feePayer123" {
			t.Errorf("expected feePayer123, got %s", feePayer)
		}
	})

	t.Run("expiration", func(t *testing.T) {
		cache := NewFeePayerCache(100 * time.Millisecond)
		cache.Set("solana-devnet", "feePayer123")
		
		// Should be found immediately
		_, found := cache.Get("solana-devnet")
		if !found {
			t.Error("expected to find cached value immediately")
		}
		
		// Wait for expiration
		time.Sleep(150 * time.Millisecond)
		
		_, found = cache.Get("solana-devnet")
		if found {
			t.Error("expected cache to expire")
		}
	})

	t.Run("clear", func(t *testing.T) {
		cache := NewFeePayerCache(5 * time.Minute)
		cache.Set("solana-devnet", "feePayer1")
		cache.Set("solana-mainnet", "feePayer2")
		
		cache.Clear()
		
		_, found1 := cache.Get("solana-devnet")
		_, found2 := cache.Get("solana-mainnet")
		
		if found1 || found2 {
			t.Error("expected cache to be cleared")
		}
	})
}

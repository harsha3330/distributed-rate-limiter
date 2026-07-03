package ratelimiter

import (
	"sync"
	"time"

	"github.com/harsha3330/distributed-rate-limiter/internal/bucket"
	"github.com/harsha3330/distributed-rate-limiter/internal/policy"
)

type RateLimiter struct {
	mu sync.Mutex
	bs *bucket.BucketStore
	ps *policy.PolicyStore
}

func NewRateLimiter(bs *bucket.BucketStore, ps *policy.PolicyStore) *RateLimiter {
	return &RateLimiter{
		bs: bs,
		ps: ps,
	}
}

type RateLimiterResponse struct {
	Available                 bool      `json:"available"`
	WindowEndTime             time.Time `json:"window_endtime"`
	ExpiryInSeconds           uint64    `json:"expiry_in_seconds"`
	RequestCount              uint64    `json:"request_count"`
	RequestsAvailableInWindow uint64    `json:"requests_available_in_window"`
}

func (rl *RateLimiter) Check(key bucket.BucketKey) (*RateLimiterResponse, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	policyId := key.PolicyID
	currentTime := time.Now()
	policy, err := rl.ps.GetPolicy(policyId)
	if err != nil {
		return nil, err
	}
	bkt, created, err := rl.bs.GetOrCreateBucket(key)
	if err != nil {
		return nil, err
	}
	if created {
		bkt.Requests = 1
		return &RateLimiterResponse{
			Available:                 true,
			WindowEndTime:             bkt.WindowStart.Add(policy.Window),
			ExpiryInSeconds:           uint64(policy.Window.Seconds()),
			RequestCount:              1,
			RequestsAvailableInWindow: policy.Limit - 1,
		}, nil
	}
	windowEnd := bkt.WindowStart.Add(policy.Window)
	if currentTime.After(windowEnd) {
		bkt.Requests = 1
		bkt.WindowStart = currentTime
		return &RateLimiterResponse{
			Available:                 true,
			WindowEndTime:             currentTime.Add(policy.Window),
			ExpiryInSeconds:           uint64(policy.Window.Seconds()),
			RequestCount:              1,
			RequestsAvailableInWindow: uint64(policy.Limit - 1),
		}, nil
	} else if bkt.Requests < policy.Limit {
		bkt.Requests++
		return &RateLimiterResponse{
			Available:                 true,
			WindowEndTime:             windowEnd,
			ExpiryInSeconds:           uint64(windowEnd.Sub(currentTime).Seconds()),
			RequestCount:              bkt.Requests,
			RequestsAvailableInWindow: policy.Limit - bkt.Requests,
		}, nil
	}

	return &RateLimiterResponse{
		Available:                 false,
		WindowEndTime:             windowEnd,
		ExpiryInSeconds:           uint64(windowEnd.Sub(currentTime).Seconds()),
		RequestCount:              bkt.Requests,
		RequestsAvailableInWindow: 0,
	}, nil
}

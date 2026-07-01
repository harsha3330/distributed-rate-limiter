package bucket

import (
	"fmt"
	"sync"
	"time"

	policy "github.com/harsha3330/distributed-rate-limiter/internal/policy"
)

type BucketKey struct {
	PolicyID string `json:"policy_id"`
	Subject  string `json:"subject"`
}

type Bucket struct {
	Count uint32    `json:"count"`
	Start time.Time `json:"start"`
}

type BucketStore struct {
	mu    sync.Mutex
	state map[BucketKey]*Bucket
	ps    *policy.PolicyStore
}

type BucketCheckResponse struct {
	Available         bool      `json:"available"`
	RequestCount      uint64    `json:"request_count"`
	WindowStart       time.Time `json:"window_start"`
	ResetAfterSeconds int64     `json:"reset_after_seconds"`
}

func NewBucketStore() *BucketStore {
	ps := policy.NewPolicyStore()
	return &BucketStore{
		state: make(map[BucketKey]*Bucket),
		ps:    ps,
	}
}

func (bs *BucketStore) Check(key BucketKey) (*BucketCheckResponse, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bucket, ok := bs.state[key]
	if !ok {
		return nil, fmt.Errorf("Bucket key with subject : %s and policy : %s does not exists", key.PolicyID, key.Subject)
	}
	policy, err := bs.ps.GetPolicy(key.PolicyID)
	if err != nil {
		return nil, fmt.Errorf("Bucket key's Policy : %s does not exists", key.PolicyID)
	}

	bucketExpiry := bucket.Start.Add(policy.Window)

	if time.Now().After(bucketExpiry) {
		start := time.Now()
		bs.state[key] = &Bucket{
			Count: 1,
			Start: start,
		}
		return &BucketCheckResponse{
			Available:         true,
			RequestCount:      1,
			WindowStart:       start,
			ResetAfterSeconds: int64(policy.Window.Seconds()),
		}, nil
	}

	windowStart := bs.state[key].Start
	if bucket.Count+1 < policy.Limit {
		bs.state[key].Count++
		return &BucketCheckResponse{
			Available:         true,
			RequestCount:      uint64(bs.state[key].Count),
			WindowStart:       windowStart,
			ResetAfterSeconds: int64(time.Since(windowStart)),
		}, nil
	}

	return &BucketCheckResponse{
		Available:         false,
		RequestCount:      uint64(bs.state[key].Count),
		WindowStart:       windowStart,
		ResetAfterSeconds: int64(time.Since(windowStart)),
	}, nil
}

func (bs *BucketStore) GetBucket(key BucketKey) (*Bucket, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if bucket, ok := bs.state[key]; ok {
		return bucket, nil
	}
	return nil, fmt.Errorf("Bucket with subject: %s ,policy: %s not found", key.Subject, key.PolicyID)
}

func (bs *BucketStore) DeleteBucket(key BucketKey) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if _, ok := bs.state[key]; !ok {
		return fmt.Errorf("Bucket with subject: %s ,policy: %s not found", key.Subject, key.PolicyID)
	}
	delete(bs.state, key)
	return nil
}

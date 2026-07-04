package bucket

import (
	"fmt"
	"sync"
	"time"
)

type BucketKey struct {
	PolicyID string `json:"policy_id"`
	Subject  string `json:"subject"`
}

type Bucket struct {
	Requests    uint64    `json:"requests"`
	WindowStart time.Time `json:"window_start"`
}

type BucketInfo struct {
	Requests    uint64    `json:"requests"`
	WindowStart time.Time `json:"window_start"`
	PolicyID    string    `json:"policy_id"`
	Subject     string    `json:"subject"`
}

type BucketStore struct {
	mu    sync.RWMutex
	state map[BucketKey]*Bucket
}

type BucketCheckResponse struct {
	Available         bool      `json:"available"`
	RequestCount      uint64    `json:"request_count"`
	WindowStart       time.Time `json:"window_start"`
	ResetAfterSeconds int64     `json:"reset_after_seconds"`
}

func NewBucketStore() *BucketStore {
	return &BucketStore{
		state: make(map[BucketKey]*Bucket),
	}
}

func (bs *BucketStore) ListBucket() []*BucketInfo {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	buckets := make([]*BucketInfo, 0)
	for bkey, bval := range bs.state {
		buckets = append(buckets, &BucketInfo{
			Requests:    bval.Requests,
			WindowStart: bval.WindowStart,
			PolicyID:    bkey.PolicyID,
			Subject:     bkey.Subject,
		})
	}
	return buckets
}

func (bs *BucketStore) GetOrCreateBucket(key BucketKey) (*Bucket, bool, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bucket, ok := bs.state[key]
	if ok {
		return bucket, false, nil
	}
	bucket = &Bucket{
		Requests:    0,
		WindowStart: time.Now(),
	}
	bs.state[key] = bucket
	return bucket, true, nil
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

func (bs *BucketStore) UpdateBucket(key BucketKey, bucket *Bucket) error {
	if _, ok := bs.state[key]; !ok {
		return fmt.Errorf("bucket with subject %q and policy %q not found", key.Subject, key.PolicyID)
	}
	bs.state[key] = bucket
	return nil
}

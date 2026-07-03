package policy

import (
	"fmt"
	"sync"
	"time"
)

type Policy struct {
	ID     string        `json:"id"`
	Window time.Duration `json:"window"`
	Limit  uint64        `json:"limit"`
}

type PolicyStore struct {
	mu    sync.RWMutex
	state map[string]*Policy
}

func NewPolicyStore() *PolicyStore {
	return &PolicyStore{
		state: make(map[string]*Policy),
	}
}

func (ps *PolicyStore) ListPolicy() []*Policy {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	policies := make([]*Policy, 0)
	for _, policy := range ps.state {
		policies = append(policies, policy)
	}
	return policies
}

func (ps *PolicyStore) GetPolicy(ID string) (*Policy, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if policy, ok := ps.state[ID]; ok {
		return policy, nil
	}
	return nil, fmt.Errorf("Policy with ID: %s not found", ID)
}

func (ps *PolicyStore) CreatePolicy(p *Policy) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.state[p.ID]; ok {
		return fmt.Errorf("Policy with ID: %s already exists", p.ID)
	}
	ps.state[p.ID] = p
	return nil
}

func (ps *PolicyStore) UpdatePolicy(p *Policy) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.state[p.ID]; !ok {
		return fmt.Errorf("Policy with ID: %s does not exists", p.ID)
	}
	ps.state[p.ID] = p
	return nil
}

func (ps *PolicyStore) DeletePolicy(ID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.state[ID]; !ok {
		return fmt.Errorf("Policy with ID: %s does not exists", ID)
	}
	delete(ps.state, ID)
	return nil
}

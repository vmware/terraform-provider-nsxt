package nsxt

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type providerRuntimeContext struct {
	mu   sync.RWMutex
	cond *sync.Cond

	configID    atomic.Value // stores string
	initialized atomic.Bool
}

func newProviderRuntimeContext() *providerRuntimeContext {
	r := &providerRuntimeContext{}
	// Cond must use a Locker with Lock/Unlock (not RLock/RUnlock),
	// so we bind it to the RWMutex write lock.
	r.cond = sync.NewCond(&r.mu)
	// Initialize atomic.Value to avoid panics on Load before first Store.
	r.configID.Store("")
	return r
}

func InitConfigIDOnce(m interface{}, configID string) error {
	if configID == "" {
		return fmt.Errorf("configID is empty")
	}
	c := m.(nsxtClients)
	r := c.Runtime

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized.Load() {
		return nil
	}
	r.configID.Store(configID)
	r.initialized.Store(true)
	r.cond.Broadcast()
	return nil
}

func GetConfigID(m interface{}) (string, error) {
	c := m.(nsxtClients)
	r := c.Runtime

	// Fast path: once initialized, callers do not lock at all.
	if r.initialized.Load() {
		id, _ := r.configID.Load().(string)
		if id == "" {
			return "", fmt.Errorf("provider configID initialized but empty")
		}
		return id, nil
	}

	// Slow path: wait for the singleton resource to initialize configID.
	r.mu.Lock()
	defer r.mu.Unlock()
	for !r.initialized.Load() {
		r.cond.Wait()
	}
	id, _ := r.configID.Load().(string)
	if id == "" {
		return "", fmt.Errorf("provider configID initialized but empty")
	}
	return id, nil
}

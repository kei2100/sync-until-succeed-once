package sync

import (
	"sync"
	"sync/atomic"
)

// UntilSucceedOnce is similar to sync.Once, but perform the action until it succeeds
type UntilSucceedOnce struct {
	m    sync.Mutex
	done uint32
}

// Do is similar to sync.Once.Do, but perform the action until it succeeds
func (o *UntilSucceedOnce) Do(fn func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		if err := fn(); err != nil {
			return err
		}
		atomic.StoreUint32(&o.done, 1)
	}
	return nil
}

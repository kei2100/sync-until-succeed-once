package sync_test

import (
	"fmt"
	gosync "sync"
	"sync/atomic"

	"github.com/kei2100/sync-until-succeed-once"
)

func ExampleUntilSucceedOnce_Do() {
	var usOnce sync.UntilSucceedOnce
	cache := cacheSomething{errorsOccursCount: 5}

	const n = 10
	var wg gosync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := usOnce.Do(cache.sometimesErrorsOccurs); err != nil {
				// handle errors as needed.
			}
		}()
	}
	wg.Wait()

	// Output: 6
	fmt.Println(cache.callCount)
}

type cacheSomething struct {
	errorsOccursCount int32
	callCount         int32
}

func (c *cacheSomething) sometimesErrorsOccurs() error {
	incr := atomic.AddInt32(&c.callCount, 1)
	if incr <= c.errorsOccursCount {
		return fmt.Errorf("return errors until incr(%d) <= errrosOccursCount(%d)", incr, c.errorsOccursCount)
	}
	return nil
}

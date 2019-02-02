package helper

import (
	"time"
)

func Debounce(f func()) func() {
	waitingInvoke := false
	return func() {
		invokeTime := time.After(300 * time.Millisecond)
		if !waitingInvoke {
			waitingInvoke = true
			go func() {
				for {
					select {
					case <-invokeTime:
						f()
						waitingInvoke = false
						return
					}
				}
			}()
		}
	}
}

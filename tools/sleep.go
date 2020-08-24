package tools

import "time"

// Sleep sleeps the current thread for a certain duration d.
func Sleep(d time.Duration) {
	time.Sleep(d)
}

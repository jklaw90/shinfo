package utils

import (
	"time"

	"golang.org/x/time/rate"
)

// LimitWithPreviousErr used for health checks and checks where you would only want to allow the check every so often
// will also return previous err if limiter not allowed
func LimitOncePer(durration time.Duration, fn func() error) func() error {
	var err error
	limiter := rate.NewLimiter(rate.Every(durration), 1)
	return func() error {
		if limiter.Allow() {
			err = fn()
		}
		return err
	}
}

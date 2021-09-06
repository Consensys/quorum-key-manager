package common

import (
	"context"
)

// Runnable manages running long living task
type Runnable interface {
	// Start long living task in a parallel goroutine

	// It MUST return an error if and only if it failed at starting the long living task
	Start(context.Context) error

	// Stop gracefully interupts long living task

	// Stop SHOULD make sure underlying long living task has gracefully interupted execution before returning

	// In case context timeouts or is canceled Stop MUST
	// 1. [optional] Try to kill long living task (MUST not block on doing this)
	// 2. Return immediately with an error

	// In any other case Stop MUST not return before the long live task
	// has gracefully interupted execution
	Stop(context.Context) error

	// Close clean the runnable

	// Init, Start, Stop, MUST NOT be called after Close
	Close() error

	// Error returns any possible error met by the runnable

	// It MAY be
	// - an error that raised Init, Start, Stop or Close
	// - an error that raised on the long living task which force prematured Stop
	Error() error
}

// Checkable allows to
type Checkable interface {
	ID() string

	// Live MUST return an error if the long living task is not running otherwise nil
	CheckLiveness(context.Context) error

	// Ready MUST return an error if the long living task is not running otherwise nil
	CheckReadiness(context.Context) error
}

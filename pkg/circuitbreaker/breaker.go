package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var (
	ErrCircuitOpen = errors.New("circuit is open")
)

type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	successThreshold int

	state           State
	failureCount    int
	successCount    int
	lastStateChange time.Time

	mutex sync.RWMutex
}

type Options struct {
	FailureThreshold int
	ResetTimeout     time.Duration
	SuccessThreshold int
}

func DefaultOptions() Options {
	return Options{
		FailureThreshold: 5,
		ResetTimeout:     10 * time.Second,
		SuccessThreshold: 2,
	}
}

func NewCircuitBreaker(opts Options) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: opts.FailureThreshold,
		resetTimeout:     opts.ResetTimeout,
		successThreshold: opts.SuccessThreshold,
		state:            StateClosed,
		lastStateChange:  time.Now(),
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.AllowRequest() {
		return ErrCircuitOpen
	}

	err := fn()

	cb.RecordResult(err == nil)

	return err
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastStateChange) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.lastStateChange = time.Now()
			cb.failureCount = 0
			cb.successCount = 0
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) RecordResult(success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateClosed:
		if !success {
			cb.failureCount++
			if cb.failureCount >= cb.failureThreshold {
				cb.moveToOpenLocked()
			}
		} else {
			cb.failureCount = 0
		}
	case StateHalfOpen:
		if success {
			cb.successCount++
			if cb.successCount >= cb.successThreshold {
				cb.moveToClosedLocked()
			}
		} else {
			cb.moveToOpenLocked()
		}
	}
}

func (cb *CircuitBreaker) moveToOpenLocked() {
	cb.state = StateOpen
	cb.lastStateChange = time.Now()
	cb.failureCount = 0
	cb.successCount = 0
}

func (cb *CircuitBreaker) moveToClosedLocked() {
	cb.state = StateClosed
	cb.lastStateChange = time.Now()
	cb.failureCount = 0
	cb.successCount = 0
}

func (cb *CircuitBreaker) moveToHalfOpen() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateOpen {
		cb.state = StateHalfOpen
		cb.lastStateChange = time.Now()
		cb.failureCount = 0
		cb.successCount = 0
	}
}

func (cb *CircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return cb.state
}

package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(Options{
		FailureThreshold: 2,
		ResetTimeout:     100 * time.Millisecond,
		SuccessThreshold: 1,
	})

	// Test initial state
	if cb.GetState() != StateClosed {
		t.Errorf("Initial state should be Closed, got %v", cb.GetState())
	}

	// Test transition to open
	for i := 0; i < 2; i++ {
		err := cb.Execute(func() error {
			return errors.New("failure")
		})
		if err == nil {
			t.Error("Expected error but got nil")
		}
	}

	if cb.GetState() != StateOpen {
		t.Errorf("State should be Open, got %v", cb.GetState())
	}

	// Test rejection while open
	err := cb.Execute(func() error {
		return nil
	})
	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Test half-open state
	err = cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Should be closed after success
	if cb.GetState() != StateClosed {
		t.Errorf("State should be Closed, got %v", cb.GetState())
	}
}

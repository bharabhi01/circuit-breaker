package main

import (
	"fmt"
	"time"

	"github.com/bharabhi01/circuit-breaker/pkg/circuitbreaker"
)

func main() {
	options := circuitbreaker.DefaultOptions()
	options.ResetTimeout = 3 * time.Second
	cb := circuitbreaker.NewCircuitBreaker(options)

	callCount := 0

	unstableService := func() error {
		callCount++

		if callCount <= 10 {
			fmt.Println("Service call failed")
			return fmt.Errorf("service is unavailable")
		}

		fmt.Println("Service call succeeded")
		return nil
	}

	for i := 1; i <= 15; i++ {
		fmt.Printf("Attempt %d:\n", i)
		fmt.Println("Circuit State: ", getStateName(cb.GetState()))

		err := cb.Execute(unstableService)
		if err != nil {
			if err == circuitbreaker.ErrCircuitOpen {
				fmt.Println("Circuit is open, request rejected")
			} else {
				fmt.Println("Service call failed:", err)
			}
		} else {
			fmt.Println("Service call succeeded")
		}

		time.Sleep(1 * time.Second)
	}
}

func getStateName(state circuitbreaker.State) string {
	switch state {
	case circuitbreaker.StateClosed:
		return "CLOSED"
	case circuitbreaker.StateOpen:
		return "OPEN"
	case circuitbreaker.StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

# Circuit Breaker

A simple, idiomatic Go implementation of the Circuit Breaker pattern.

## Usage

```go
// Create a circuit breaker with default settings
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultOptions())

// Use it to protect a function call
err := cb.Execute(func() error {
    // Call to an external service or any function that might fail
    return externalService()
})

if err == circuitbreaker.ErrCircuitOpen {
    // Circuit is open, take alternative action
} else if err != nil {
    // Service returned an error
} else {
    // Service call succeeded
}
```

## Options

- `FailureThreshold`: Number of consecutive failures before opening the circuit
- `ResetTimeout`: How long the circuit stays open before allowing a test request
- `SuccessThreshold`: Number of successful test requests needed to close the circuit

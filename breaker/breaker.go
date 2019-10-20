package breaker

import (
        "github.com/ae-gis/volt/log"
        "github.com/afex/hystrix-go/hystrix"
)

type CircuitBreaker struct {
        name          string
        maxConcurrent int
        interval      int
        timeout       int

        fallbackFunc fallbackFunc
}

type fallbackFunc func(error) error

// SetCommandBreaker the circuit breaker
func NewBreaker(name string, timeout, maxConcurrent int, args ...interface{}) *CircuitBreaker {
        cb := &CircuitBreaker{
                name:          name,
                maxConcurrent: maxConcurrent,
                timeout:       timeout,
        }
        if len(args) == 1 {
                switch args[0].(type) {
                case fallbackFunc:
                        cb.fallbackFunc = args[0].(fallbackFunc)
                }
        }

        hystrix.ConfigureCommand(cb.name, hystrix.CommandConfig{
                MaxConcurrentRequests: cb.maxConcurrent,
                Timeout:               cb.timeout,
                ErrorPercentThreshold: 25,
        })

        return cb
}

// callBreaker command circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) (err error) {
        if cb.name == "" {
                return fn()
        }

        err = hystrix.Do(cb.name, func() error {
                return fn()
        }, nil)

        if err != nil {
                log.Error("Call Breaker",
                        log.Field("Hystrix Do", err))
        }
        return err
}

package gopool

import "errors"

var (
	// configuration parameter is illegal, 0 or negaitve for buffer size etc.
	ERR_INVALID_CONFIG error = errors.New("invalid config parameter provided.")
	// pool.Start() called on an already runnig pool
	ERR_START_ON_RUNNING_POOL error = errors.New("cannot start pool: already running")
	// pool.Shutdown() called on an stopped / not started pool
	ERR_STOP_ON_CLOSED_POOL error = errors.New("cannot shutdown pool: already stopped")
	// Task submitted was nil
	ERR_NIL_WORK error = errors.New("cannot submit nil task pool")
	// Task submitted in closed pool
	ERR_SUBMIT_IN_CLOSED_POOL error = errors.New("cannot submit task: pool not running")
)

// Work Work is a unit of work to be submitted to pool
type Work interface {
	// Do implement this func as the processing block of Task;
	//  returns the Result
	Do() Result
}

// Result Result is the result of Task execution.
type Result interface {
	// Get Get func needs to implemented to get
	// result back from the processing result of Task
	Get() any
}

// PoolLogger PoolLogger is logger for the Pool.
// Supply an implementation which satisfies this interface
type PoolLogger interface {
	// Debug prints debug message with supplied args
	Debug(msg string, args ...any)
	// Info prints info message with supplied args
	Info(msg string, args ...any)
	// Warn prints warn message with supplied args
	Warn(msg string, args ...any)
	// Error prints error message with supplied args
	Error(msg string, args ...any)
}

// NOOPLogger sample logger, logs nothing
type NOOPLogger struct{}

func (p *NOOPLogger) Debug(msg string, args ...any) {}
func (p *NOOPLogger) Info(msg string, args ...any)  {}
func (p *NOOPLogger) Warn(msg string, args ...any)  {}
func (p *NOOPLogger) Error(msg string, args ...any) {}

var _ PoolLogger = &NOOPLogger{}

package gopool

// PoolLogger PoolLogger is a logger for the Pool.
// Supply an implementation which satisfies this interface
type PoolLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

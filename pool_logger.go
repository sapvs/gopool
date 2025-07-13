package gopool

// PoolLogger PoolLogger is a logger for the Pool.
// Supply an implementation which satisfies this interface
type PoolLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NOOPLogger struct{}

func (p *NOOPLogger) Debug(msg string, args ...any) {}

func (p *NOOPLogger) Info(msg string, args ...any) {}

func (p *NOOPLogger) Warn(msg string, args ...any) {}

func (p *NOOPLogger) Error(msg string, args ...any) {}

var _ PoolLogger = &NOOPLogger{}

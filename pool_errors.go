package gopool

import "errors"

var (
	// conifguration parameter is illegal, 0 or negaitvve for buffer size etc.
	ERR_INVALID_CONFIG error = errors.New("invalid config parameter provided.")
	// pool.Start() called on an already runnig pool
	ERR_START_ON_RUNNING_POOL error = errors.New("cannot start a running pool")
	// pool.Shutdown() called on an stopped / not started pool
	ERR_STOP_ON_CLOSED_POOL error = errors.New("cannot stop a non running pool")
	// Task submitted was nil
	ERR_NIL_TASK error = errors.New("cannot submit nil task pool")
	// Task submitted in closed pool
	ERR_SUBMIT_IN_CLOSED_POOL error = errors.New("cannot submit task to closed pool")
)

package gopool

// Task Task is a unit of work to be submitted to pool
type Task interface {
	// Do implement this func as the processing block of Task;
	//  returns the Result
	Do() Result
}

// Result Result is the result of Task execution.
type Result interface {
	// Result Result func needs to implemented to get
	// result back from the processing result of Task
	Result() any
}

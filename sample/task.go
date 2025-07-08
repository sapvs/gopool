package sample

import (
	"strings"

	"github.com/sapvs/gopool"
)

type ATask struct {
	Message string
}

func (p *ATask) Do() gopool.Result {
	return &AResult{strings.ToUpper(p.Message)}
}

type AResult struct {
	result string
}

func (p *AResult) Result() any {
	return p.result
}

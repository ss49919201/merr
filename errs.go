package errs

import (
	"errors"
	"sync"
)

var formatter = func(e *errs) string {
	var result string
	suffix := ",\n"
	for i, err := range e.Errors {
		if len(e.Errors)-1 == i {
			suffix = ""
		}
		result += err.Error() + suffix
	}
	return result
}

type Errs interface {
	Append(err error)
	Error() string
	Is(target error) bool
	As(target interface{}) bool
}

type errs struct {
	mx     sync.Mutex
	Errors []error
}

func New() Errs {
	return &errs{
		mx:     sync.Mutex{},
		Errors: nil,
	}
}

func (e *errs) Error() string {
	return formatter(e)
}

func (e *errs) Append(err error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.Errors == nil {
		e.Errors = make([]error, 0)
	}
	e.Errors = append(e.Errors, err)
}

func (e *errs) Is(target error) bool {
	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *errs) As(target interface{}) bool {
	for _, err := range e.Errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

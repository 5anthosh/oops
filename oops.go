package oops

import (
	"runtime"
	"strings"
)

//Error is context for error
type Error struct {
	error
	File     string
	Line     int
	FuncName string
	Module   string
}

func (err Error) Error() string {
	return err.error.Error()
}

//T stores error with info like
// the file name and line number of the
// source code corresponding to the where function got called
func T(err error) Error {
	var err1 Error
	if err == nil {
		return err1
	}
	switch err.(type) {
	case Error:
		return err.(Error)
	}
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return newError(err, f, pc)
}

func newError(err error, f *runtime.Func, pc []uintptr) Error {
	var err1 Error
	if f != nil {
		file, line := f.FileLine(pc[0])
		funN := f.Name()
		funI := strings.Split(funN, ".")
		mod := funI[0]
		err1 = Error{
			error:    err,
			File:     file,
			Line:     line,
			FuncName: funN,
			Module:   mod,
		}
	}
	return err1
}

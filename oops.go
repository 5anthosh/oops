//Package oops makes errors in Go traceable. It provies traceback function to get more information when you return error from the function
package oops

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var emptyString = ""

//ErrorTraceFormat prints trace detail in a format
type ErrorTraceFormat func(funcName string, lineNo int, file string) string

//ErrorHeaderFormat prints error detail in a format
type ErrorHeaderFormat func(err string, info string) string

var defaultTraceFormat = "\n \t at %s line %d %s "
var defaultErrorFormat = "🔴  Error : %s \n%s "

func dTraceFormat(funcName string, lineNo int, file string) string {
	return fmt.Sprintf(defaultTraceFormat, funcName, lineNo, file)
}

func dErrorFormat(err string, info string) string {
	if len(info) > 0 {
		return fmt.Sprintf(defaultErrorFormat, err, "ℹ️  Info  : "+info)
	}
	return fmt.Sprintf(defaultErrorFormat, err, "")
}

//Error is a error with more information
type Error struct {
	error
	info        string
	stackTrace  []stack
	skip        int
	traceFormat ErrorTraceFormat
	errorFormat ErrorHeaderFormat
}

//stack stores single stack information
type stack struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	FuncName string `json:"func_name,omitempty"`
}

//JSON converts Error into json format to for structured logging
func (err *Error) JSON() map[string]interface{} {
	if err == nil {
		return nil
	}
	json := make(map[string]interface{})
	json["error"] = err.error
	json["info"] = err.info
	json["stack_trace"] = err.stackTrace
	return json
}

//TraceFormat registers Errorformat function to print trace in a format
func (err *Error) TraceFormat(f ErrorTraceFormat) *Error {
	if err == nil {
		return nil
	}
	err.traceFormat = f
	return err
}

//ErrorFormat registers Errorformat function to print error in a format
func (err *Error) ErrorFormat(f ErrorHeaderFormat) *Error {
	if err == nil {
		return nil
	}
	err.errorFormat = f
	return err
}

func (s stack) format(f string) string {
	return fmt.Sprintf(f, s.FuncName, s.Line, s.File)
}
func (err *Error) Error() string {
	if err == nil {
		return emptyString
	}
	return err.errorWithSkip(err.skip)
}

//Info lets you add more information about the error
func (err *Error) Info(value string) *Error {
	if err == nil {
		return nil
	}
	err.info = value
	return err
}

//Skip skips n functions from bottom of the stack
func (err *Error) Skip(n int) *Error {
	if err == nil {
		return nil
	}
	if n > len(err.stackTrace) {
		n = len(err.stackTrace)
	}
	err.skip = n
	return err
}

//Line sets line no where error occured
func (err *Error) Line(value int) *Error {
	if err == nil {
		return nil
	}
	err.stackTrace[0].Line = value
	return err
}

//Func sets function name where error occured
func (err *Error) Func(value string) *Error {
	if err == nil {
		return nil
	}
	err.stackTrace[0].FuncName = value
	return err
}

func (err *Error) errorWithSkip(skip int) string {
	if err == nil {
		return emptyString
	}
	if skip > len(err.stackTrace) {
		skip = len(err.stackTrace)
	}
	var sb strings.Builder
	st := err.error.Error()
	sb.WriteString(err.errorFormat(st, err.info))
	for _, stack := range err.stackTrace[:len(err.stackTrace)-skip] {
		sb.WriteString(err.traceFormat(stack.FuncName, stack.Line, stack.File))
	}
	return sb.String()
}

//Origin prints where error got originated, not the trace
func (err *Error) Origin() string {
	if err == nil {
		return emptyString
	}
	return err.errorWithSkip(len(err.stackTrace) - 1)
}

//T add error with more information like stacktrace
//corresponding to the where function got called
func T(err error) *Error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case *Error:
		return err.(*Error)
	}
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc)
	var st []stack
	for {
		f, more := frames.Next()
		if !more {
			break
		}
		st = append(st, stack{
			File:     formatFileName(f.File),
			Line:     f.Line,
			FuncName: f.Function,
		})
	}
	return &Error{error: err, stackTrace: st, traceFormat: dTraceFormat, errorFormat: dErrorFormat}
}

func formatFileName(fileName string) string {
	homDir, err := os.UserHomeDir()
	if err == nil {
		fileName = strings.TrimLeft(fileName, homDir)
	}
	return fileName
}

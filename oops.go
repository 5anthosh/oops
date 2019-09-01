//Package oops makes errors in Go traceable.
//It provides traceback function to get more information when you return error from the function
package oops

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var emptyString = ""

//ErrorFormat let us specfiy format for error and stacktrace
type ErrorFormat func(err string, info string, stacktrace []Stack) string

var defaultTraceFormat = "\n \t at %s line %d %s "
var defaultErrorFormat = "🔴  Error : %s \n%s "

func dTraceFormat(funcName string, lineNo int, file string) string {
	return fmt.Sprintf(defaultTraceFormat, funcName, lineNo, file)
}

func dErrorHeaderFormat(err string, info string) string {
	if len(info) > 0 {
		return fmt.Sprintf(defaultErrorFormat, err, "ℹ️  Info  : "+info)
	}
	return fmt.Sprintf(defaultErrorFormat, err, "")
}

func dErrorFormat(err string, info string, stackTrace []Stack) string {
	var sb strings.Builder
	errF := dErrorHeaderFormat(err, info)
	sb.WriteString(errF)
	for _, trace := range stackTrace {
		sb.WriteString(dTraceFormat(trace.FuncName, trace.Line, trace.File))
	}
	return sb.String()
}

//Error is a error with more information
type Error struct {
	error
	info       string
	stackTrace []Stack
	skip       int
	format     ErrorFormat
}

//Stack stores single Stack information
type Stack struct {
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

//Format registers format for error strings
func (err *Error) Format(f ErrorFormat) *Error {
	if err == nil {
		return nil
	}
	err.format = f
	return err
}

func (s Stack) format(f string) string {
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
	return err.format(err.error.Error(), err.info, err.stackTrace)
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
	var st []Stack
	for {
		f, more := frames.Next()
		if !more {
			break
		}
		st = append(st, Stack{
			File:     formatFileName(f.File),
			Line:     f.Line,
			FuncName: f.Function,
		})
	}
	return &Error{error: err, stackTrace: st, format: dErrorFormat}
}

func formatFileName(fileName string) string {
	homDir, err := os.UserHomeDir()
	if err == nil {
		fileName = strings.TrimLeft(fileName, homDir)
	}
	return fileName
}

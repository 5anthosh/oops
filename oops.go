package oops

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

//ErrorFormat prints trace detail in a format
type ErrorFormat func(funcName string, lineNo int, file string) string

var defaultErrorFormat = "\n at %s line %d %s "

func dErrorFormat(funcName string, lineNo int, file string) string {
	return fmt.Sprintf(defaultErrorFormat, funcName, lineNo, file)
}

//Error is a error with more information
type Error struct {
	error
	info       string
	stackTrace []Stack
	skip       int
	errFormat  ErrorFormat
}

//Stack stores single stack information
type Stack struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	FuncName string `json:"func_name,omitempty"`
}

//JSON converts Error into json format to for structured logging
func (err Error) JSON() map[string]interface{} {
	json := make(map[string]interface{})
	json["error"] = err.error
	json["info"] = err.info
	json["stack_trace"] = err.stackTrace
	return json
}

//Format registers Errorformat function to print trace in a format
func (err Error) Format(f ErrorFormat) Error {
	err.errFormat = f
	return err
}

func (s Stack) format(f string) string {
	return fmt.Sprintf(f, s.FuncName, s.Line, s.File)
}
func (err Error) Error() string {
	return err.errorWithSkip(err.skip)
}

//Info lets you add more information about the error
func (err Error) Info(value string) Error {
	err.info = value
	return err
}

//Skip skips n functions from bottom of the stack
func (err Error) Skip(n int) Error {
	if n > len(err.stackTrace) {
		n = len(err.stackTrace)
	}
	err.skip = n
	return err
}

//Line sets line no where error occured
func (err Error) Line(value int) Error {
	err.stackTrace[0].Line = value
	return err
}

//Func sets function name where error occured
func (err Error) Func(value string) Error {
	err.stackTrace[0].FuncName = value
	return err
}

func (err Error) errorWithSkip(skip int) string {
	if skip > len(err.stackTrace) {
		skip = len(err.stackTrace)
	}
	var sb strings.Builder
	st := err.error.Error()
	sb.WriteString("🔴  Error : ")
	sb.WriteString(st)
	sb.WriteByte('\n')

	if len(err.info) > 0 {
		sb.WriteString("ℹ️  Info  : ")
		sb.WriteString(err.info)
		sb.WriteByte('\n')
	}
	for _, stack := range err.stackTrace[:len(err.stackTrace)-skip] {
		sb.WriteString(err.errFormat(stack.FuncName, stack.Line, stack.File))
	}
	return sb.String()
}

//Origin prints where error got originated, not the trace
func (err Error) Origin() string {
	return err.errorWithSkip(len(err.stackTrace) - 1)
}

//T add error with more information like stacktrace
//corresponding to the where function got called
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
	return Error{error: err, stackTrace: st, errFormat: dErrorFormat}
}

func formatFileName(fileName string) string {
	homDir, err := os.UserHomeDir()
	if err == nil {
		fileName = strings.TrimLeft(fileName, homDir)
	}
	return fileName
}

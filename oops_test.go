package oops

import (
	"errors"
	"strings"
	"testing"
)

// TODO: write the tests for Error funcs

var errorMessage = "testing"
var infoMessage = "info message"
var file = "github.com/5anthosh/oops/oops_test.go"

func test2() *Error {
	return test()
}
func test() *Error {
	return T(errors.New(errorMessage)).Info(infoMessage)
}

func TestT(t *testing.T) {
	err := test2()
	if err.error.Error() != errorMessage {
		t.Errorf("want %q instead got %q", errorMessage, err.error.Error())
	}
	if err.info != "info message" {
		t.Errorf("want %q instead got %q", infoMessage, err.info)
	}

	errortrace := err.stackTrace
	top := errortrace[0]
	teststack(top, t, "github.com/5anthosh/oops.test", 19, file)
	belowTest := errortrace[1]
	teststack(belowTest, t, "github.com/5anthosh/oops.test2", 16, file)
	testFunc := errortrace[2]
	teststack(testFunc, t, "github.com/5anthosh/oops.TestT", 23, file)
}

func teststack(s stack, t *testing.T, funcName string, line int, fileName string) {
	if s.FuncName != funcName {
		t.Errorf("want %q instead got %q", funcName, s.FuncName)
	}
	if s.Line != line {
		t.Errorf("want %d instead got %d", line, s.Line)
	}
	if !strings.HasSuffix(s.File, fileName) {
		t.Errorf("%q does not have suffix %q", s.File, fileName)
	}
}

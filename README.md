# oops

oops makes errors in go traceable.
It provies traceback function to get more information when you return error from the function

> The convention says that an error should either be handled (whatever that means) or returned to the caller.
> But with abtraction, it is difficult to find where the error got originated
> so Traceable function returns error with more information

## example

```go
package main

import (
	"errors"

	"github.com/5anthosh/oops"
)

func main() {
	err := f1().(oops.Error)
	println("FileName ", err.File)
	println("Line no", err.Line)
	println("FuncName", err.FuncName)
	println("Module", err.Module)
}

func f1() error {
	return f2()
}

func f2() error {
	return f3()
}

func f3() error {
	return oops.T(errors.New("dummy one")) // Line number 26
}
```

### Run the program

```sh
$ go run test.go
FileName test.go
Line no 26
FuncName main.f3
Module main
```

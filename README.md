# oops 🙊

oops makes errors in Go traceable.
It provies traceback function to get more information when you return error from the function

> The convention says that an error should either be handled (whatever that means) or returned to the caller.
> But with more abstraction, it is difficult to find where the error got originated
> so Traceable function returns error with more information

## Example

```go
package main

import (
	"errors"

	"github.com/5anthosh/oops"
)

func main() {
	err := func1().(oops.Error)
	println(err.Skip(1).Error())
}

func func1() error {
	return func2()
}

func func2() error {
	return func3()
}

func func3() error {
	return oops.T(errors.New("dummy one")).Info("this is just testing")
}

```

### Run the program

```sh
$ go run test.go
🔴  Error : dummy one
ℹ️   Info  : this is just testing
         at main.func3 line 23 Desktop/Files/test.go
         at main.func2 line 19 Desktop/Files/test.go
         at main.func1 line 15 Desktop/Files/test.go
         at main.main line 10 Desktop/Files/test.go
```
- [![GoDoc](https://godoc.org/github.com/5anthosh/oops?status.svg)](https://godoc.org/github.com/5anthosh/oops)

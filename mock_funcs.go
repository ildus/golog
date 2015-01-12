package golog

import "os"
import "fmt"

var osExit func(code int) = os.Exit
var mocked bool = false

func mockFuncs() {
	osExit = func(code int) {
		fmt.Printf("Fake exit with code %d\n", code)
	}
}

func useStdFuncs() {
	osExit = os.Exit
}

package main

import (
	"github.com/cuiweixie/toylang/eval"
)

func main() {
	testFileName := "D:\\code\\go\\gopath\\src\\github.com\\cuiweixie\\toylang\\cmd\\basic\\test.toy"
	eval.EvalFile(testFileName)
}

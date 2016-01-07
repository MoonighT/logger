package main

import (
	"github.com/MoonighT/logger"
)

func main() {
	logger.Init("./log/test.log", 2, 3600*24, 1024*1024*10, 3)
	test := "abc"
	logger.Detailf("this is a test %s", test)
}

package main

import "fmt"

var (
	ErrNonDir = fmt.Errorf("Non directory")
	ErrNoConf = fmt.Errorf("rcc.json is missing")
)

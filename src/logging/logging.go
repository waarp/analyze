package logging

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug   = log.New(ioutil.Discard, "DEBUG  : ", 0)
	Info    = log.New(os.Stderr, "INFO   : ", 0)
	Warning = log.New(os.Stderr, "WARNING: ", 0)
	Error   = log.New(os.Stderr, "ERROR  : ", 0)
)

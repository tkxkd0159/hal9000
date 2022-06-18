package types

import "os"

type Code int

const (
	EXIT Code = iota
	KEEP
)

type Fstream struct {
	In  *os.File
	Out *os.File
	Err *os.File
}

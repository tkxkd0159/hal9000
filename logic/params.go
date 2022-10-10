package logic

import (
	"time"
)

const (
	ReQueryDelay    = time.Duration(1)
	SeqRecoverDelay = time.Duration(4)
	IBCDelay        = time.Duration(30)
)

const (
	QueryErrPrefix = "[QUERY ERROR] : "
)

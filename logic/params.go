package logic

import (
	"time"
)

const (
	ReQueryDelay = time.Second * 1
	IBCDelay     = time.Second * 30 // This value must be set higher than the ibc timeout.
)

const (
	QueryErrPrefix = "[QUERY ERROR] : "
)

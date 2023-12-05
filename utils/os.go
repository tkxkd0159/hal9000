package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/tkxkd0159/HAL9000/utils/types"
)

var ToB64Str = base64.StdEncoding.EncodeToString

func CloseFds(fds ...*os.File) {
	for _, fd := range fds {
		err := fd.Close()
		CheckErr(err, "", 1)
	}
}

func CheckErr(err error, moreMsg string, action types.Code) {
	LogErrWithFd(os.Stderr, err, moreMsg, action)
}

func LogErrWithFd(fd *os.File, err error, msg string, action types.Code) {
	if err != nil {
		var logger *log.Logger
		_, f, l, ok := runtime.Caller(1)
		if ok {
			logger = log.New(fd, fmt.Sprintf(" ⛈ : %s:%d ", f, l), log.LstdFlags)
		} else {
			logger = log.New(fd, " ⛈ : ", log.LstdFlags)
		}
		logger.Printf("\n %s: %v\n", msg, err)

		if action == types.EXIT {
			panic("process panic after logging")
		}
	}
}

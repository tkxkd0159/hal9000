package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/Carina-labs/HAL9000/utils/types"
)

var (
	ToB64Str = base64.StdEncoding.EncodeToString
)

func CloseFds(fds ...*os.File) {
	for _, fd := range fds {
		err := fd.Close()
		CheckErr(err, "", 1)
	}
}

func CheckErr(err error, moreMsg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			panic(fmt.Sprintf("%s: \n %+v\n", moreMsg, err))
		}
	case types.KEEP:
		if err != nil {
			log.Printf("%s: \n %+v\n", moreMsg, err)
		}
	}
}

func LogErrWithFd(fd *os.File, err error, msg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			fmt.Fprintf(fd, "\n %s: \n %v", msg, err)
			panic(err)
		}
	case types.KEEP:
		if err != nil {
			var logger *log.Logger
			_, f, l, ok := runtime.Caller(1)
			if ok {
				logger = log.New(fd, fmt.Sprintf(" ⛈ : %s:%d ", f, l), log.LstdFlags)
			} else {
				logger = log.New(fd, " ⛈ : ", log.LstdFlags)
			}
			logger.Printf("\n %s: %v\n", msg, err)
		}
	}
}

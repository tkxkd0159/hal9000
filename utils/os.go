package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/Carina-labs/HAL9000/utils/types"
	"log"
	"os"
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
			panic(fmt.Sprintf("%s: \n %v", moreMsg, err))

		}
	case types.KEEP:
		if err != nil {
			log.Printf("%s: \n %v\n", moreMsg, err)
		}
	}
}

func LogErrWithFd(fd *os.File, err error, msg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			panic(fmt.Sprintf("%s: \n %v", msg, err))
		}
	case types.KEEP:
		if err != nil {
			l := log.New(fd, "ERROR (check) : ", log.Llongfile|log.LstdFlags)
			l.Printf("%s: \n %v\n", msg, err)
		}
	}
}

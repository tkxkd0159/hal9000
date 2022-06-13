package utils

import (
	"github.com/Carina-labs/HAL9000/utils/types"
	"log"
	"os"
)

func CheckErr(err error, msg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			log.Fatalf("%s: \n %v", msg, err)
		}
	case types.KEEP:
		if err != nil {
			log.Printf("%s: \n %v", msg, err)
		}
	}
}

func SetPipe(store *os.File) (*os.File, *os.File, error) {
	store = os.Stdout
	return os.Pipe()
}

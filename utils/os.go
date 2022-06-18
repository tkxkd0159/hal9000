package utils

import (
	"errors"
	"fmt"
	"github.com/Carina-labs/HAL9000/utils/types"
	"io/fs"
	"log"
	"os"
	"path"
)

func CheckErr(err error, msg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			panic(fmt.Sprintf("%s: \n %v", msg, err))

		}
	case types.KEEP:
		if err != nil {
			log.Printf("%s: \n %v\n", msg, err)
		}
	}
}

func CheckErrWithFP(fp *os.File, err error, msg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			panic(fmt.Sprintf("%s: \n %v", msg, err))
		}
	case types.KEEP:
		if err != nil {
			l := log.New(fp, "ERROR (check) : ", log.Llongfile|log.LstdFlags)
			l.Printf("%s: \n %v\n", msg, err)
		}
	}
}

func SetDir(dir string) (string, error) {
	ckrDir, err := os.Getwd()
	CheckErr(err, "cannot get working directory", 0)
	var fm fs.FileMode = 0740

	dir = path.Join(ckrDir, dir)
	err = os.Mkdir(dir, fm)
	if os.IsExist(err) {
		return "", errors.New("** this directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	return dir, nil
}

func GetDir(dir string) (string, error) {
	ckrDir, err := os.Getwd()
	CheckErr(err, "cannot get working directory", 0)

	dir = path.Join(ckrDir, dir)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return "", errors.New("** this directory does not exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	return dir, nil
}

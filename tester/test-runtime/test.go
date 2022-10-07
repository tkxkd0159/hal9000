package main

import (
	"fmt"

	rpc "github.com/Carina-labs/HAL9000/rpc/types"
)

func main() {
	fmt.Println("tester")
	tep := rpc.NewTypedEventParser("aaa", "sss")
	fmt.Println(tep)
}

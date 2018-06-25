package main

import (
	"fmt"
)

func main() {
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}

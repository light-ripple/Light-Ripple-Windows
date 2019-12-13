package main

import (
	"fmt"

	"zxq.co/ripple/agplwarning"
)

func main() {
	err := agplwarning.Warn("agplwarning", "AGPLWarning")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("License agreed")
}

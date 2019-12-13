package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.zxq.co/ripple/schiavolib"
)

var (
	url      = "general"
	messages = make(chan string, 20)
)

func main() {
	schiavo.ForceDo = true
	for i := 0; i < 10; i++ {
		go sender()
	}
	fmt.Println("schiavosay")
	fmt.Print("> ")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if strings.Index(sc.Text(), "/switch ") == 0 {
			url = sc.Text()[len("/switch "):]
			fmt.Println("=> Switched to", url)
			fmt.Print("> ")
			continue
		}
		messages <- sc.Text()
		fmt.Print("> ")
	}
}

func sender() {
	for m := range messages {
		err := schiavo.Channel(url).Send(m)
		if err != nil {
			fmt.Println()
			fmt.Println("=>", err)
			fmt.Print("> ")
		}
	}
}

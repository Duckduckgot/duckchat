package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Duckduckgot/duckchat"
)

func main() {
	agent := duckchat.NewAgent(duckchat.GPT4o)

	for {
		fmt.Print("user > ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := agent.Send(line)
		if err != nil {
			fmt.Println("Unexpected Error:", err.Error())
			continue
		}

		fmt.Println("Assistant >", response)
	}
}

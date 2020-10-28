package app

import "fmt"

type ConsoleMessage struct {
	Message string
	Process *Process
}

type ConsoleMessageChannel chan *ConsoleMessage

func (c *ConsoleMessageChannel) PrintRelevant(processes ProcessList) {
	for message := range *c {
		for _, p := range processes {
			if p == message.Process {
				fmt.Printf(message.Message)
			}
		}
	}
}

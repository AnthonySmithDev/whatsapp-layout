package main

import (
	"github.com/AnthonySmithDev/whatsapp-tui-layout/tui"
	"github.com/AnthonySmithDev/whatsapp-tui-layout/ws"
)

func main() {
	// channel := make(chan *whatsmeow.Client)
	client := ws.RunWs()
	tui.NewProgram(client)
}
package main

import (
	"github.com/AnthonySmithDev/whatsapp-tui-layout/app"
	"github.com/AnthonySmithDev/whatsapp-tui-layout/tui"
)

func main() {
	client := app.Connect()
	app.Handler(client)
	tui.NewProgram(client)
}
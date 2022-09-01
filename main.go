package main

import (
	// "time"

	"github.com/AnthonySmithDev/whatsapp-tui-layout/app"
	"github.com/AnthonySmithDev/whatsapp-tui-layout/tui"
)

func main() {
	app.Connect()
	app.NewDatabase()
	app.Handler()
	// time.Sleep(10 * time.Second)
	tui.NewProgram()
}
package main

import (
	"github.com/AnthonySmithDev/whatsapp-tui-layout/app"
)

func main() {
	first := app.Connect()
	app.Database()
	app.Handler()
	app.NewTui(first)
}
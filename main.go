package main

import (
	"GatewayWorker/console"
	"github.com/ctfang/command"
)

func main() {
	app := command.New()

	app.SetConfig("D:\\GoLanProject\\GatewayWorker\\config.ini")
	app.IniConfig()

	AddCommands(&app)
	app.Run()
}

func AddCommands(app *command.Console) {
	app.AddCommand(console.Gateway{})
	app.AddCommand(console.Register{})
	app.AddCommand(console.Worker{})
	app.AddCommand(console.Test{})
}

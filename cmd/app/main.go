package main

import (
	"backend/config"
	"backend/internal/app"
)

func main() {
	app.Run(config.New())
}

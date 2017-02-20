package main

import (
	"github.com/reechou/robot-account/config"
	"github.com/reechou/robot-account/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}

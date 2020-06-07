package main

import (
	"github.com/kreempuff/geoffrey/orchestrator"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	a := orchestrator.CreateCLI()
	err := a.Run(os.Args)
	if err != nil {
		logrus.Error(err)
	}
}

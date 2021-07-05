package main

import (
	"github.com/consensys/quorum-key-manager/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	command := cmd.NewCommand()
	if err := command.Execute(); err != nil {
		log.WithError(err).Fatalf("execution failed")
	}
	log.Info("execution completed")
}

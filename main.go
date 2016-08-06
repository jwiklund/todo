package main

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/docopt/docopt-go"
	_ "github.com/pkg/errors"
)

func main() {
	log.WithField("component", "main").Debug("Startup")
}

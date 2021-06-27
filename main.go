package main

import (
	log "github.com/sirupsen/logrus"
	"technopark-dbms/internal/app"
)

func main() {
	log.SetLevel(log.FatalLevel)
	app.RunServer(":5000")
}

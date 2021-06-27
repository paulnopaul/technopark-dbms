package main

import (
	log "github.com/sirupsen/logrus"
	"technopark-dbms/internal/app"
)

func main() {
	log.SetLevel(log.ErrorLevel)
	app.RunServer(":5000")
}

package main

import (
	"github.com/fellah/stop"

	"github.com/fellah/tcache/jobs"
	"github.com/fellah/tcache/log"
)

func main() {
	go jobs.Start()

	<-stop.Ch

	jobs.End()
	log.Info.Println("EXIT")
}
